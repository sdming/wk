// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"bytes"
	"fmt"
	"io"
)

// ContentResult is raw content
type ContentResult struct {

	// ContentType
	ContentType string

	// Data
	Data interface{}
}

// Content return *ContentResult
func Content(data interface{}, contentType string) *ContentResult {
	return &ContentResult{
		Data:        data,
		ContentType: contentType,
	}
}

// Execute write Data to response
func (c *ContentResult) Execute(ctx *HttpContext) {
	if c.ContentType != "" {
		if ctype := ctx.Header("Content-Type"); ctype != "" {
			ctx.SetHeader("Content-Type", c.ContentType)
		}
	}

	if w, ok := c.Data.(io.WriterTo); ok {
		w.WriteTo(ctx.Resonse)
		return
	}

	if r, ok := c.Data.(io.Reader); ok {
		io.Copy(ctx.Resonse, r)
		return
	}

	if b, ok := c.Data.([]byte); ok {
		//  Write([]byte) (int, error)
		ctx.Resonse.Write(b)
		return
	}

	if s, ok := c.Data.(string); ok {
		io.WriteString(ctx.Resonse, s)
		return
	}

	fmt.Fprintln(ctx.Resonse, c.Data)
}

// WriteTo implement io.WriteTo
func (c *ContentResult) WriteTo(w io.Writer) (n int64, err error) {
	if wt, ok := c.Data.(io.WriterTo); ok {
		return wt.WriteTo(w)
	}

	if b, ok := c.Data.([]byte); ok {
		buf := bytes.NewBuffer(b)
		return buf.WriteTo(w)
	}

	ni, err := fmt.Fprint(w, c.Data)
	return int64(ni), err
}

// DataResult is wrap of simple type
type DataResult struct {

	// Data
	Data interface{}
}

// Data return *DataResult 
func Data(data interface{}) *DataResult {
	return &DataResult{
		Data: data,
	}
}

// Execute write Data to response 
func (c *DataResult) Execute(ctx *HttpContext) {
	fmt.Fprintln(ctx.Resonse, c.Data)
}
