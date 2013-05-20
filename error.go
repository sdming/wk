// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"fmt"
	"net/http"
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
// TODO: cutomer error view page
func (e *ErrorResult) Execute(ctx *HttpContext) error {
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
