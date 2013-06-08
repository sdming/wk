// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"fmt"
	"net/http"
	"strings"
)

// ErrorResult return error to client
type ErrorResult struct {
	Message string
	Stack   []byte
	State   interface{}
}

// String
func (e *ErrorResult) String() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintln(e.Message, e.Stack, e.State)
}

// Error return *ErrorResult
func Error(message string) *ErrorResult {
	return &ErrorResult{
		Message: message,
	}
}

// Execute write response
func (e *ErrorResult) Execute(ctx *HttpContext) error {
	if ctx.Server.Config.ErrorPageEnable {
		accetps := ctx.Accept()
		if strings.Contains(accetps, "text/html") {
			if f := ctx.Server.MapPath("public/error.html"); isFileExists(f) {
				http.ServeFile(ctx.Resonse, ctx.Request, f)
				return nil
			}
			if ctx.Server.Config.ViewEnable {
				if f := ctx.Server.MapPath("views/error.html"); isFileExists(f) {
					ctx.ViewData["ctx"] = ctx
					ctx.ViewData["error"] = e
					return executeViewFile("error.html", ctx)
				}
			}
		}
		if strings.Contains(accetps, "text/plain") {
			if f := ctx.Server.MapPath("public/error.txt"); isFileExists(f) {
				http.ServeFile(ctx.Resonse, ctx.Request, f)
				return nil
			}
		}
	}

	http.Error(ctx.Resonse, e.Message, http.StatusInternalServerError)
	return nil
}

// // executeErrorResult
// func executeErrorResult(ctx *HttpContext, err error) {
// 	e := &ErrorResult{
// 		Err: err,
// 	}
// 	e.Execute(ctx)
// }
