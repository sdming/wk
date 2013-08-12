// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
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
func (c *ContentResult) Execute(ctx *HttpContext) error {
	if c.ContentType != "" {
		if ctype := ctx.ReqHeader("Content-Type"); ctype != "" {
			ctx.SetHeader("Content-Type", c.ContentType)
		}
	}

	if w, ok := c.Data.(io.WriterTo); ok {
		w.WriteTo(ctx.Resonse)
		return nil
	}

	if r, ok := c.Data.(io.Reader); ok {
		io.Copy(ctx.Resonse, r)
		return nil
	}

	if b, ok := c.Data.([]byte); ok {
		//  Write([]byte) (int, error)
		ctx.Resonse.Write(b)
		return nil
	}

	if s, ok := c.Data.(string); ok {
		io.WriteString(ctx.Resonse, s)
		return nil
	}

	fmt.Fprintln(ctx.Resonse, c.Data)
	return nil
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
func (c *DataResult) Execute(ctx *HttpContext) error {
	_, err := fmt.Fprintln(ctx.Resonse, c.Data)
	return err
}

// TextResult is plaintext
type TextResult string

// Text return a *TextResult
func Text(data string) TextResult {
	return TextResult(data)
}

// Execute write Data to response
func (t TextResult) Execute(ctx *HttpContext) error {
	if ctype := ctx.ReqHeader("Content-Type"); ctype != "" {
		ctx.SetHeader("Content-Type", ContentTypeText)
	}
	_, err := ctx.Resonse.Write([]byte(t))
	return err
}
