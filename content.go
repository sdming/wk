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

func Content(data interface{}, contentType string) *ContentResult {
	return &ContentResult{
		Data:        data,
		ContentType: contentType,
	}
}

// Execute render raw content
func (c *ContentResult) Execute(ctx *HttpContext) {
	if c.ContentType != "" {
		if ctype := ctx.Header("Content-Type"); ctype != "" {
			ctx.SetHeader("Content-Type", c.ContentType)
		}
	}

	if r, ok := c.Data.(io.WriterTo); ok {
		r.WriteTo(ctx.Resonse)
		return
	}

	if r, ok := c.Data.(io.Reader); ok {
		io.Copy(ctx.Resonse, r)
		return
	}

	fmt.Fprintln(ctx.Resonse, c.Data)
}

// WriteTo writes data to w
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

// DataResult is 
type DataResult struct {

	// Data
	Data interface{}
}

func Data(data interface{}) *DataResult {
	return &DataResult{
		Data: data,
	}
}

// Execute render raw content
func (c *DataResult) Execute(ctx *HttpContext) {
	fmt.Fprintln(ctx.Resonse, c.Data)
}
