// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic demo

*/
package controller

import (
	"encoding/json"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/model"
	"log"
)

type DemoController struct {
	data []*model.Data
}

func NewDemoController() *DemoController {
	return &DemoController{
		data: make([]*model.Data, 0),
	}
}

func (c *DemoController) deleteByInt(v int) {
	for i := 0; i < len(c.data); {
		if c.data[i].Int == v {
			c.data = append(c.data[:i], c.data[i+1:]...)
		} else {
			i++
		}
	}
}

func (c *DemoController) getByInt(v int) []*model.Data {
	data := make([]*model.Data, 0)
	for i := 0; i < len(c.data); i++ {
		if c.data[i].Int == v {
			data = append(data, c.data[i])
		}
	}
	return data
}

var enableLog bool = false

// url: /demo/all/
func (c *DemoController) All(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	if enableLog {
		log.Println("DemoController All")
	}

	return wk.Json(c.data), nil
}

// url: /demo/add/?int=32&str=string&uint=1024&float=1.1&byte=64
func (c *DemoController) Add(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	if enableLog {
		log.Println("DemoController Add")
	}

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

// url: /demo/delete/32
func (c *DemoController) Delete(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	if enableLog {
		log.Println("DemoController Delete")
	}

	l := len(c.data)
	if i, ok := ctx.RouteData.Int("id"); ok {
		c.deleteByInt(i)
	}
	return wk.Data(l - len(c.data)), nil
}

// url: /demo/put/32?str=s&uint=64&float=3.14&byte=8
func (c *DemoController) Put(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	if enableLog {
		log.Println("DemoController Put")
	}

	var id int

	if i, ok := ctx.RouteData.Int("id"); ok {
		id = i
	} else {
		return wk.Data(0), nil
	}

	var count int = 0
	for _, d := range c.data {
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

// url: /demo/int/32
func (c *DemoController) Int(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	if enableLog {
		log.Println("DemoController Int")
	}

	if id, ok := ctx.RouteData.Int("id"); ok {
		return wk.Json(c.getByInt(id)), nil
	}
	return wk.Data(""), nil
}

// url: /demo/rangecount/?start=1&end=99
func (c *DemoController) RangeCount(ctx *wk.HttpContext) (result wk.HttpResult, err error) {

	if enableLog {
		log.Println("DemoController RangeCount")
	}

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

// url: post /demo/post
func (c *DemoController) Post(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	if enableLog {
		log.Println("DemoController Post")
	}

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

// url: /demo/clear/
func (c *DemoController) Clear(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	if enableLog {
		log.Println("DemoController Clear")
	}

	c.data = make([]*model.Data, 0)
	return wk.Data(len(c.data)), nil
}

func (c *DemoController) OnActionExecuting(action *wk.ActionContext) {
	if enableLog {
		log.Println("DemoController OnActionExecuting", action.Name)
	}
}

func (c *DemoController) OnActionExecuted(action *wk.ActionContext) {
	if enableLog {
		log.Println("DemoController OnActionExecuted", action.Name)
	}
}

func (c *DemoController) OnException(action *wk.ActionContext) {
	if enableLog {
		log.Println("DemoController OnActionExecuted", action.Name, action.Err)
	}
}
