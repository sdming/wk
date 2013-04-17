// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"errors"
	"log"
	"reflect"
)

var (
	// can not find view 
	errViewNotFound = errors.New("can not find view")

	// internal error
	errInternalError = errors.New(msgServerInternalErr)

	// can not find actionin method
	errNoAction = errors.New(msgNoAction)

	// httpresult is nil
	errNoResult = &ErrorResult{
		Err: errors.New(msgNoResult),
	}
)

var (
	// action did not return a result
	resultVoid = &VoidResult{}

	// httpresult is nil
	resultNotFound = &NotFoundResult{}

	// no result
	resultNoResult = &ErrorResult{Err: errors.New(msgNoResult)}
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
	Execute(*HttpContext)
}
