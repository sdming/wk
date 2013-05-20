// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"bytes"
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

// ViewResult
type ViewResult struct {
	File string
	Data interface{}
}

// View return *ViewResult
func View(file string, data interface{}) *ViewResult {
	return &ViewResult{
		File: file,
		Data: data,
	}
}

// String
func (v *ViewResult) String() string {
	if v == nil {
		return "<nil>"
	}
	return "view:" + v.File
}

// Execute
func (v *ViewResult) Execute(ctx *HttpContext) error {
	// buffer := &bytes.Buffer{}
	// err := GoHtml.Execte(buffer, v.File, v.Data)
	// if err != nil {
	// 	return err
	// }

	// ctx.ContentType(v.ContentType())
	// buffer.WriteTo(ctx.Resonse)
	return v.Write(ctx.Resonse.Header(), ctx.Resonse)
}

// ContentType return mime type
func (v *ViewResult) Type() string {
	ctype := mime.TypeByExtension(filepath.Ext(v.File))
	return ctype
}

// Write execute template and write output to body
func (v *ViewResult) Write(header http.Header, body io.Writer) error {
	buffer := &bytes.Buffer{}
	err := DefaultViewEngine.Execte(buffer, v.File, v.Data)
	if err != nil {
		return err
	}

	header.Set(HeaderContentType, v.Type())
	_, err = buffer.WriteTo(body)
	return err
}
