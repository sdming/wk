// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// RouteData is wrap of route data
type RouteData map[string]string

// Int parse route value as int
func (r RouteData) Int(name string) (int, bool) {
	if s, ok := r[name]; ok {
		if i, err := strconv.Atoi(s); err == nil {
			return i, true
		}
	}
	return 0, false
}

// Int return route value as int or v
func (r RouteData) IntOr(name string, v int) int {
	if s, ok := r[name]; ok {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}
	return v
}

// Bool parse route value as bool
func (r RouteData) Bool(name string) (bool, bool) {
	if s, ok := r[name]; ok {
		if b, err := strconv.ParseBool(s); err == nil {
			return b, true
		}
	}
	return false, false
}

// BoolOr return route value as bool or v
func (r RouteData) BoolOr(name string, v bool) bool {
	if s, ok := r[name]; ok {
		if b, err := strconv.ParseBool(s); err == nil {
			return b
		}
	}
	return v
}

// Float parse route value as float64
func (r RouteData) Float(name string) (float64, bool) {
	if s, ok := r[name]; ok {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

// FloatOr return route value as float64 or v
func (r RouteData) FloatOr(name string, v float64) float64 {
	if s, ok := r[name]; ok {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}
	return v
}

// Str return route value
func (r RouteData) Str(name string) (string, bool) {
	if s, ok := r[name]; ok {
		return s, true
	}
	return "", false
}

// Str return route value or v if name donesn't exist
func (r RouteData) StrOr(name string, v string) string {
	if s, ok := r[name]; ok {
		return s
	}
	return v
}

// HttpContext is wrap of request & response
type HttpContext struct {
	// Server is current http server
	Server *HttpServer

	// Request is *http.Request
	Request *http.Request

	// Resonse is http.ResponseWriter
	Resonse http.ResponseWriter

	// Method is http method
	Method string

	// RequestPath is
	RequestPath string

	// PhysicalPath is path of static file
	PhysicalPath string

	// RouteData
	RouteData RouteData

	// ViewData
	ViewData map[string]interface{}

	// Result
	Result HttpResult

	// Error
	Error error

	// Flash is flash variables of request life cycle
	Flash map[string]interface{}

	// Session
	Session Session

	// SessionIsNew is true if create session in this request
	SessionIsNew bool
}

// String
func (ctx *HttpContext) String() string {
	return fmt.Sprintf("%s %s %v %v \n;", ctx.Method, ctx.RequestPath, ctx.Result, ctx.Error)
}

// RouteValue return route value by name
func (ctx *HttpContext) RouteValue(name string) (string, bool) {
	v, ok := ctx.RouteData[name]
	return v, ok
}

// FV is alias of FormValue
func (ctx *HttpContext) FV(name string) string {
	return ctx.Request.FormValue(name)
}

// FormValue is alias of Request FormValue
func (ctx *HttpContext) FormValue(name string) string {
	return ctx.Request.FormValue(name)
}

// FormInt parse form value as int
func (ctx *HttpContext) FormInt(name string) (int, bool) {
	if s := ctx.FormValue(name); s != "" {
		if i, err := strconv.Atoi(s); err == nil {
			return i, true
		}
	}
	return 0, false
}

// FormIntOr return form value as int or v
func (ctx *HttpContext) FormIntOr(name string, v int) int {
	if s := ctx.FormValue(name); s != "" {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}
	return v
}

// FormBool parse form value as bool
func (ctx *HttpContext) FormBool(name string) (bool, bool) {
	if s := ctx.FormValue(name); s != "" {
		if b, err := strconv.ParseBool(s); err == nil {
			return b, true
		}
	}
	return false, false
}

// FormBoolOr return form value as bool or v
func (ctx *HttpContext) FormBoolOr(name string, v bool) bool {
	if s := ctx.FormValue(name); s != "" {
		if b, err := strconv.ParseBool(s); err == nil {
			return b
		}
	}
	return v
}

// FormFloat parse form value as float64
func (ctx *HttpContext) FormFloat(name string) (float64, bool) {
	if s := ctx.FormValue(name); s != "" {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

// FormFloatOr return form value as float64 or v
func (ctx *HttpContext) FormFloatOr(name string, v float64) float64 {
	if s := ctx.FormValue(name); s != "" {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}
	return v
}

// // QueryValue return value from request URL query
// func (ctx *HttpContext) QueryValue(name string) []string {
// 	return ctx.Request.URL.Query()[name]
// }

// ReqHeader return request header by name
func (ctx *HttpContext) ReqHeader(name string) string {
	return ctx.Request.Header.Get(name)
}

// ResHeader return response header by name
func (ctx *HttpContext) ResHeader(name string) string {
	return ctx.Resonse.Header().Get(name)
}

// SetHeader set resonse http header
func (ctx *HttpContext) SetHeader(key string, value string) {
	ctx.Resonse.Header().Set(key, value)
}

// AddHeader add response http header
func (ctx *HttpContext) AddHeader(key string, value string) {
	ctx.Resonse.Header().Add(key, value)
}

// ContentType set response header Content-Type
func (ctx *HttpContext) ContentType(ctype string) {
	ctx.Resonse.Header().Set(HeaderContentType, ctype)
}

// Status write status code to http header
func (ctx *HttpContext) Status(code int) {
	ctx.Resonse.WriteHeader(code)
}

// Accept return request header Accept
func (ctx *HttpContext) Accept() string {
	return ctx.Request.Header.Get(HeaderAccept)
}

// Write writes b to resposne
func (ctx *HttpContext) Write(b []byte) (int, error) {
	return ctx.Resonse.Write(b)
}

// Expires set reponse header Expires
func (ctx *HttpContext) Expires(t string) {
	ctx.SetHeader(HeaderExpires, t)
}

// SetCookie set cookie to response
func (ctx *HttpContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.Resonse, cookie)
}

// Cookie return cookie from request
func (ctx *HttpContext) Cookie(name string) (*http.Cookie, error) {
	return ctx.Request.Cookie(name)
}

// SessionId return sessionid of current context
func (ctx *HttpContext) SessionId() string {
	return string(ctx.Session)
}

// Flush flush response immediately
func (ctx *HttpContext) Flush() {
	// _, buf, _ := ctx.Resonse.(http.Hijacker).Hijack()
	// if buf != nil {
	// 	buf.Flush()
	// }
	f, ok := ctx.Resonse.(http.Flusher)
	if ok && f != nil {
		f.Flush()
	}
}

// GetFlash return value in Context.Flash
func (ctx *HttpContext) GetFlash(key string) (v interface{}, ok bool) {
	if ctx.Flash == nil {
		return nil, false
	}
	v, ok = ctx.Flash[key]
	return
}

// SetFlash set value to Context.Flash
func (ctx *HttpContext) SetFlash(key string, v interface{}) {
	if ctx.Flash == nil {
		ctx.Flash = make(map[string]interface{})
	}
	ctx.Flash[key] = v
}

// ReadBody read Request.Body
func (ctx *HttpContext) ReadBody() ([]byte, error) {
	return ioutil.ReadAll(ctx.Request.Body)
}
