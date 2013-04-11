// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	//"github.com/sdming/kiss"
	"bytes"
	"errors"
	"fmt"
	"github.com/sdming/kiss/gotype"
	"reflect"
	"runtime/debug"
	"strings"
)

// ActionSubscriber  
type ActionSubscriber interface {
	OnActionExecuting(action *ActionContext)
	OnActionExecuted(action *ActionContext)
	OnException(action *ActionContext)
}

type ActionContext struct {
	Name       string
	Context    *HttpContext
	Controller reflect.Value
	Result     HttpResult
	Err        error
}

// controller
type Controller struct {

	// Handler
	Handler reflect.Value

	//Methods cache method reflect info 
	Methods map[string]*gotype.MethodInfo

	//Binder Binder

	// ActionSubscriber 
	actionSubscriber ActionSubscriber

	//server
	server *HttpServer
}

func newController(x interface{}) (c *Controller, err error) {

	var handler reflect.Value
	if v, ok := x.(reflect.Value); ok {
		handler = v
	} else {
		handler = reflect.ValueOf(x)
	}

	c = &Controller{Handler: handler, Methods: make(map[string]*gotype.MethodInfo)}
	if subscriber, ok := x.(ActionSubscriber); ok {
		c.actionSubscriber = subscriber
	}

	for i := 0; i < handler.NumMethod(); i++ {
		m := handler.Method(i)
		typ := m.Type()
		fmt.Println("method", m.Kind(), m, typ.PkgPath())

		if !m.IsValid() || m.Kind() != reflect.Func || typ.Name() == "" {
			continue
		}

		// not export
		if typ.PkgPath() != "" {
			continue
		}

		c.Methods[typ.Name()] = gotype.GetMethodInfo(m)
	}
	return c, nil
}

// findAction
func (c *Controller) findAction(ctx *HttpContext) (method *gotype.MethodInfo, err error) {

	actionName, ok := ctx.RouteData[_action]
	if !ok || actionName == "" {
		switch ctx.Method {
		case HttpVerbsConnect, HttpVerbsOptions, HttpVerbsTrace:
			actionName = "_NotSupport" + ctx.Method
		case HttpVerbsGet, HttpVerbsDelete, HttpVerbsPost:
			actionName = ctx.Method
		case HttpVerbsHead:
			actionName = "Get"
		}
	}
	actionName = strings.ToLower(actionName)

	if method, ok = c.Methods[actionName]; ok {
		return
	}

	actionName = _defaultAction
	if method, ok = c.Methods[actionName]; ok {
		return
	}

	err = errNoAction
	return
}

func (c *Controller) Execute(ctx *HttpContext) {
	result, err := c.invoke(ctx)
	if err != nil {
		ctx.Error = err
		return
	}
	ctx.Result = result
}

func (c *Controller) invoke(ctx *HttpContext) (result HttpResult, err error) {
	fmt.Println("invoke")
	method, err := c.findAction(ctx)
	fmt.Println("findAction", method, err)

	if err != nil {
		return nil, err
	}

	if LogLevel >= LogDebug {
		Logger.Println("controller dispatch request to action", method)
	}

	var actionContext *ActionContext
	actionContext = &ActionContext{
		Controller: c.Handler,
		Name:       method.Name,
		Context:    ctx,
		Err:        nil,
	}

	c.server.Fire(_route, _eventStartAction, c, actionContext, ctx)

	if c.actionSubscriber != nil {
		c.actionSubscriber.OnActionExecuting(actionContext)
	}

	var handle func(*HttpContext) (HttpResult, error)
	result, err = safeCallAction(ctx, handle)

	actionContext.Err = err
	actionContext.Result = result

	if err != nil && c.actionSubscriber != nil {
		c.actionSubscriber.OnException(actionContext)
	}

	if c.actionSubscriber != nil {
		c.actionSubscriber.OnActionExecuted(actionContext)
	}

	c.server.Fire(_route, _eventEndAction, c, actionContext, ctx)

	return result, err
}

func safeCallAction(ctx *HttpContext, handle func(*HttpContext) (HttpResult, error)) (result HttpResult, err error) {
	defer func() {
		if x := recover(); x != nil {
			if LogLevel >= LogError {
				var buf bytes.Buffer
				fmt.Fprintf(&buf, "call action panic %v : %v \n", ctx.Request.URL, x)
				buf.Write(debug.Stack())
				Logger.Println(buf.String())
			}

			if e, ok := x.(error); ok {
				err = e
			} else {
				err = errors.New(fmt.Sprintln("call action fail", x))
			}
		}

	}()

	return handle(ctx)
}
