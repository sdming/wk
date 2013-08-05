// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic basic

*/
package controller

import (
	"github.com/sdming/kiss/kson"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/boot"
)

type Driver struct {
	Driver   string
	Host     string
	User     string
	Password string
	A        string
	B        string
}

type Config struct {
	Log_Level string
	Listen    uint
	Roles     []Role
	Db_Log    Db
	Env       map[string]string
}

type Role struct {
	Name  string
	Allow []string
	Deny  []string
}

type Db struct {
	Driver   string
	Host     string
	User     string
	Password string
}

type ConfigController struct {
	webConfig *wk.WebConfig
}

func (c *ConfigController) node() *kson.Node {
	return c.webConfig.AppConfig
}

var config *ConfigController

func init() {
	boot.Boot(RegisterConfigRoute)
}

func RegisterConfigRoute(server *wk.HttpServer) {
	config = &ConfigController{}
	//config.node = server.Config.AppConfig
	config.webConfig = server.Config

	server.RouteTable.Path("/config/{action}").ToController(config)
}

// url: /config/refresh
func (c *ConfigController) Refresh(ctx *wk.HttpContext) (wk.HttpResult, error) {
	s := c.node().ChildString("key_string")
	return wk.Data(s), nil
}

// url: /config/dump
func (c *ConfigController) Dump(ctx *wk.HttpContext) (wk.HttpResult, error) {
	return wk.Data(c.node().MustChild("key_config").Dump()), nil
}

// url: /config/child
func (c *ConfigController) Child(ctx *wk.HttpContext) (wk.HttpResult, error) {
	_, ok := c.node().Child("key_string")
	return wk.Data(ok), nil
}

// url: /config/query
func (c *ConfigController) Query(ctx *wk.HttpContext) (wk.HttpResult, error) {
	// maybe support "a b[@field=xxx] c d " later
	n, ok := c.node().Query("key_config Db_Log Host")
	if ok {
		return wk.Data(n.Literal), nil
	}
	return wk.Data(ok), nil
}

// url: /config/childstringordefault
func (c *ConfigController) ChildStringOrDefault(ctx *wk.HttpContext) (wk.HttpResult, error) {
	//ChildIntOrDefault, ChildUintOrDefault, ChildFloatOrDefault, ChildBoolOrDefault, ChildStringOrDefault
	s := c.node().ChildStringOrDefault("key_string_not", "default value")
	return wk.Data(s), nil
}

// url: /config/childint
func (c *ConfigController) ChildInt(ctx *wk.HttpContext) (wk.HttpResult, error) {
	//ChildInt, ChildUint, ChildFloat, ChildBool, ChildString
	i := c.node().ChildInt("key_int")
	return wk.Data(i), nil
}

// url: /config/bool
func (c *ConfigController) Bool(ctx *wk.HttpContext) (wk.HttpResult, error) {
	//Int, Uint, Float, Bool, String
	b, err := c.node().MustChild("key_bool").Bool()
	if err != nil {
		return nil, err
	}
	return wk.Data(b), nil
}

// url: /config/slice
func (c *ConfigController) Slice(ctx *wk.HttpContext) (wk.HttpResult, error) {
	data, err := c.node().MustChild("key_array").Slice()
	if err != nil {
		return nil, err
	}
	return wk.Data(data), nil
}

// url: /config/map
func (c *ConfigController) Map(ctx *wk.HttpContext) (wk.HttpResult, error) {
	data, err := c.node().MustChild("key_map").Map()
	if err != nil {
		return nil, err
	}
	return wk.Data(data), nil
}

// url: /config/value
func (c *ConfigController) Value(ctx *wk.HttpContext) (wk.HttpResult, error) {
	v := Driver{
		Driver:   "driver",
		Host:     "host",
		User:     "user",
		Password: "password",
		A:        "aaa",
		B:        "bbb",
	}

	err := c.node().MustChild("key_struct").Value(&v)
	if err != nil {
		return nil, err
	}
	return wk.Data(v), nil
}

// url: /config/composite
func (c *ConfigController) Composite(ctx *wk.HttpContext) (wk.HttpResult, error) {
	conf := &Config{}
	err := c.node().MustChild("key_config").Value(conf)
	if err != nil {
		return nil, err
	}
	return wk.Data(conf), nil
}
