// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"fmt"
	"github.com/sdming/pathexp"
	"reflect"
)

type Binder func(*HttpContext) []reflect.Value

// Route is collection of RouteRule
type RouteTable []*RouteRule

func (r RouteTable) addRouteRule(method, path string) *RouteRule {
	rule, err := NewRouteRule(method, path)
	if err != nil {
		panic(err)
	}

	r = append(r, rule)
	return rule
}

// Get match HttpVerbsGet
func (r RouteTable) Get(path string) *RouteRule {
	return r.addRouteRule(HttpVerbsGet, path)
}

// Put match HttpVerbsPut
func (r RouteTable) Put(path string) *RouteRule {
	return r.addRouteRule(HttpVerbsPut, path)
}

// Post match HttpVerbsPost
func (r RouteTable) Post(path string) *RouteRule {
	return r.addRouteRule(HttpVerbsPost, path)
}

// Delete match HttpVerbsDelete
func (r RouteTable) Delete(path string) *RouteRule {
	return r.addRouteRule(HttpVerbsDelete, path)
}

// Path match any http method
func (r RouteTable) Path(path string) *RouteRule {
	return r.addRouteRule(_any, path)
}

// Add a RouteRule
func (r RouteTable) Add(rule *RouteRule) error {
	exp, err := pathexp.Compile(rule.Pattern)
	if err != nil {
		return nil
	}
	rule.pathex = exp

	r = append(r, rule)
	return nil
}

// Match return matched rule and route data
func (r RouteTable) Match(ctx *HttpContext) (rule *RouteRule, data map[string]string, ok bool) {
	for _, rule = range r {
		if data, ok = rule.Match(ctx); ok {
			return
		}
	}
	return
}

func newRouteTable() RouteTable {
	return make([]*RouteRule, _defaultSize)
}

// RouteRule is rule of route 
type RouteRule struct {

	// Methos is http method of request
	Method string

	// Pattern is path pattern
	Pattern string

	// Handl process request
	Handle Handler

	// Binder create parameter from request
	Binder Binder

	// Formatter format result 
	Formatter interface{}

	pathex *pathexp.Pathex
}

func (r *RouteRule) String() string {
	return fmt.Sprintln(r.Method, r.Pattern, reflect.TypeOf(r.Handle))
}

// match 
func (r *RouteRule) Match(ctx *HttpContext) (data map[string]string, ok bool) {
	if ctx.Method != r.Method && r.Method != _any && r.Method != "" && !(ctx.Method == HttpVerbsHead && r.Method == HttpVerbsGet) {
		return
	}

	routeData := r.pathex.FindAllStringSubmatch(ctx.RequestPath)
	if routeData == nil {
		return
	}

	data = make(map[string]string)
	for _, x := range routeData {
		data[x[0]] = x[1]
	}

	ok = true
	return
}

func NewRouteRule(method, path string) (rule *RouteRule, err error) {
	var exp *pathexp.Pathex
	exp, err = pathexp.Compile(path)
	if err != nil {
		return
	}

	rule = &RouteRule{
		Pattern: path,
		Method:  method,
		pathex:  exp,
	}
	return
}

// To map path to handle
func (r *RouteRule) To(handle Handler) {
	r.Handle = handle
}

// // To map path to handle
// func (r *RouteRule) ToMethod(method reflect.MethodValue) {
// 	return
// }

// To map path to handle
func (r *RouteRule) ToFunc(function interface{}) {
	return
}

// To map path to Controller
func (r *RouteRule) ToController(controller interface{}) {
	return
}

// RouteProcessor 
type RouteProcessor struct {
	server *HttpServer
}

func newRouteProcessor() *RouteProcessor {
	return &RouteProcessor{}
}

func (r *RouteProcessor) Register(server *HttpServer) {
	r.server = server
}

// Execute 
func (r *RouteProcessor) Execute(ctx *HttpContext) {
	route, routeData, ok := r.server.RouteTable.Match(ctx)
	fmt.Println("route", ok, route, routeData)

	if !ok {
		return
	}

	ctx.RouteData = routeData
	route.Handle(ctx)
	ctx.Write([]byte("route" + route.Pattern))
}

/*
map
	1: pa,p2,p3
	2:["a","b","c"]

H(httpcontext) httpresult, error
*/
