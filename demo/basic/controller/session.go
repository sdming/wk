// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
session demo

*/
package controller

import (
	"fmt"
	"github.com/sdming/kiss/kson"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/boot"
	"github.com/sdming/wk/session"
	"log"
	"time"
)

type Session struct {
}

func NewSession() *Session {
	return &Session{}
}

var sessionDemo *Session

func init() {
	boot.Boot(RegisterSessionRoute)
}

func RegisterSessionRoute(server *wk.HttpServer) {
	sessionDemo = NewSession()

	// url: /session/xxx/xxx
	// route to Session
	server.RouteTable.Path("/session/{action}").ToController(sessionDemo)
}

// url: /session/id
func (c *Session) Id(ctx *wk.HttpContext) (wk.HttpResult, error) {
	id := ctx.SessionId()
	return wk.Data(id), nil
}

// url: /session/add?k=test&v=101
func (c *Session) Add(ctx *wk.HttpContext) (wk.HttpResult, error) {
	ok, err := ctx.Session.Add(ctx.FV("k"), ctx.FV("v"))
	return wk.Data(ok), err
}

// url: /session/get?k=test
func (c *Session) Get(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	v, _, err := ctx.Session.Get(ctx.FV("k"))
	return wk.Data(v), err
}

// url: /session/set?k=test&v=101
func (c *Session) Set(ctx *wk.HttpContext) (wk.HttpResult, error) {
	err := ctx.Session.Set(ctx.FV("k"), ctx.FV("v"))
	return wk.Data(err == nil), err
}

// url: /session/remove?k=test
func (c *Session) Remove(ctx *wk.HttpContext) (wk.HttpResult, error) {
	err := ctx.Session.Remove(ctx.FV("k"))
	return wk.Data(err == nil), err
}

// url: /session/abandon
func (c *Session) Abandon(ctx *wk.HttpContext) (wk.HttpResult, error) {
	err := ctx.Session.Abandon()
	return wk.Data(true), err
}

// url: /session/keys
func (c *Session) Keys(ctx *wk.HttpContext) (wk.HttpResult, error) {
	keys, err := ctx.Session.Keys()
	return wk.Data(fmt.Sprintln(keys)), err
}

type DebugDriver struct {
	Option *Option
}

type Option struct {
	Int     int
	String  string
	Float64 float64
	Bool    bool
	Map     map[string]string
	No      string
}

func newOption() *Option {
	return &Option{
		Int:     1,
		String:  "string",
		Float64: 1.1,
		Bool:    false,
		No:      "no config",
	}
}

func newDebugDriver() *DebugDriver {
	return &DebugDriver{Option: newOption()}
}

func (d *DebugDriver) Name() string {
	log.Println("Name")
	return "session_debug"
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
	return true, nil
}

func (d *DebugDriver) Keys(sessionId string) ([]string, error) {
	log.Println("Keys", sessionId)
	return make([]string, 0), nil
}

func (d *DebugDriver) Init(options string) error {
	log.Println("Init", options)
	kson.Unmarshal([]byte(options), d.Option)
	log.Printf("option: %#v", d.Option)
	return nil
}

func init() {
	session.Register("session_debug", newDebugDriver())
}
