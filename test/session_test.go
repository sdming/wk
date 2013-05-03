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

func test(t *testing.T, m string, err error, expect, actual interface{}) {
	t.Log(m, err, expect, actual)

	if err != nil {
		t.Error(m, "error", err)
		return
	}

	if expect != actual {
		t.Error(m, "expect", expect, actual, actual)
	}
}

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
		test(t, "Name", nil, driverName, name)
	}

	if err := driver.New(id, time.Second); err != nil {
		t.Error(t, "New", err)
		return
	}

	ok, err = driver.Add(id, key, value)
	test(t, "Add", err, true, ok)

	ok, err = driver.Add(id, key, value)
	test(t, "Add Again", err, false, ok)

	v, ok, err := driver.Get(id, key)
	test(t, "Get err", err, true, ok)
	test(t, "Get ok", err, value, v)

	value = value + "_new"
	err = driver.Set(id, key, value)
	test(t, "Set", err, nil, nil)

	v, ok, err = driver.Get(id, key)
	test(t, "Get after set", err, value, v)

	keys, err := driver.Keys(id)
	test(t, "Keys", err, key, keys[0])

	err = driver.Remove(id, key)
	test(t, "Remove", err, nil, nil)

	v, ok, err = driver.Get(id, key)
	test(t, "Get after remove", err, false, ok)
	test(t, "Get after remove", err, nil, v)

	keys, err = driver.Keys(id)
	test(t, "Keys after remove", err, 0, len(keys))

	ok, err = driver.Exists(id)
	test(t, "Exists", err, true, ok)

	err = driver.Abandon(id)
	test(t, "Abandon", err, nil, nil)

	ok, err = driver.Exists(id)
	test(t, "Exists after Abandon", err, false, ok)
}
