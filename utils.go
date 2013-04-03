// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"time"
)

func isFileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func isDirExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return true
	}
	return false
}

// Return the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

// returns the name of the calling method, Caller(N)
func methodNameN(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}
	return f.Name()
}

const formatTime = "Mon, 02 Jan 2006 15:04:05 GMT"

func webTime(t time.Time) string {
	return t.Format(formatTime)
}

// safeCall
func safeCall(fn reflect.Value, args []reflect.Value) (result []reflect.Value, err error) {
	defer func() {
		if x := recover(); x != nil {
			if e, ok := x.(error); ok {
				err = e
			} else {
				err = errors.New(fmt.Sprintf("call method %s fail, %s ", fn.Type(), x))
			}
		}

	}()

	return fn.Call(args), nil
}
