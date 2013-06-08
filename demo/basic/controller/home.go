// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic basic

*/
package controller

import (
	"errors"
	"github.com/sdming/wk"
)

type HomeController struct {
}

var home *HomeController

func RegisterHomeRoute(server *wk.HttpServer) {
	home = &HomeController{}

	server.RouteTable.Get("/?").To(HomeIndex)
	server.RouteTable.Get("/about?").To(HomeAbout)
	server.RouteTable.Get("/error?").To(HomeError)
	server.RouteTable.Get("/panic?").To(HomePanic)
}

func HomeIndex(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/home.html"), nil
}

func HomeAbout(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/about.html"), nil
}

func HomeError(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return nil, errors.New("error happend")
}

func HomePanic(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	panic("panic happend")
	return nil, nil
}
