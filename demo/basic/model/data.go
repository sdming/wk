// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic demo

*/
package model

import (
	"errors"
	"fmt"
	"github.com/sdming/kiss/kson"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/boot"
	"strconv"
)

type Data struct {
	Str   string
	Uint  uint64
	Int   int
	Float float32
	Byte  byte
}

func newData(i int) *Data {
	return &Data{
		Str:   "string:" + strconv.Itoa(i),
		Uint:  uint64(i * 100),
		Int:   i * 10,
		Float: 0.1 + float32(i),
		Byte:  byte(i),
	}
}

func DataTop(count int) []*Data {
	if count <= 0 || count > 100 {
		panic("count is invalid")
	}

	data := make([]*Data, count)
	for i := 0; i < count; i++ {
		data[i] = newData(i)
	}

	return data
}

func DataByInt(i int) *Data {
	if i < 0 {
		i = 0
	}
	return newData(i)
}

func DataByIntRange(start, end int) []*Data {
	data := make([]*Data, 2)
	data[0] = newData(start)
	data[1] = newData(end)
	return data
}

func DataSet(s string, u uint64, i int, f float32, b byte) *Data {
	return &Data{
		Str:   s,
		Uint:  u,
		Int:   i,
		Float: f,
		Byte:  b,
	}
}

func DataAnonymous(data struct {
	Str   string
	Uint  uint64
	Int   int
	Float float32
	Byte  byte
},) string {

	return fmt.Sprintln(data)
}

func DataDelete(i int) string {
	return fmt.Sprintf("delete %d", i)
}

func DataPost(data Data) string {
	return data.String()
}

func DataPostPtr(data *Data) string {
	return data.String()
}

func (d *Data) String() string {
	if d == nil {
		return "<nil>"
	}
	return fmt.Sprintf("string:%s;uint:%d;int:%d;float:%f;byte:%c",
		d.Str, d.Uint, d.Int, d.Float, d.Byte)
}

func DataTopHandle(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	if count, ok := ctx.RouteData.Int("count"); !ok {
		err = errors.New("parameter invalid:" + "count")
	} else {
		data := DataTop(count)
		result = wk.Json(data)
	}
	return
}

func DataType() Data {
	return Data{
		Str:   "string",
		Uint:  64,
		Int:   32,
		Float: 3.14,
		Byte:  8,
	}
}

func init() {
	boot.Boot(RegisterDataRoute)
}

func RegisterDataRoute(server *wk.HttpServer) {
	// url: /data/top/10
	// func: DataTopHandle(ctx *wk.HttpContext) (result wk.HttpResult, err error)
	// route to func (*wk.HttpContext) (wk.HttpResult, error)
	server.RouteTable.Get("/data/top/{count}").To(DataTopHandle)

	// url: /data/datatype
	// func: DataType() Data
	// return content-type according to accepts
	server.RouteTable.Get("/data/datatype?").ToFunc(DataType)

	// url: /data/int/1
	// func: DataByInt(i int) *Data
	// route to a function, convert parameter by index(p0,p1,p2...)
	server.RouteTable.Get("/data/int/{p0}?").ToFunc(DataByInt)

	// url: /data/range/1-9
	// func: DataByIntRange(start, end int) []*Data
	// route to a function, convert parameter by index(p0,p1,p2...)
	server.RouteTable.Get("/data/range/{p0}-{p1}").ToFunc(DataByIntRange)

	// url: /data/int/1/xml
	// func: DataByInt(i int) *Data
	// return xml
	server.RouteTable.Get("/data/int/{p0}/xml").ToFunc(DataByInt).ReturnXml()

	// url: /data/int/1/json
	// func: DataByInt(i int) *Data
	// return json
	server.RouteTable.Get("/data/int/{p0}/json").ToFunc(DataByInt).ReturnJson()

	// url: /data/int/1/kson
	// func: DataByInt(i int) *Data
	// return custome formatted data
	server.RouteTable.Get("/data/int/{p0}/kson").ToFunc(DataByInt).Return(formatKson)

	// url: /data/name/1
	// func: DataByInt(i int) *Data
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/name/{id}").ToFunc(DataByInt).
		BindByNames("id")

	// url: /data/namerange/1-9
	// func: DataByIntRange(start, end int) []*Data
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/namerange/{start}-{end}").ToFunc(DataByIntRange).
		BindByNames("start", "end")

	// url: /data/namerange/?start=1&end=9
	// func: DataByIntRange(start, end int) []*Data
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/namerange/").ToFunc(DataByIntRange).
		BindByNames("start", "end")

	// url: post /data/post?
	// form:{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}}
	// func: DataPost(data Data) string
	// route http post to function, build struct parameter from form
	server.RouteTable.Post("/data/post?").ToFunc(DataPost).BindToStruct()

	// url: post /data/postptr?
	// form:{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}}
	// func DataPostPtr(data *Data) string
	// route http post to function, build struct parameter from form
	server.RouteTable.Post("/data/postptr?").ToFunc(DataPostPtr).BindToStruct()

	// url: delete /data/delete/1
	// func: DataDelete(i int) string
	// route http delete to function
	server.RouteTable.Delete("/data/delete/{p0}").ToFunc(DataDelete)

	// url: get /data/set?str=string&uint=1024&int=32&float=3.14&byte=64
	// func: DataSet(s string, u uint64, i int, f float32, b byte) *Data
	// test diffrent parameter type
	server.RouteTable.Get("/data/set?").ToFunc(DataSet).
		BindByNames("str", "uint", "int", "float", "byte")

	// url: get /data/anonymous?str=string&uint=1024&int=32&float=3.14&byte=64
	// func: DataAnonymous(...) string
	// test anonymous struct
	server.RouteTable.Get("/data/anonymous?").ToFunc(DataAnonymous).BindToStruct()

}

func formatKson(ctx *wk.HttpContext, x interface{}) (wk.HttpResult, bool) {
	b, _ := kson.Marshal(x)
	return wk.Content(string(b), "text/plain"), true
}
