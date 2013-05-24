// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"github.com/sdming/kiss/gotype"
	"reflect"
	"strings"
)

// ActionSubscriber is interface that subscrib controller events
type ActionSubscriber interface {
	OnActionExecuting(action *ActionContext)
	OnActionExecuted(action *ActionContext)
	OnException(action *ActionContext)
}

// ActionContext is data that pass to ActionSubscriber
type ActionContext struct {
	// Name is action method name
	Name       string
	Context    *HttpContext
	Controller reflect.Value
	Result     HttpResult
	Err        error
}

// Controller is wrap of controller handler
type Controller struct {

	// Handler
	Handler reflect.Value

	//Methods is cache of method reflect info
	Methods map[string]*gotype.MethodInfo

	// actionSubscriber
	actionSubscriber ActionSubscriber

	//server
	server *HttpServer
}

// newController
func newController(x interface{}) (c *Controller, err error) {

	handler := reflect.ValueOf(x)
	handlerType := handler.Type()

	c = &Controller{Handler: handler, Methods: make(map[string]*gotype.MethodInfo)}
	c.actionSubscriber, _ = x.(ActionSubscriber)

	for i := 0; i < handler.NumMethod(); i++ {
		method := handlerType.Method(i)
		methodValue := handler.Method(i)
		if !methodValue.IsValid() || methodValue.Kind() != reflect.Func || method.Name == "" {
			continue
		}

		// not export
		if method.PkgPath != "" {
			continue
		}

		methdoInfo := gotype.GetMethodInfo(method)
		c.Methods[strings.ToLower(methdoInfo.Method.Name)] = methdoInfo

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

	actionName = actionName + strings.ToLower(ctx.Method)
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

// Execute call action method
func (c *Controller) Execute(ctx *HttpContext) {
	ctx.Result, ctx.Error = c.invoke(ctx)
}

// invoke
func (c *Controller) invoke(ctx *HttpContext) (result HttpResult, err error) {
	method, err := c.findAction(ctx)
	if err != nil {
		return nil, err
	}

	if LogLevel >= LogDebug {
		Logger.Println("controller dispatch request to action", method.Method.Name)
	}

	var actionContext *ActionContext
	actionContext = &ActionContext{
		Controller: c.Handler,
		Name:       method.Method.Name,
		Context:    ctx,
		Err:        nil,
	}

	c.server.Fire(_route, _eventStartAction, c, actionContext, ctx)

	if c.actionSubscriber != nil {
		c.actionSubscriber.OnActionExecuting(actionContext)
	}

	//var handle func(*HttpContext) (HttpResult, error)
	in := make([]reflect.Value, 2)
	var out []reflect.Value
	in[0] = c.Handler
	in[1] = reflect.ValueOf(ctx)
	out, err = safeCall(method.Func, in)

	if err == nil {
		if out[0].IsNil() {
			result = nil
		} else {
			result = out[0].Interface().(HttpResult)
		}
		if !out[1].IsNil() {
			err = out[1].Interface().(error)
		} else {
			err = nil
		}
	} else {
		result = nil
	}

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
