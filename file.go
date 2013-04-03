// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"net/http"
)

// StaticFileResult is static file 
type FileResult struct {
	Path string
}

// String
func (file *FileResult) String() string {
	if file == nil {
		return "<nil>"
	}
	return file.Path
}

func File(path string) *FileResult {
	return &FileResult{
		Path: path,
	}
}

// Execute rende a static file
func (file *FileResult) Execute(ctx *HttpContext) {
	http.ServeFile(ctx.Resonse, ctx.Request, file.Path)
}

// StreamFileResult is static file ,:TODO
type StreamFileResult struct {
}
