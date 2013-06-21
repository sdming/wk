// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"log"
	"net"
	"net/http"
	//_ "net/http/pprof"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"
)

// HttpServer
type HttpServer struct {
	// config
	Config *WebConfig

	// net.Listener
	Listener net.Listener

	// // *ServeMux
	// Mux *http.ServeMux

	// http server
	server *http.Server

	// Processes
	Processes ProcessTable

	// RouteTable
	RouteTable *RouteTable

	//server variables
	Variables map[string]interface{}

	// ViewEngine
	ViewEngine ViewEngine
}

// Fire can fire a event
func (srv *HttpServer) Fire(moudle, name string, source, data interface{}, context *HttpContext) {

	var e *EventContext

	for _, sub := range Subscribers {
		if (sub.Moudle == _any || sub.Moudle == moudle) && (sub.Name == _any || sub.Name == name) {
			if e == nil {
				e = &EventContext{
					Moudle:  moudle,
					Name:    name,
					Source:  source,
					Data:    data,
					Context: context,
				}
			}

			sub.Handler.On(e)
		}
	}
}

// DefaultServer return a http server with default config
func NewDefaultServer() (srv *HttpServer, err error) {
	var conf *WebConfig

	conf, err = ReadDefaultConfigFile()
	if err != nil {
		conf = NewDefaultConfig()
	}

	return NewHttpServer(conf)
}

// NewHttpServer return a http server with provided config
func NewHttpServer(config *WebConfig) (srv *HttpServer, err error) {
	srv = &HttpServer{
		Config: config,
	}
	srv.init()
	return srv, nil
}

// init
func (srv *HttpServer) init() error {
	srv.Variables = make(map[string]interface{})
	srv.RouteTable = newRouteTable()

	// copy hander, maybe does not need this?
	l := len(Processes)
	srv.Processes = make([]*Process, l)
	for i := 0; i < l; i++ {
		Processes[i].Handler.Register(srv)
		srv.Processes[i] = Processes[i]
	}

	return nil
}

// listenAndServe
func (srv *HttpServer) listenAndServe() (err error) {
	srv.server = &http.Server{
		Addr:           srv.Config.Address,
		Handler:        srv,
		ReadTimeout:    time.Duration(srv.Config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(srv.Config.WriteTimeout) * time.Second,
		MaxHeaderBytes: srv.Config.MaxHeaderBytes,
	}

	return srv.server.ListenAndServe()
}

// error return error message to client
func (srv *HttpServer) error(ctx *HttpContext, message string, code int) {
	if LogLevel >= LogError {
		Logger.Println(message)
		Logger.Println(string(debug.Stack()))
	}

	if code == 0 {
		code = http.StatusInternalServerError
	}

	if message == "" {
		message = http.StatusText(code)
	}

	http.Error(ctx.Resonse, message, code)
}

// Setup initialize server instance
func (srv *HttpServer) Setup() (err error) {
	if Logger == nil {
		Logger = log.New(os.Stdout, _serverName+_version, log.Ldate|log.Ltime)
	}
	Logger.Println("http server is starting")

	Logger.Println("Address:", srv.Config.Address,
		"\n\t RootDir:", srv.Config.RootDir,
		"\n\t ConfigDir:", srv.Config.ConfigDir,
		"\n\t PublicDir:", srv.Config.PublicDir,
		"\n\t Debug:", srv.Config.Debug,
	)

	if err = srv.configSession(); err != nil {
		return
	}

	if srv.Config.ViewEnable {
		if err = srv.configViewEngine(); err != nil {
			return
		}
	}

	if len(srv.Processes) == 0 {
		return errors.New("server processes is empty")
	}

	for _, p := range srv.Processes {
		Logger.Println("process", p.Method, p.Path, p.Name)
	}

	return nil
}

// Start start server instance and listent request
func (srv *HttpServer) Start() (err error) {
	if err = srv.Setup(); err != nil {
		Logger.Println("http server setup fail:", err)
		return
	}

	if err = srv.listenAndServe(); err != nil {
		Logger.Println("http server listen fail:", err)
		return
	}

	Logger.Println("http server server is listen on:", srv.Config.Address)
	return nil
}

// // Close stop listen and close http server
// func (s *HttpServer) Close() error {
// 	Logger.Println("http server is closing")

// 	if srv.Listener != nil {
// 		srv.Listener.Close()
// 	}
// 	Logger.Println("http server closed")

// 	return nil
// }

// ServeHTTP
func (srv *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer func() {
		err := recover()
		if err == nil {
			return
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "http server panic %v : %v\n", r.URL, err)
		buf.Write(debug.Stack())
		Logger.Println(buf.String())

		http.Error(w, msgServerInternalErr, codeServerInternaError)

	}()

	ctx := srv.buildContext(w, r)

	if LogLevel >= LogDebug {
		Logger.Println("request start", ctx.Method, ctx.RequestPath)
	}
	srv.Fire(_wkWebServer, _eventStartRequest, srv, nil, ctx)

	srv.doServer(ctx)

	srv.Fire(_wkWebServer, _eventEndRequest, srv, nil, ctx)
	if LogLevel >= LogDebug {
		Logger.Println("request end", ctx.Method, ctx.RequestPath, ctx.Result, ctx.Error)
	}

}

// doServer
func (srv *HttpServer) doServer(ctx *HttpContext) {
	ctx.SetHeader(HeaderServer, _serverName)
	srv.exeSession(ctx)

	for _, h := range srv.Processes {
		if h.match(ctx) {
			srv.exeProcess(ctx, h)
		}
	}
}

// buildContext
func (s *HttpServer) buildContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	_ = r.ParseForm()
	ctx := &HttpContext{
		Resonse:     w,
		Request:     r,
		Method:      r.Method,
		RequestPath: cleanPath(strings.TrimSpace(r.URL.Path)),
		Server:      s,
	}

	if s.Config.ViewEnable {
		ctx.ViewData = make(map[string]interface{})
	}

	return ctx
}

// exeProcess execute all HttpProcessor
func (srv *HttpServer) exeProcess(ctx *HttpContext, p *Process) (err error) {

	if LogLevel >= LogDebug {
		Logger.Println("process start:", p.Name)
	}

	defer func() {
		if x := recover(); x != nil {
			if LogLevel >= LogError {
				Logger.Println("execute process recover:", p.Name, x)
				Logger.Println(string(debug.Stack()))
			}

			if e, ok := x.(error); ok {
				err = e
			} else {
				err = errors.New(fmt.Sprintln(x))
			}
			ctx.Error = err
		}
	}()

	srv.Fire(p.Name, _eventStartExecute, p, nil, ctx)

	p.Handler.Execute(ctx)

	srv.Fire(p.Name, _eventEndExecute, p, nil, ctx)

	if LogLevel >= LogDebug {
		Logger.Println("process end", p.Name, err)
	}

	return nil
}

// MapPath return physical path
func (srv *HttpServer) MapPath(file string) string {
	return path.Join(srv.Config.RootDir, file)
	//info, err := os.Stat(f)
	// if err != nil || info.IsDir() {
	// 	return ""
	// }
	// return f
}

// // serverFile
// func (srv *HttpServer) serverFile(ctx *HttpContext, file string, status int) error {
// 	fullPath := path.Join(srv.Config.RootDir, file)
// 	_, err := os.Stat(fullPath)
// 	if err != nil {
// 		return err
// 	}
// 	var f *os.File
// 	f, err = os.Open(fullPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	ctx.ContentType(mime.TypeByExtension(filepath.Ext(fullPath)))
// 	ctx.Status(status)
// 	//io.Copy(dst·3, src·4)
// 	//http.ServeContent(ctx.Resonse, ctx.Request, fullPath, info.ModTime(), f)
// 	http.ServeFile(ctx.Resonse, ctx.Request, fullPath)
// 	return nil

// }
