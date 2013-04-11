// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic demo

*/
package main

import (
	"encoding/json"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/model"
)

type DemoController struct {
	data []*model.Data
}

func newDemoController() *DemoController {
	return &DemoController{
		data: make([]*model.Data, 0),
	}
}

func (c *DemoController) deleteByInt(v int) {
	for i := 0; i < len(c.data); i++ {
		if c.data[i].Int == v {
			c.data = append(c.data[:i], c.data[i+1:]...)
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

// url: post /demo/post
func (c *DemoController) Post(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
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

// url: /demo/all
func (c *DemoController) All(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.Json(c.data), nil
}

// url: /demo/delete/1
func (c *DemoController) Delete(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	l := len(c.data)
	if i, ok := ctx.RouteData.Int("id"); ok {
		c.deleteByInt(i)
	}
	return wk.Data(l - len(c.data)), nil
}

// url: /demo/put/1?str=string&uint=1024&float=1.1&byte=64
func (c *DemoController) Put(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
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

// url: /demo/clear
func (c *DemoController) Clear(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	c.data = make([]*model.Data, 0)
	return wk.Data(len(c.data)), nil
}

// url: /demo/add?int=32&str=string&uint=1024&float=1.1&byte=64
func (c *DemoController) Add(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	data := &model.Data{
		Int:   ctx.FormIntOr("int", 0),
		Uint:  uint64(ctx.FormIntOr("uint", 0)),
		Str:   ctx.FormValue("str"),
		Float: float32(ctx.FormFloatOr("float", 0.0)),
		Byte:  byte(ctx.FormIntOr("byte", 0)),
	}
	c.data = append(c.data, data)
	return wk.Data(len(c.data)), nil
}

// url: /demo/int/1
func (c *DemoController) Int(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	if id, ok := ctx.RouteData.Int("id"); ok {
		return wk.Json(c.getByInt(id)), nil
	}
	return wk.Data(""), nil
}

// url: /demo/RangeCount?start=1&end=9
func (c *DemoController) RangeCount(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	start := ctx.RouteData.IntOr("start", 0)
	end := ctx.RouteData.IntOr("end", 0)
	var count int = 0

	for _, d := range c.data {
		if d.Int >= start && d.Int <= end {
			count++
		}
	}
	return wk.Data(count), nil
}
