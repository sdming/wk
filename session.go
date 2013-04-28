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
	Session := session.GetDriver(name)
	if Session == nil {
		return errors.New("session driver is nil:" + name)
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
		Logger.Println("ctx.Cookie error, create new session", err)
		id = session.NewId()
		Logger.Println("session.NewId()", id)
		err := SessionDriver.New(id, time.Duration(srv.Config.SessionTimeout)*time.Second)
		if err != nil {
			Logger.Println("error", "Session.New", err)
			return
		}

		cookie = &http.Cookie{
			Name:  cookieSessionId,
			Value: id,
			Path:  `/`,
			//Domain     string
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
		}
		Logger.Println("cookie", cookie)
		ctx.SetCookie(cookie)

	} else {
		id = cookie.Value
	}

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
