// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic basic

*/
package controller

import (
	"github.com/sdming/wk"
)

type HomeController struct {
}

var home *HomeController

func RegisterHomeRoute(server *wk.HttpServer) {
	home = &HomeController{}

	server.RouteTable.Get("/?").To(HomeIndex)
	server.RouteTable.Get("/about?").To(HomeAbout)
}

func HomeIndex(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/home.html"), nil
}

func HomeAbout(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/about.html"), nil
}
