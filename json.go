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

// JsonResult  ("application/json")
type JsonResult struct {
	NeedIndent bool
	Prefix     string
	Indent     string
	Data       interface{}
}

func Json(a interface{}) *JsonResult {
	return &JsonResult{Data: a,
		NeedIndent: false,
		Prefix:     _jsonPrefix,
		Indent:     _jsonIndent,
	}
}

func (j *JsonResult) marshal() ([]byte, error) {
	if j.NeedIndent {
		return json.MarshalIndent(j.Data, j.Prefix, j.Indent)
	}
	return json.Marshal(j.Data)
}

// Execute marshal result as json
func (j *JsonResult) Execute(ctx *HttpContext) {
	b, err := j.marshal()
	if err != nil {
		executeErrorResult(ctx, err)
		return
	}

	ctx.ContentType(ContentTypeJson)
	ctx.Write(b)
}

// Read reads marshaled data
func (j *JsonResult) Read(p []byte) (n int, err error) {
	b, err := j.marshal()
	if err != nil {
		return 0, err
	}
	buffer := bytes.NewBuffer(b)
	return buffer.Read(p)
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
