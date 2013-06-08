// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"bytes"
	"mime"
	"path/filepath"
)

// ViewResult
type ViewResult struct {
	File string
}

// View return *ViewResult
func View(file string) *ViewResult {
	return &ViewResult{
		File: file,
	}
}

// String
func (v *ViewResult) String() string {
	return v.File
}

func executeViewFile(file string, ctx *HttpContext) error {
	view := &ViewResult{
		File: file,
	}
	return view.Execute(ctx)
}

// Execute
func (v *ViewResult) Execute(ctx *HttpContext) error {
	buffer := &bytes.Buffer{}

	err := ctx.Server.ViewEngine.Execte(buffer, v.File, ctx.ViewData)
	if err != nil {

		return err
	}

	ctx.SetHeader(HeaderContentType, v.Type())
	//ctx.SetHeader(HeaderContentLength, strconv.Itoa(len(buffer.Bytes())))
	_, err = buffer.WriteTo(ctx.Resonse)
	return err
}

// ContentType return mime type
func (v *ViewResult) Type() string {
	ctype := mime.TypeByExtension(filepath.Ext(v.File))
	return ctype
}

// // Write execute template and write output to body
// func (v *ViewResult) Write(header http.Header, body io.Writer) error {
// 	buffer := &bytes.Buffer{}
// 	err := DefaultViewEngine.Execte(buffer, v.File, v.Data)
// 	if err != nil {
// 		return err
// 	}

// 	header.Set(HeaderContentType, v.Type())
// 	_, err = buffer.WriteTo(body)
// 	return err
// }
