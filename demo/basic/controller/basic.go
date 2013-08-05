// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic basic

*/
package controller

import (
	"encoding/json"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/boot"
	"github.com/sdming/wk/demo/basic/model"
)

type BasicController struct {
	data []*model.Data
}

func NewBasicController() *BasicController {
	return &BasicController{
		data: make([]*model.Data, 0),
	}
}

var basic *BasicController

func init() {
	boot.Boot(RegisterBasicRoute)
}

func RegisterBasicRoute(server *wk.HttpServer) {
	basic = NewBasicController()

	// url: /basic/xxx/xxx
	// route to controller
	server.RouteTable.Path("/basic#/{action}/{id}").ToController(basic)
}

func (c *BasicController) deleteByInt(v int) {
	for i := 0; i < len(c.data); {
		if c.data[i].Int == v {
			c.data = append(c.data[:i], c.data[i+1:]...)
		} else {
			i++
		}
	}
}

func (c *BasicController) getByInt(v int) []*model.Data {
	data := make([]*model.Data, 0)
	for i := 0; i < len(c.data); i++ {
		if c.data[i].Int == v {
			data = append(data, c.data[i])
		}
	}
	return data
}

// url: /basic/all/
func (c *BasicController) All(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.Json(c.data), nil
}

// url: /basic/add/?int=32&str=string&uint=1024&float=1.1&byte=64
func (c *BasicController) Add(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	data := &model.Data{
		Int:   ctx.FormIntOr("int", 0),
		Uint:  uint64(ctx.FormIntOr("uint", 0)),
		Str:   ctx.FormValue("str"),
		Float: float32(ctx.FormFloatOr("float", 0.0)),
		Byte:  byte(ctx.FormIntOr("byte", 0)),
	}
	c.data = append(c.data, data)
	return wk.Data(data.String()), nil
}

// url: /basic/delete/32
func (c *BasicController) Delete(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	l := len(c.data)
	if i, ok := ctx.RouteData.Int("id"); ok {
		c.deleteByInt(i)
	}
	return wk.Data(l - len(c.data)), nil
}

// url: /basic/set/32?str=s&uint=64&float=3.14&byte=8
func (c *BasicController) Set(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	var id int

	if i, ok := ctx.RouteData.Int("id"); ok {
		id = i
	} else {
		return wk.Data(0), nil
	}

	var count int = 0
	for i := 0; i < len(c.data); i++ {
		d := c.data[i]

		if d.Int != id {
			continue
		}
		if x := ctx.FormValue("str"); x != "" {
			d.Str = x
		}
		if x, ok := ctx.FormInt("uint"); ok {
			d.Uint = uint64(x)
		}
		if x, ok := ctx.FormFloat("float"); ok {
			d.Float = float32(x)
		}
		if x, ok := ctx.FormInt("byte"); ok {
			d.Byte = byte(x)
		}
		count++
	}

	return wk.Data(count), nil
}

// url: /basic/int/32
func (c *BasicController) Int(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	if id, ok := ctx.RouteData.Int("id"); ok {
		return wk.Json(c.getByInt(id)), nil
	}
	return wk.Data(""), nil
}

// url: /basic/rangecount/?start=1&end=99
func (c *BasicController) RangeCount(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	start := ctx.FormIntOr("start", 0)
	end := ctx.FormIntOr("end", 0)
	var count int = 0

	for _, d := range c.data {
		if d.Int >= start && d.Int <= end {
			count++
		}
	}
	return wk.Data(count), nil
}

// url: post /basic/post/
func (c *BasicController) Post(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	var body []byte
	if body, err = ctx.ReadBody(); err != nil {
		return nil, err
	}

	data := &model.Data{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}
	c.data = append(c.data, data)
	return wk.Data(true), nil
}

// url: /basic/clear/
func (c *BasicController) Clear(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	c.data = make([]*model.Data, 0)
	return wk.Data(len(c.data)), nil
}
