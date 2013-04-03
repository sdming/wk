// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"errors"
	"log"
)

var (
	// can not find view 
	errViewNotFound = errors.New("can not find view")

	// internal error
	errInternalError = errors.New(msgServerInternalErr)

	// errNoaction mean can not find actionin controller
	errNoaction = errors.New(msgNoAction)
)

var (
	// ResultVoid mean action did not return a result
	resultVoid = &VoidResult{}

	// resultNotFound mean can not find result
	resultNotFound = &NotFoundResult{}

	// no result
	resultNoResult = &ErrorResult{Err: errors.New(msgNoResult)}
)

var (
	// logger
	Logger *log.Logger

	// LogLevel is level of log
	LogLevel int = LogDebug

	// EnableProfile mean enable http profile or not
	EnableProfile bool = false
)

type ServerError struct {
	Message string
	Module  string
}

func (se *ServerError) Error() string {
	return se.Module + ":" + se.Message
}

// type Text interface {
// 	Text() ([]byte, error)
// }

type Handler func(*HttpContext)
