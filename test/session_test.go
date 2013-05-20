package wk_test

import (
	"github.com/sdming/wk/session"
	"testing"
	"time"
)

func TestSessionId(t *testing.T) {
	id := session.NewId()
	t.Log("NewId", id)
	if id == "" {
		t.Error("session id error")
	}
}

var driverName = "session_default"

func TestSessionBasic(t *testing.T) {
	id := session.NewId()
	key := "key"
	value := "value"

	var ok bool
	var err error

	driver := session.GetDriver(driverName)
	if driver == nil {
		t.Error("GetDriver", "return nil", driverName)
		return
	}

	if name := driver.Name(); name != driverName {
		successAndEq(t, "Name", nil, driverName, name)
	}

	if err := driver.New(id, time.Second); err != nil {
		t.Error(t, "New", err)
		return
	}

	ok, err = driver.Add(id, key, value)
	successAndEq(t, "Add", err, true, ok)

	ok, err = driver.Add(id, key, value)
	successAndEq(t, "Add Again", err, false, ok)

	v, ok, err := driver.Get(id, key)
	successAndEq(t, "Get err", err, true, ok)
	successAndEq(t, "Get ok", err, value, v)

	value = value + "_new"
	err = driver.Set(id, key, value)
	successAndEq(t, "Set", err, nil, nil)

	v, ok, err = driver.Get(id, key)
	successAndEq(t, "Get after set", err, value, v)

	keys, err := driver.Keys(id)
	successAndEq(t, "Keys", err, key, keys[0])

	err = driver.Remove(id, key)
	successAndEq(t, "Remove", err, nil, nil)

	v, ok, err = driver.Get(id, key)
	successAndEq(t, "Get after remove", err, false, ok)
	successAndEq(t, "Get after remove", err, nil, v)

	keys, err = driver.Keys(id)
	successAndEq(t, "Keys after remove", err, 0, len(keys))

	ok, err = driver.Exists(id)
	successAndEq(t, "Exists", err, true, ok)

	err = driver.Abandon(id)
	successAndEq(t, "Abandon", err, nil, nil)

	ok, err = driver.Exists(id)
	successAndEq(t, "Exists after Abandon", err, false, ok)
}
