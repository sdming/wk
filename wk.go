// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"errors"
	"io"
	"log"
	"net/http"
	"reflect"
)

const (
	HttpVerbsGet     = "GET"
	HttpVerbsPost    = "POST"
	HttpVerbsPut     = "PUT"
	HttpVerbsDelete  = "DELETE"
	HttpVerbsHead    = "HEAD"
	HttpVerbsTrace   = "TRACE"
	HttpVerbsConnect = "CONNECT"
	HttpVerbsOptions = "OPTIONS"
)

const (
	HeaderAccept          = "Accept"
	HeaderAcceptCharset   = "Accept-Charset"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderCacheControl    = "Cache-Control"
	HeaderContentEncoding = "Content-Encoding"
	HeaderContentLength   = "Content-Length"
	HeaderContentType     = "Content-Type"
	HeaderDate            = "Date"
	HeaderEtag            = "Etag"
	HeaderExpires         = "Expires"
	HeaderLastModified    = "Last-Modified"
	HeaderLocation        = "Location"
	HeaderPragma          = "Pragma"
	HeaderServer          = "Server"
	HeaderSetCookie       = "Set-Cookie"
	HeaderUserAgent       = "User-Agent"
)

const (
	ContentTypeStream     = "application/octet-stream"
	ContentTypeJson       = "application/json"
	ContentTypeJsonp      = "application/jsonp"
	ContentTypeJavascript = "application/javascript"
	ContentTypeHTML       = "text/html"
	ContentTypeXml        = "text/xml"
	ContentTypeCss        = "text/css"
	ContentTypeText       = "text/plain"
	ContentTypeGif        = "image/gif"
	ContentTypeIcon       = "image/x-icon"
	ContentTypeJpeg       = "image/jpeg"
	ContentTypePng        = "image/png"
)

//application/x-www-form-urlencoded
//multipart/form-data

var (
	msgServerTimeout     = "server timeout"
	msgServerInternalErr = http.StatusText(http.StatusInternalServerError)
	msgNotFound          = http.StatusText(http.StatusNotFound)
	msgNoResult          = "no result"
	msgNoView            = "view not found"
	msgNoAction          = "can not find action"
)

const (
	codeServerInternaError = http.StatusInternalServerError
)

const (
	LogError = iota
	LogInfo
	LogDebug
)

const (
	_root           = "/"
	_any            = "*"
	_route          = "_route"
	_static         = "_static"
	_render         = "_render"
	_action         = "action"
	_notFoundAction = "noaction"
	_defaultAction  = "default"
	_serverName     = "go web server "
	_version        = "0.4"
)

const (
	_wkWebServer             = "_webserver"
	_eventStartRequest       = "start_request" //request start
	_eventEndRequest         = "end_request"   //request end
	_eventStartExecute       = "start_execute" //processor start execute
	_eventEndExecute         = "end_execute"   //processor end
	_eventStartResultExecute = "start_result"  //start to execute result
	_eventEndResultExecute   = "end_result"    //result execute end
	_eventStartAction        = "start_action"  //
	_eventEndAction          = "end_action"
)

const (
	_defaultSize = 61
)

var (
	// can not find view
	errNoView = errors.New(msgNoView)

	// internal error
	errInternalError = errors.New(http.StatusText(http.StatusInternalServerError))

	// can not find actionin method
	errNoAction = errors.New(msgNoAction)

	// httpresult is nil
	errNoResult = errors.New(msgNoResult)
)

var (
	// action did not return a result
	resultVoid = &VoidResult{}

	// httpresult is nil
	resultNotFound = &NotFoundResult{}
)

var (
	// logger
	Logger *log.Logger

	// LogLevel is level of log
	LogLevel int = LogError

	// EnableProfile mean enable http profile or not
	EnableProfile bool = false
)

var (
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
)

type Handler interface {
	Execute(ctx *HttpContext)
}

type ViewData map[string]interface{}

type ContentTyper interface {
	Type() string
}

type Render interface {
	ContentTyper

	// Write write header & body
	Write(header http.Header, body io.Writer) error
}

// HttpResult is a interface that define how to write server reply to response
type HttpResult interface {
	// Execute
	Execute(ctx *HttpContext) error
}

// type RuntimeError struct {
// 	Err    error
// 	Stack  string
// 	Target string
// }

// func (re *RuntimeError) Error() string {
// 	return re.Err.Error()
// }

// func (re *RuntimeError) String() string {
// 	log.Println(fmt.Sprintf("error=%v;target=%s;statck=%s", re.Err, re.Target, re.Stack))
// 	return fmt.Sprintf("error=%v;target=%s;statck=%s", re.Err, re.Target, re.Stack)
// }
