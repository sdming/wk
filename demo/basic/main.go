// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic demo

*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sdming/kiss/kson"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/controller"
	"github.com/sdming/wk/demo/basic/model"
)

func DataTopHandle(ctx *wk.HttpContext) (result wk.HttpResult, err error) {

	if count, ok := ctx.RouteData.Int("count"); !ok {
		err = errors.New("parameter invalid:" + "count")
	} else {
		data := model.DataTop(count)
		result = wk.Json(data)
	}
	return
}

func main() {

	server, err := wk.NewDefaultServer()

	if err != nil {
		fmt.Println("DefaultServer error", err)
		return
	}

	controller := controller.NewDemoController()

	// url: /demo/xxx/xxx
	// route to controller
	server.RouteTable.Path("/demo/{action}/{id}").ToController(controller)

	// url: /data/top/10
	// func: DataTopHandle(ctx *wk.HttpContext) (result wk.HttpResult, err error)
	// route to func (*wk.HttpContext) (wk.HttpResult, error)
	server.RouteTable.Get("/data/top/{count}").To(DataTopHandle)

	// url: /data/int/1
	// func: DataByInt(i int) *Data
	// route to a function, convert parameter by index(p0,p1,p2...)
	server.RouteTable.Get("/data/int/{p0}?").ToFunc(model.DataByInt)

	// url: /data/range/1-9
	// func: DataByIntRange(start, end int) []*Data
	// route to a function, convert parameter by index(p0,p1,p2...)
	server.RouteTable.Get("/data/range/{p0}-{p1}").ToFunc(model.DataByIntRange)

	// url: /data/int/1/xml
	// func: DataByInt(i int) *Data
	// return xml
	server.RouteTable.Get("/data/int/{p0}/xml").ToFunc(model.DataByInt).ReturnXml()

	// url: /data/int/1/json
	// func: DataByInt(i int) *Data
	// return json
	server.RouteTable.Get("/data/int/{p0}/json").ToFunc(model.DataByInt).ReturnJson()

	// url: /data/int/1/kson
	// func: DataByInt(i int) *Data
	// return custome formatted data
	server.RouteTable.Get("/data/int/{p0}/kson").ToFunc(model.DataByInt).Return(formatKson)

	// url: /data/name/1
	// func: DataByInt(i int) *Data
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/name/{id}").ToFunc(model.DataByInt).
		BindByNames("id")

	// url: /data/namerange/1-9
	// func: DataByIntRange(start, end int) []*Data
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/namerange/{start}-{end}").ToFunc(model.DataByIntRange).
		BindByNames("start", "end")

	// url: /data/namerange/?start=1&end=9
	// func: DataByIntRange(start, end int) []*Data
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/namerange/").ToFunc(model.DataByIntRange).
		BindByNames("start", "end")

	// url: post /data/post?
	// form:{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}}
	// func: DataPost(data Data) string 
	// route http post to function, build struct parameter from form  
	server.RouteTable.Post("/data/post?").ToFunc(model.DataPost).BindToStruct()

	// url: post /data/postptr?
	// form:{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}}
	// func DataPostPtr(data *Data) string
	// route http post to function, build struct parameter from form
	server.RouteTable.Post("/data/postptr?").ToFunc(model.DataPostPtr).BindToStruct()

	// url: delete /data/delete/1
	// func: DataDelete(i int) string 
	// route http delete to function
	server.RouteTable.Delete("/data/delete/{p0}").ToFunc(model.DataDelete)

	// url: get /data/set?str=string&uint=1024&int=32&float=3.14&byte=64
	// func: DataSet(s string, u uint64, i int, f float32, b byte) *Data 
	// test diffrent parameter type
	server.RouteTable.Get("/data/set?").ToFunc(model.DataSet).
		BindByNames("str", "uint", "int", "float", "byte")

	server.Start()

}

func formatKson(ctx *wk.HttpContext, x interface{}) (wk.HttpResult, bool) {
	b, _ := kson.Marshal(x)
	return wk.Content(string(b), "text/plain"), true
}

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
