// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"net/http"
	"strings"
)

// NotFoundResult is wrap of http 404
type NotFoundResult struct {
}

// Execute return 404 to client
func (r *NotFoundResult) Execute(ctx *HttpContext) error {
	/*
		text/plain; text/html; text/*; 
		text/html,application/xhtml+xml,application/xml;
	*/
	if ctx.Server.Config.NotFoundPageEnable {
		accetps := ctx.Accept()
		if strings.Contains(accetps, "text/html") {
			if f := ctx.Server.MapPath("public/404.html"); isFileExists(f) {
				http.ServeFile(ctx.Resonse, ctx.Request, f)
				return nil
			}
			if ctx.Server.Config.ViewEnable {
				if f := ctx.Server.MapPath("views/404.html"); isFileExists(f) {
					ctx.ViewData["ctx"] = ctx
					return executeViewFile("404.html", ctx)
				}
			}
		}
		if strings.Contains(accetps, "text/plain") {
			if f := ctx.Server.MapPath("public/404.txt"); isFileExists(f) {
				http.ServeFile(ctx.Resonse, ctx.Request, f)
				return nil
			}
		}
	}

	http.Error(ctx.Resonse, msgNotFound, http.StatusNotFound)
	return nil
}
