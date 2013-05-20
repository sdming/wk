package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"time"
)

var (
	SessionNotExists = errors.New("session id doesn't exists")
)

const randLength = 20

// return sesstion id
func NewId() string {
	b := make([]byte, randLength)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
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

type Session struct {
	CreateAt time.Time
	Id       string
}

var drivers map[string]Driver = make(map[string]Driver)

// Register register a session driver by the provided name
func Register(name string, driver Driver) {
	if driver == nil {
		panic("Register driver is nil")
	}
	drivers[name] = driver
}

func GetDriver(name string) Driver {
	x, _ := drivers[name]
	return x
}
