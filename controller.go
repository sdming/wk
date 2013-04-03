// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	//"github.com/sdming/kiss"
	//"github.com/sdming/kiss/gotype"
	"reflect"
	"strings"
)

// ActionSubscriber  
type ActionSubscriber interface {
	OnActionExecuting(action *ActionContext)
	OnActionExecuted(action *ActionContext)
	OnException(action *ActionContext)
}

type ActionContext struct {
	// Name is action name
	Name string

	Context *HttpContext

	// Err is not nil if action return error or panic
	Err error
}

// controller
type controller struct {

	// Handler
	Handler reflect.Value

	//Methods cache method reflect info 
	Methods map[string]methodInfo

	Binder Binder

	// ActionSubscriber 
	actionSubscriber ActionSubscriber
}

type methodInfo struct {
	// name 
	Name string

	// Type
	Type reflect.Type

	// Value
	Value reflect.Value

	// NumIn is number of input paramter
	NumIn int

	// NumOut is number of output paramter
	NumOut int

	// Ins is input parameter type
	In []reflect.Type

	// Outs is output parameter type
	Out []reflect.Type

	outError      bool
	outErrorIndex int
}

func newController(x interface{}) (c controller, err error) {

	var handler reflect.Value
	if v, ok := x.(reflect.Value); ok {
		handler = v
	} else {
		handler = reflect.ValueOf(x)
	}

	c = controller{Handler: handler, Methods: make(map[string]methodInfo)}
	if subscriber, ok := x.(ActionSubscriber); ok {
		c.actionSubscriber = subscriber
	}

	for i := 0; i < handler.NumMethod(); i++ {
		m := handler.Method(i)
		typ := m.Type()
		if !m.IsValid() || m.Kind() != reflect.Func || typ.Name() == "" {
			continue
		}

		// not export
		if typ.PkgPath() != "" {
			continue
		}

		info := methodInfo{
			Name:   strings.ToLower(typ.Name()),
			Type:   typ,
			Value:  m,
			NumIn:  typ.NumIn(),
			NumOut: typ.NumOut(),
			In:     make([]reflect.Type, typ.NumIn()),
			Out:    make([]reflect.Type, typ.NumOut())}

		for j := 0; j < typ.NumIn(); j++ {
			info.In[j] = typ.In(j)
		}
		for j := 0; j < typ.NumOut(); j++ {
			info.Out[j] = typ.Out(j)

			if info.Out[j].Implements(reflect.TypeOf(err)) {
				info.outError = true
				info.outErrorIndex = j
			}
		}
		c.Methods[typ.Name()] = info
	}
	return c, nil
}

// findAction
func (c *controller) findAction(ctx *HttpContext) (method methodInfo, err error) {

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

	err = errNoaction
	return
}

func (c *controller) Handle(ctx *HttpContext) {
	result, err := c.invoke(ctx)
	if err != nil {
		ctx.Error = err
		return
	}
	ctx.Result = result
}

func (c *controller) invoke(ctx *HttpContext) (result HttpResult, err error) {
	method, err := c.findAction(ctx)
	if err != nil {
		return nil, err
	}

	if LogLevel >= LogDebug {
		Logger.Println("controller dispatch request to action", method.Name)
	}

	args, err := decodeArgs(ctx, method)
	if err != nil {
		return nil, err
	}

	var actionContext *ActionContext
	if c.actionSubscriber != nil {
		actionContext = &ActionContext{
			Name:    method.Name,
			Context: ctx,
			Err:     nil,
		}
		c.actionSubscriber.OnActionExecuting(actionContext)
	}

	ret, err := safeCall(method.Value, args)

	if c.actionSubscriber != nil {
		actionContext.Err = err
		c.actionSubscriber.OnActionExecuting(actionContext)
	}

	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return resultVoid, nil
	}

	result = convertResult(ret[0])
	return result, nil
}

// convertResult convert x to http result
func convertResult(reflect.Value) HttpResult {
	return nil
	// v := reflect.Indirect(reflect.ValueOf(i))
	// kind := v.Kind()

	// switch {
	// case gotype.IsSimple(kind):
	// 	return &ContentResult{Data: v.Interface()}
	// case gotype.IsStruct(kind) || gotype.IsCollect(kind):
	// 	switch {
	// 	case strings.Index(accept, ContentTypeJson) > -1:
	// 		return &JsonResult{Data: v.Interface()}
	// 	case strings.Index(accept, ContentTypeXml) > -1:
	// 		return &XmlResult{Data: v.Interface()}
	// 	case strings.Index(accept, ContentTypeJsonp) > -1:
	// 		return &JsonpResult{Data: v.Interface()}
	// 	case strings.Index(accept, ContentTypeJavascript) > -1:
	// 		return &JavaScriptResult{Data: v.Interface()}
	// 	}
	// }
	// return &ContentResult{Data: v.Interface()}
}

func decodeArgs(ctx *HttpContext, method methodInfo) (args []reflect.Value, err error) {
	return nil, nil
	// numIn := method.NumIn
	// args = make([]reflect.Value, numIn, numIn)

	// for i := 0; i < numIn; i++ {
	// 	in := method.In[i]
	// 	switch in.Kind() {
	// 	case reflect.Ptr:
	// 		if strings.Contains(in.Elem().Name(), "HttpContext") {
	// 			args[i] = reflect.ValueOf(ctx)
	// 		} else {
	// 			return nil, errors.New("TODO: prase in parameter " + in.Kind().String())
	// 		}
	// 	case reflect.Struct:
	// 		v := reflect.New(in)
	// 		if ctx.Method == HttpVerbsPost || ctx.Method == HttpVerbsPut {
	// 			body, err := ioutil.ReadAll(ctx.Request.Body)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			err = json.Unmarshal(body, v.Interface())
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 			v = reflect.Indirect(v)
	// 		} else {
	// 			v = reflect.Indirect(v)
	// 			var src kiss.StrGetFunc = func(name string) (string, bool) {
	// 				x := ctx.Value(name)
	// 				return x, x == ""
	// 			}
	// 			kiss.ParseStruct(v, src)
	// 		}
	// 		args[i] = v
	// 	default:
	// 		// simple type
	// 		name := "p" + strconv.Itoa(i)
	// 		str, ok := ctx.RouteData[name]
	// 		if !ok {
	// 			str = ctx.Request.FormValue(name)
	// 		}
	// 		if str == "" {
	// 			continue
	// 		}
	// 		if v, err := gotype.Atok(str, in.Kind()); err != nil {
	// 			return nil, err
	// 		} else {
	// 			args[i] = v
	// 		}
	// 	}
	// }

	// return args, nil
}
