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

// get: /doc/home
func (uc *DocController) Home(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/home.html"), nil
}
