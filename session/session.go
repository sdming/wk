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

/*


1: register
2: plugin.conf -- configable
3: web.conf -- type? off ??
4: debug -- 


web.conf

session_enable
session_timeout
session_mode

*/

type Session struct {
	CreateAt time.Time
	Id       string
}

// session enable
//参数设置？？
//Expire 过期？？
//IsNewSession	获取一个值，该值指示会话是否是与当前请求一起创建的。
//SessionID
//timeout

// Add,

// <sessionState
//   Mode="InProc"
//   stateConnectionString="tcp=127.0.0.1:42424"
//   stateNetworkTimeout="10"
//   sqlConnectionString="data source=127.0.0.1;Integrated Security=SSPI"
//   sqlCommandTimeout="30"
//   customProvider=""
//   cookieless="false"
//   regenerateExpiredSessionId="false"
//   timeout="20"
//   sessionIDManagerType="Your.ID.Manager.Type,
//     CustomAssemblyNameInBinFolder"
// />

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
