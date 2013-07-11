// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic demo

*/
package controller

import (
	"github.com/sdming/wk"
)

func init() {
}

type DocController struct {
}

func NewDocController() *DocController {
	return &DocController{}
}

var docController *DocController

func RegisterDocRoute(srv *wk.HttpServer) {
	docController = NewDocController()

	// url: /doc/xxx
	// route to DocController
	srv.RouteTable.Path("/doc/{action}").ToController(docController)
}

// get: /doc/index
func (uc *DocController) Index(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/index.html"), nil
}

// get: /doc/start
func (uc *DocController) Start(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/start.html"), nil
}

// get: /doc/document
func (uc *DocController) Document(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/document.html"), nil
}

// get: /doc/demo
func (uc *DocController) Demo(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/demo.html"), nil
}

// get: /doc/sessiondemo
func (uc *DocController) SessionDemo(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/session_demo.html"), nil
}

// get: /doc/configdemo
func (uc *DocController) ConfigDemo(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/config_demo.html"), nil
}

// get: /doc/otherdemo
func (uc *DocController) OtherDemo(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/other_demo.html"), nil
}

// get: /doc/basicdemo
func (uc *DocController) BasicDemo(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/basic_demo.html"), nil
}

// get: /doc/routedemo
func (uc *DocController) RouteDemo(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/route_demo.html"), nil
}

// get: /doc/home
func (uc *DocController) Home(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/home.html"), nil
}
