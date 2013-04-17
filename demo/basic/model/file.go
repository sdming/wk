// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*

*/
package model

import (
	"github.com/sdming/wk"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func RegisterFileRoute(server *wk.HttpServer) {
	// url: get /file/time.txt
	server.RouteTable.Get("/file/time.txt").To(FileHelloTime)

	// url: get /js/js_bundle.js
	server.RouteTable.Get("/js/js_bundle.js").To(FileJsBundling)
}

func FileHelloTime(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	s := "hello, time is " + time.Now().String()
	reader := strings.NewReader(s)
	return wk.FileStream("", "hellotime.txt", reader, time.Now()), nil
}

// TODO: close reader?
func FileJsBundling(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	files := []string{"main.js", "plugins.js"}
	base := "./public/js/"
	readers := make([]io.Reader, len(files))

	var modtime time.Time

	for i := 0; i < len(files); i++ {
		f, err := os.Open(filepath.Join(base, files[i]))

		if err != nil {
			return nil, err
		}

		d, err := f.Stat()
		if err != nil {
			return nil, err
		}

		if d.ModTime().After(modtime) {
			modtime = d.ModTime()
		}

		readers[i] = f
	}
	return wk.FileStream("application/javascript", "", io.MultiReader(readers[:]...), modtime), nil

}
