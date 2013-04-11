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

	controller := newDemoController()

	// url: /demo/xxx/xxx
	// route to controller
	server.RouteTable.Path("/demo/{action}/{id}").ToController(controller)

	// url: /data/top/10
	// route to func (*wk.HttpContext) (wk.HttpResult, error)
	server.RouteTable.Get("/data/top/{count}").To(DataTopHandle)

	// url: /data/int/1
	// route to a function, convert parameter by index(p0,p1,p2...)
	server.RouteTable.Get("/data/int/{p0}?").ToFunc(model.DataByInt)

	// url: /data/range/1-9
	// route to a function, convert parameter by index(p0,p1,p2...)
	server.RouteTable.Get("/data/range/{p0}-{p1}").ToFunc(model.DataByIntRange)

	// url: /data/int/1/xml
	// return xml
	server.RouteTable.Get("/data/int/{p0}/xml").ToFunc(model.DataByInt).ReturnXml()

	// url: /data/int/1/json
	// return json
	server.RouteTable.Get("/data/int/{p0}/json").ToFunc(model.DataByInt).ReturnJson()

	// url: /data/int/1/kson
	// return formatted data
	server.RouteTable.Get("/data/int/{p0}/kson").ToFunc(model.DataByInt).Return(formatKson)

	// url: /data/name/1
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/name/{id}").ToFunc(model.DataByInt).
		BindByNames("id")

	// url: /data/namerange/1-9
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/namerange/{start}-{end}").ToFunc(model.DataByIntRange).
		BindByNames("start", "end")

	// url: /data/namerange/?start=1&end=9
	// route to a function, convert parameter by name
	server.RouteTable.Get("/data/namerange/").ToFunc(model.DataByIntRange).
		BindByNames("start", "end")

	// url: post /data/post?
	// form:{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}}
	server.RouteTable.Post("/data/post?").ToFunc(model.DataPost).BindToStruct()

	// url: post /data/postptr?
	// form:{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}}
	server.RouteTable.Post("/data/postptr?").ToFunc(model.DataPostPtr).BindToStruct()

	// url: delete /data/delete/1
	server.RouteTable.Delete("/data/delete/{p0}").ToFunc(model.DataDelete)

	// url: get /data/set?str=string&uint=1024&int=32&float=3.14&byte=64
	server.RouteTable.Get("/data/set?").ToFunc(model.DataSet).
		BindByNames("str", "uint", "int", "float", "byte")

	server.Start()

}

func formatKson(ctx *wk.HttpContext, x interface{}) (wk.HttpResult, bool) {
	b, _ := kson.Marshal(x)
	return wk.Content(string(b), "text/plain"), true
}

// package user

// import (
// 	"fmt"
// 	"strconv"
// 	"github.com/sdming/gomvc"
// )

// //controller
// type UserController struct {
// }

// //get http://localhost:8080/user/user/1000 (application/json; application/xml;text/plain;text/html)
// //return struct, marshal to json/xml/text accoring to "Accept"
// func (*UserController) User(id int) User {
// 	return GetById(id)
// }

// //get http://localhost:8080/user/string/1000
// //return string
// func (*UserController) String(id int) string {
// 	return GetById(id).String()
// }

// //get http://localhost:8080/user/int/1000
// //return int
// func (*UserController) Int(id int) int {
// 	return GetById(id).Id
// }

// //get http://localhost:8080/user/json/1000
// //return json
// func (*UserController) Json(id int) gomvc.HttpResult {
// 	u := GetById(id)
// 	return gomvc.Json(u)
// }

// //get http://localhost:8080/user/xml/1000
// //return xml
// func (*UserController) Xml(id int) gomvc.HttpResult {
// 	u := GetById(id)
// 	return gomvc.Xml(u)
// }

// //get http://localhost:8080/user/Slice/10
// //return slice
// func (*UserController) Slice(count int) []User {
// 	users := Take(count)
// 	return users
// }

// //get http://localhost:8080/user/search/123456_10_20
// //muti parameter
// func (*UserController) Search(zipcode string, ageFrom, ageTo int) []User {
// 	users := Search(zipcode, ageFrom, ageTo)
// 	return users
// }

// //get http://localhost:8080/user/struct?Id=1000&Name=hello&Zipcode=000000&Age=18
// //struct as parameter
// //validation(TODO):required,pattern,type,max,min,range,maxLength,minLength, rangeLength,
// //  number,date,time,zipcode,alphanumeric,lettersonly,email,url,greaterThan,lessThan
// func (*UserController) Struct(p struct {
// 	Id      int    `required, min:0`
// 	Name    string `required, rangeLength:1-10`
// 	Age     int    `default:1, range:0-99`
// 	Zipcode string `pattern:[0-9]+`
// },) User {
// 	return User{Id: p.Id, Name: p.Name, Age: p.Age, ZipCode: p.Zipcode}
// }

// //post http://localhost:8080/user/post
// //unmarshal parameter from json
// func (*UserController) Post(u User) string {
// 	return u.String()
// }

// //GET http://localhost:8080/user/content
// //access http context
// func (*UserController) Content(ctx *gomvc.HttpContext) string {
// 	return fmt.Sprintln("request", ctx.RequestPath)
// }

// //PUT http://localhost:8080/user/Put/1000
// //put all them together
// func (*UserController) Put(id int, u User, ctx *gomvc.HttpContext) string {
// 	return fmt.Sprintf("%v %v %v", ctx.Method, id, u.String())
// }

// //http://localhost:8080/user/error
// //raise error
// func (*UserController) Error() string {
// 	n := 0
// 	i := 100 / n
// 	return strconv.Itoa(i)
// }
