package session_test

import (
	"github.com/sdming/wk/session"
	"log"
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

var driverName = "default"

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

type Driver interface {
	// Initialize setup this driver
	Init(options string) error

	// return name of this driver
	Name() string

	// Get return item value of session marked with sessionId, return false if key exists
	Add(sessionId, key string, value interface{}) (bool, error)

	// Get return item value of session marked with sessionId, return false if key doesn't exists
	Get(sessionId, key string) (interface{}, bool, error)

	// Set add or update a item 
	Set(sessionId, key string, value interface{}) error

	// Remove remove item from session
	Remove(sessionId, key string) error

	// New create a new session entry
	New(sessionId string, timeout time.Duration) error

	// Abandon mark session abandon
	Abandon(sessionId string) error

	// Exists return false if sessionId doesn't exists
	Exists(sessionId string) (bool, error)

	// Keys all item keys of session
	Keys(sessionId string) ([]string, error)
}

type DebugDriver struct {
}

type Options struct {
	Int    int
	String string
	Int64  int64
	Float  float64
	Slice  []string
}

func newDebugDriver() *DebugDriver {
	return &DebugDriver{}
}

func (d *DebugDriver) Name() string {
	log.Println("Name")
	return "debug"
}

func (d *DebugDriver) Add(sessionId, key string, value interface{}) (bool, error) {
	log.Println("Add", sessionId, key, value)
	return true, nil
}

func (d *DebugDriver) Get(sessionId, key string) (interface{}, bool, error) {
	log.Println("Get", sessionId, key)
	return nil, false, nil
}

func (d *DebugDriver) Set(sessionId, key string, value interface{}) error {
	log.Println("Set", sessionId, key, value)
	return nil
}

func (d *DebugDriver) Remove(sessionId, key string) error {
	log.Println("Remove", sessionId, key)
	return nil
}

func (d *DebugDriver) New(sessionId string, timeout time.Duration) error {
	log.Println("New", sessionId, timeout)
	return nil
}

func (d *DebugDriver) Abandon(sessionId string) error {
	log.Println("Abandon", sessionId)
	return nil
}

func (d *DebugDriver) Exists(sessionId string) (bool, error) {
	log.Println("Exists", sessionId)
	return false, nil
}

func (d *DebugDriver) Keys(sessionId string) ([]string, error) {
	log.Println("Keys", sessionId)
	return make([]string, 0), nil
}

func (d *DebugDriver) Init(options string) error {
	log.Println("Init", options)
	return nil
}

func init() {
	session.Register("debug", newDebugDriver())
}
