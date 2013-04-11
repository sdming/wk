// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"errors"
	"fmt"
	"github.com/sdming/kiss/gotype"
	"github.com/sdming/pathexp"
	"reflect"
	"strconv"
	"strings"
)

// Route is collection of RouteRule
type RouteTable struct {
	Routes []*RouteRule
}

//type RouteTable  []*RouteRule

func (r *RouteTable) addRouteRule(method, path string) *RouteRule {
	rule, err := NewRouteRule(method, path)
	if err != nil {
		panic(err)
	}
	r.Routes = append(r.Routes, rule)
	return rule
}

// Get match HttpVerbsGet
func (r *RouteTable) Get(path string) *RouteRule {
	return r.addRouteRule(HttpVerbsGet, path)
}

// Put match HttpVerbsPut
func (r *RouteTable) Put(path string) *RouteRule {
	return r.addRouteRule(HttpVerbsPut, path)
}

// Post match HttpVerbsPost
func (r *RouteTable) Post(path string) *RouteRule {
	return r.addRouteRule(HttpVerbsPost, path)
}

// Delete match HttpVerbsDelete
func (r *RouteTable) Delete(path string) *RouteRule {
	return r.addRouteRule(HttpVerbsDelete, path)
}

// Path match any http method
func (r *RouteTable) Path(path string) *RouteRule {
	return r.addRouteRule(_any, path)
}

// Add a RouteRule
func (r *RouteTable) Add(rule *RouteRule) error {
	exp, err := pathexp.Compile(rule.Pattern)
	if err != nil {
		return nil
	}
	rule.pathex = exp

	r.Routes = append(r.Routes, rule)
	return nil
}

// Match return matched rule and route data
func (r *RouteTable) Match(ctx *HttpContext) (rule *RouteRule, data map[string]string, ok bool) {
	if r == nil || r.Routes == nil || len(r.Routes) == 0 {
		return
	}

	for _, rule = range r.Routes {
		//fmt.Println("test route match", rule)
		if data, ok = rule.Match(ctx); ok {
			//fmt.Println("test route match success", rule, ok)
			return
		}
	}
	return
}

func newRouteTable() *RouteTable {
	return &RouteTable{
		Routes: make([]*RouteRule, 0, 101),
	}
}

// RouteRule is rule of route 
type RouteRule struct {

	// Methos is http method of request
	Method string

	// Pattern is path pattern
	Pattern string

	// Handler process request
	Handler Handler

	pathex *pathexp.Pathex
}

func (r *RouteRule) String() string {
	return fmt.Sprint(r.Method, " ", r.Pattern, " handle by ", reflect.TypeOf(r.Handler))
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
func (r *RouteRule) HandleBy(handler Handler) {
	r.Handler = handler
}

// To map path to handle
func (r *RouteRule) To(handler func(*HttpContext) (HttpResult, error)) {
	r.Handler = RouteFunc(handler)
}

// // To map path to handle
// func (r *RouteRule) ToMethod(method reflect.MethodValue) {
// 	return
// }

// To map path to handle
func (r *RouteRule) ToFunc(function interface{}) *FuncServer {
	fv := reflect.ValueOf(function)
	handler := &FuncServer{
		Func:      function,
		funcValue: fv,
	}
	r.Handler = handler
	return handler
}

// To map path to Controller
func (r *RouteRule) ToController(controller interface{}) *Controller {
	handler, _ := newController(controller)
	r.Handler = handler
	return handler
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
	if ctx.Result != nil {
		return
	}

	route, routeData, ok := r.server.RouteTable.Match(ctx)
	if !ok {
		return
	}

	ctx.RouteData = routeData
	route.Handler.Execute(ctx)

	if ctx.Error != nil {
		if LogLevel >= LogError {
			Logger.Printf("route execute error: %s; request: %s", ctx.Error, ctx.Request.URL)
		}
	}

}

// FuncServer TODO:
type FuncServer struct {
	Binder    func(*HttpContext) ([]reflect.Value, error)
	Func      interface{}
	funcValue reflect.Value
	Formatter FormatFunc
}

// 
func (f *FuncServer) ReturnXml() *FuncServer {
	f.Formatter = formatXml
	return f
}

// 
func (f *FuncServer) ReturnJson() *FuncServer {
	f.Formatter = formatJson
	return f
}

// 
func (f *FuncServer) Return(fn FormatFunc) *FuncServer {
	f.Formatter = fn
	return f
}

// BindByIndex TODO:
func (f *FuncServer) BindByIndex() {
	binder := newIndexBinder(f.funcValue)

	f.Binder = func(ctx *HttpContext) ([]reflect.Value, error) {
		return binder.Bind(ctx)
	}
}

// BindByNames TODO:
func (f *FuncServer) BindByNames(name ...string) {
	binder := newNamedBinder(name[:], f.funcValue)

	f.Binder = func(ctx *HttpContext) ([]reflect.Value, error) {
		return binder.Bind(ctx)
	}
}

// BindToStruct TODO:
func (f *FuncServer) BindToStruct() {
	binder := newStructBinder(f.funcValue)

	f.Binder = func(ctx *HttpContext) ([]reflect.Value, error) {
		return binder.Bind(ctx)
	}
}

type structBinder struct {
	method *gotype.MethodInfo
}

func newStructBinder(fv reflect.Value) *structBinder {
	method := gotype.GetMethodInfo(fv)

	return &structBinder{
		method: method,
	}
}

func (binder *structBinder) Bind(ctx *HttpContext) ([]reflect.Value, error) {
	numIn := binder.method.NumIn
	args := make([]reflect.Value, numIn, numIn)

	for i := 0; i < numIn; i++ {
		in := binder.method.In[i]
		args[i] = reflect.Zero(in)

		if in.Kind() == reflect.Struct {
			args[i] = reflect.New(in).Elem()
		} else if in.Kind() == reflect.Ptr && in.Elem().Kind() == reflect.Struct {
			args[i] = reflect.New(in.Elem())
		} else {
			continue
		}

		uType := gotype.UnderlyingType(in)
		uValue := gotype.Underlying(args[i])
		filedNum := uType.NumField()

		for f := 0; f < filedNum; f++ {
			field := uType.Field(f)
			fieldValue := uValue.Field(f)
			fieldType := field.Type

			if !fieldValue.CanSet() {
				continue
			}

			fieldKind := fieldType.Kind()
			name := strings.ToLower(field.Name)

			str, ok := ctx.RouteData[name]
			if !ok {
				str = ctx.Request.FormValue(name)
			}
			if str == "" {
				continue
			}

			if gotype.IsSimple(fieldKind) {
				gotype.Value(fieldValue).Parse(str)
			} else {
				//TODO
			}
		}
	}

	return args, nil
}

type namedBinder struct {
	method  *gotype.MethodInfo
	argsMap []string
}

func newNamedBinder(args []string, fv reflect.Value) *namedBinder {
	method := gotype.GetMethodInfo(fv)

	return &namedBinder{
		method:  method,
		argsMap: args,
	}
}

func (binder *namedBinder) Bind(ctx *HttpContext) ([]reflect.Value, error) {
	numIn := binder.method.NumIn

	if len(binder.argsMap) != numIn {
		return nil, errors.New("args length doesn not match")
	}

	args := make([]reflect.Value, numIn, numIn)

	for i := 0; i < numIn; i++ {
		in := binder.method.In[i]
		name := binder.argsMap[i]
		str, ok := ctx.RouteData[name]
		if !ok {
			str = ctx.Request.FormValue(name)
		}
		if str == "" {
			args[i] = reflect.Zero(in)
			continue
		}

		if v, err := gotype.Atok(str, in.Kind()); err != nil {
			return nil, err
		} else {
			args[i] = v
		}
	}
	return args, nil
}

type indexBinder struct {
	method *gotype.MethodInfo
}

func newIndexBinder(fv reflect.Value) *indexBinder {
	method := gotype.GetMethodInfo(fv)
	return &indexBinder{
		method: method,
	}
}

func (binder *indexBinder) Bind(ctx *HttpContext) ([]reflect.Value, error) {
	numIn := binder.method.NumIn
	args := make([]reflect.Value, numIn, numIn)

	for i := 0; i < numIn; i++ {
		in := binder.method.In[i]
		name := "p" + strconv.Itoa(i)

		str, ok := ctx.RouteData[name]
		if !ok {
			str = ctx.Request.FormValue(name)
		}

		if str == "" {
			args[i] = reflect.Zero(in)
			continue
		}
		if v, err := gotype.Atok(str, in.Kind()); err != nil {
			return nil, err
		} else {
			args[i] = v
		}
	}
	return args, nil
}

func (f *FuncServer) Execute(ctx *HttpContext) {
	if f.Binder == nil {
		f.BindByIndex()
	}
	args, err := f.Binder(ctx)
	if err != nil {
		ctx.Error = err
		return
	}

	result, err := safeCall(f.funcValue, args)
	if err != nil {
		ctx.Error = err
		return
	}

	if len(result) == 0 {
		ctx.Result = resultVoid
		return
	}

	if httpResult, ok := result[0].Interface().(HttpResult); ok {
		ctx.Result = httpResult
		return
	}

	if f.Formatter != nil {
		if formatted, ok := f.Formatter(ctx, result[0].Interface()); ok {
			ctx.Result = formatted
			return
		}
	}

	ctx.Result = convertResult(ctx, result[0])
}

type RouteFunc func(*HttpContext) (HttpResult, error)

func (f RouteFunc) Execute(ctx *HttpContext) {
	result, err := f(ctx)
	if err != nil {
		ctx.Error = err
		return
	}
	ctx.Result = result
}

/*
map
	1: pa,p2,p3
	2:["a","b","c"]

H(httpcontext) httpresult, error
*/
