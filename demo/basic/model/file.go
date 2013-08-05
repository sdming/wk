// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*

*/
package model

import (
	"bytes"
	"fmt"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/boot"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	boot.Boot(RegisterFileRoute)
}

var serverPublicBase string

func RegisterFileRoute(server *wk.HttpServer) {
	serverPublicBase = server.Config.PublicDir

	// url: get /file/time.txt
	server.RouteTable.Get("/file/time.txt").To(FileHelloTime)

	// url: get /js/js_bundle.js
	server.RouteTable.Get("/js/js_bundle.js").To(FileJsBundling)

	// url: post /file/fileajax
	server.RouteTable.Post("/file/fileajax").To(FileAjax)

	// url: post /file/upload
	server.RouteTable.Post("/file/upload").To(FileUpload)

	// url: get /file/upload
	server.RouteTable.Get("/file/upload").To(FileUploadView)

	// url: get /file/absolute
	server.RouteTable.Get("/file/absolute").To(FileAbsolute)

	// url: get /file/relative
	server.RouteTable.Get("/file/relative").To(FileRelative)
}

func FileAbsolute(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.File(path.Join(ctx.Server.Config.RootDir, "public/humans.txt")), nil
}

func FileRelative(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.File("~/public/humans.txt"), nil
}

func FileHelloTime(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	s := "hello, time is " + time.Now().String()
	reader := strings.NewReader(s)
	return wk.FileStream("", "hellotime.txt", reader, time.Now()), nil
}

func FileJsBundling(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	files := []string{"js/main.js", "js/plugins.js"}
	bundle := make([]string, len(files))

	for i := 0; i < len(files); i++ {
		bundle[i] = filepath.Join(serverPublicBase, files[i])
	}
	return &wk.BundleResult{Files: bundle}, nil
}

func FileUploadView(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	return wk.View("doc/upload_demo.html"), nil
}

type UploadStatus struct {
	Files []FileInfo `json:"files"`
}

type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

func fileInfo(header *multipart.FileHeader) FileInfo {
	info := FileInfo{
		Name: header.Filename,
		Type: header.Header.Get("Content-Type"),
	}

	if f, err := header.Open(); err == nil {
		info.Size, err = f.Seek(0, os.SEEK_END)
		if err != nil {
			log.Println("Parse file info, seek error", err)
		}
		f.Close()
	} else {
		log.Println("Parse file info, open error", err)
	}
	return info
}

func FileAjax(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	ctx.Request.ParseMultipartForm(32 << 20) // 32 MB
	status := UploadStatus{make([]FileInfo, 0)}
	for _, headers := range ctx.Request.MultipartForm.File {
		for _, header := range headers {
			info := fileInfo(header)
			status.Files = append(status.Files, info)
		}
	}
	return wk.Json(status), nil
}

func FileUpload(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	ctx.Request.ParseMultipartForm(32 << 20) // 32 MB
	b := &bytes.Buffer{}
	b.WriteString("files uploaded:\n")
	for _, headers := range ctx.Request.MultipartForm.File {
		for _, header := range headers {
			info := fileInfo(header)
			b.WriteString(fmt.Sprintf("Name: %s; Content-Type: %s; Size: %d \n", info.Name, info.Type, info.Size))
		}
	}
	ctx.ViewData["msg"] = b.String()
	return wk.View("share/message.html"), nil
}
