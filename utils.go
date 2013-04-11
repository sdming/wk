// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"errors"
	"fmt"
	"github.com/sdming/kiss/gotype"
	"net/http"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
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

// http://doc.golang.org/src/pkg/net/http/fs.go
// modtime is the modification time of the resource to be served, or IsZero().
// return value is whether this request is now complete.
func checkLastModified(w http.ResponseWriter, r *http.Request, modtime time.Time) bool {
	if modtime.IsZero() {
		return false
	}

	// The Date-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if t, err := time.Parse(formatTime, r.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(1*time.Second)) {
		w.WriteHeader(http.StatusNotModified)
		return true
	}
	w.Header().Set("Last-Modified", modtime.UTC().Format(formatTime))
	return false
}

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)

// http://doc.golang.org/src/pkg/net/http/server.go
func htmlEscape(s string) string {
	return htmlReplacer.Replace(s)
}

func formatXml(ctx *HttpContext, x interface{}) (HttpResult, bool) {
	return &XmlResult{Data: x}, true
}

func formatJson(ctx *HttpContext, x interface{}) (HttpResult, bool) {
	return &JsonResult{Data: x}, true
}

// convertResult convert reflect.Value to http result
func convertResult(ctx *HttpContext, v reflect.Value) HttpResult {
	i := v.Interface()

	if r, ok := i.(HttpResult); ok {
		return r
	}

	if Formatter != nil {
		for _, f := range Formatter {
			if formatted, ok := f(ctx, v.Interface()); ok {
				return formatted
			}
		}
	}

	kind := reflect.Indirect(v).Kind()

	switch {
	case gotype.IsSimple(kind):
		return &ContentResult{Data: i}
	case gotype.IsStruct(kind) || gotype.IsCollect(kind):
		accept := ctx.Accept()
		switch {
		case strings.Index(accept, "xml") > -1:
			return &XmlResult{Data: v.Interface()}
		case strings.Index(accept, "jsonp") > -1:
			return &JsonpResult{Data: v.Interface()}
		default:
			return &JsonResult{Data: v.Interface()}
		}
	}
	return &ContentResult{Data: i}
}
