// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	_jsonPrefix = ""
	_jsonIndent = "\t"
)

// JsonResult marshal data and write to response
// ContentType is "application/json"
type JsonResult struct {
	NeedIndent bool
	// Prefix     string
	// Indent     string
	Data   interface{}
	buffer io.Reader
}

// Json return *JsonResult
func Json(a interface{}) *JsonResult {
	return &JsonResult{
		Data:       a,
		NeedIndent: false,
		// Prefix:     _jsonPrefix,
		// Indent:     _jsonIndent,
	}
}

// marshal
func (j *JsonResult) marshal() ([]byte, error) {
	if j.NeedIndent {
		//return json.MarshalIndent(j.Data, j.Prefix, j.Indent)
		return json.MarshalIndent(j.Data, _jsonPrefix, _jsonIndent)
	}
	return json.Marshal(j.Data)
}

// Execute encode result as json and write to response
func (j *JsonResult) Execute(ctx *HttpContext) {
	if !j.NeedIndent {
		w := ctx.Resonse.(io.Writer)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(j.Data)

		if err != nil {
			executeErrorResult(ctx, err)
		}
		return
	}

	b, err := j.marshal()
	if err != nil {
		executeErrorResult(ctx, err)
		return
	}

	ctx.ContentType(ContentTypeJson)
	ctx.Write(b)

	//
}

// ContentType return "application/json"
func (j *JsonResult) ContentType() string {
	return ContentTypeJson
}

// Read reads marshaled data
// TODO: bug fix
func (j *JsonResult) Read(p []byte) (n int, err error) {
	if j.buffer == nil {
		b, err := j.marshal()
		if err != nil {
			return 0, err
		}
		j.buffer = bytes.NewBuffer(b)
	}

	return j.buffer.Read(p)
}

// WriteTo writes marshaled data to w
func (j *JsonResult) WriteTo(w io.Writer) (n int64, err error) {
	b, err := j.marshal()
	if err != nil {
		return 0, err
	}
	buffer := bytes.NewBuffer(b)
	return buffer.WriteTo(w)
}
