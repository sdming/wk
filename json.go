// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	_jsonPrefix = ""
	_jsonIndent = "\t"
)

// JsonResult marshal data and write to response
// ContentType is "application/json"
type JsonResult struct {
	NeedIndent bool
	Data       interface{}
	buffer     io.Reader
}

// Json return *JsonResult
func Json(a interface{}) *JsonResult {
	return &JsonResult{
		Data:       a,
		NeedIndent: false,
	}
}

// marshal
func (j *JsonResult) marshal() ([]byte, error) {
	if j.NeedIndent {
		return json.MarshalIndent(j.Data, _jsonPrefix, _jsonIndent)
	}
	return json.Marshal(j.Data)
}

// Execute encode result as json and write to response
func (j *JsonResult) Execute(ctx *HttpContext) error {
	ctx.ContentType(ContentTypeJson)

	if !j.NeedIndent {
		w := ctx.Resonse.(io.Writer)
		encoder := json.NewEncoder(w)
		return encoder.Encode(j.Data)
	}

	b, err := j.marshal()
	if err != nil {
		return err
	}

	_, err = ctx.Write(b)
	return err
}

// ContentType return "application/json"
func (j *JsonResult) Type() string {
	return ContentTypeJson
}

// WriteTo writes marshaled data to w
func (j *JsonResult) Write(header http.Header, body io.Writer) error {
	header.Set(HeaderContentType, j.Type())

	b, err := j.marshal()
	if err != nil {
		return err
	}

	_, err = body.Write(b)
	return err
}

// // Read reads marshaled data
// // TODO: bug fix
// func (j *JsonResult) Read(p []byte) (n int, err error) {
// 	if j.buffer == nil {
// 		b, err := j.marshal()
// 		if err != nil {
// 			return 0, err
// 		}
// 		j.buffer = bytes.NewBuffer(b)
// 	}

// 	return j.buffer.Read(p)
// }

// // WriteTo writes marshaled data to w
// func (j *JsonResult) WriteTo(w io.Writer) (n int64, err error) {
// 	b, err := j.marshal()
// 	if err != nil {
// 		return 0, err
// 	}
// 	buffer := bytes.NewBuffer(b)
// 	return buffer.WriteTo(w)
// }
