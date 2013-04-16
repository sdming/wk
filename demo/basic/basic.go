// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic demo

*/
package main

import (
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
		fmt.Println("NewDefaultServer error", err)
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

	//demo, show to define custome httpresult
	model.RegisterQrRoute(server)

	enableEventTrace := false
	if enableEventTrace {
		model.RegisterEventTrace(server)
	}

	enableCompress := false
	if enableCompress {
		server.Processes.InsertBefore("_render", wk.NewCompressProcess("compress_test", "*", "/js/"))
	}

	enableFile := true
	if enableFile {
		model.RegisterFileRoute(server)
	}

	enableBigpipe := true
	if enableBigpipe {
		model.RegisterBigPipeRoute(server)
	}

	server.Start()

}

func formatKson(ctx *wk.HttpContext, x interface{}) (wk.HttpResult, bool) {
	b, _ := kson.Marshal(x)
	return wk.Content(string(b), "text/plain"), true
}
