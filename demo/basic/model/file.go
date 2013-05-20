// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*

*/
package model

import (
	"github.com/sdming/wk"
	"log"
	"path/filepath"
	"strings"
	"time"
)

var serverPublicBase string

func RegisterFileRoute(server *wk.HttpServer) {
	serverPublicBase = server.Config.PublicDir

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
	files := []string{"js/main.js", "js/plugins.js"}
	bundle := make([]string, len(files))

	for i := 0; i < len(files); i++ {
		bundle[i] = filepath.Join(serverPublicBase, files[i])
	}
	log.Println("FileJsBundling", bundle)
	return &wk.BundleResult{Files: bundle}, nil

}
