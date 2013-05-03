// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"errors"
	"github.com/sdming/wk/session"
	"net/http"
	"time"
)

var cookieSessionId = "_go_session_"

var SessionDriver session.Driver

// setupSession
func (srv *HttpServer) configSession() error {
	if !srv.Config.SessionEnable {
		Logger.Println("Session: Enable=", srv.Config.SessionEnable)
		return nil
	}

	name := srv.Config.SessionDriver
	if name == "" {
		name = defaultSessionDriver
	}
	SessionDriver = session.GetDriver(name)
	if SessionDriver == nil {
		return errors.New("session driver is nil:" + name)
	}

	//name = SessionDriver.Name()
	if srv.Config.PluginConfig != nil {
		if node, ok := srv.Config.PluginConfig.Child(name); ok && node != nil {
			SessionDriver.Init(node.Dump())
		} else {
			SessionDriver.Init("")
		}
	} else {
		SessionDriver.Init("")
	}

	Logger.Printf("Session: Enable=%v; Timeout =%d; Driver=%s %v \n",
		srv.Config.SessionEnable, srv.Config.SessionTimeout, srv.Config.SessionDriver, SessionDriver)

	return nil
}

// setupSession
func (srv *HttpServer) exeSession(ctx *HttpContext) {
	if !srv.Config.SessionEnable {
		return
	}

	var id string

	cookie, err := ctx.Cookie(cookieSessionId)
	if err != nil {
		//Logger.Println("ctx.Cookie error, create new session", err)
		srv.newSession(ctx)
	} else {
		id = cookie.Value
		if ok, err := SessionDriver.Exists(id); !ok || err != nil {
			srv.newSession(ctx)
		} else {
			ctx.Session = Session(id)
			ctx.SessionIsNew = false
		}
	}
}

func (srv *HttpServer) newSession(ctx *HttpContext) {
	id := session.NewId()
	err := SessionDriver.New(id, time.Duration(srv.Config.SessionTimeout)*time.Second)
	if err != nil {
		return
	}

	cookie := &http.Cookie{
		Name:     cookieSessionId,
		Value:    id,
		Path:     `/`,
		HttpOnly: true,
	}
	//Logger.Println("cookie", cookie)
	ctx.SetCookie(cookie)
	ctx.SessionIsNew = true
	ctx.Session = Session(id)
}

type Session string

func (s Session) Add(key string, value interface{}) (bool, error) {
	return SessionDriver.Add(string(s), key, value)
}

func (s Session) Get(key string) (interface{}, bool, error) {
	return SessionDriver.Get(string(s), key)
}

func (s Session) Set(key string, value interface{}) error {
	return SessionDriver.Set(string(s), key, value)
}

func (s Session) Remove(key string) error {
	return SessionDriver.Remove(string(s), key)
}

func (s Session) Abandon() error {
	return SessionDriver.Abandon(string(s))
}

func (s Session) Keys() ([]string, error) {
	return SessionDriver.Keys(string(s))
}
