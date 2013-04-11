// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"errors"
	"net/http"
)

// ErrorResult 
type ErrorResult struct {
	Err error
	Tag string
}

// String
func (e *ErrorResult) String() string {
	if e == nil {
		return "<nil>"
	}
	return e.Err.Error()
}

// Error return a *ErrorResult
func Error(msg string) *ErrorResult {
	return &ErrorResult{
		Err: errors.New(msg),
	}
}

// Execute, TODO: cutomer error view page
func (e *ErrorResult) Execute(ctx *HttpContext) {
	http.Error(ctx.Resonse, e.String(), http.StatusInternalServerError)
}

func executeErrorResult(ctx *HttpContext, err error) {
	e := &ErrorResult{
		Err: err,
	}
	e.Execute(ctx)
}
