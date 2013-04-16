// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"bytes"
	"encoding/xml"
	"io"
)

const (
	_xmlPrefix = ""
	_xmlIndent = "\t"
)

// XmlResult  marshal data to xml and write to response
// ContentType: "application/xml"
type XmlResult struct {
	buffer     io.Reader
	Data       interface{}
	NeedIndent bool
	// Prefix     string
	// Indent     string
}

// Xml return *XmlResult
func Xml(a interface{}) *XmlResult {
	return &XmlResult{
		Data:       a,
		NeedIndent: false,
		// Prefix:     _xmlPrefix,
		// Indent:     _xmlIndent,
	}
}

// marshal
func (x *XmlResult) marshal() ([]byte, error) {
	if x.NeedIndent {
		return xml.MarshalIndent(x.Data, _xmlPrefix, _xmlIndent)
	}
	return xml.Marshal(x.Data)
}

// Execute encode result as xml and write to response
func (x *XmlResult) Execute(ctx *HttpContext) {
	if !x.NeedIndent {
		w := ctx.Resonse.(io.Writer)
		encoder := xml.NewEncoder(w)
		err := encoder.Encode(x.Data)

		if err != nil {
			executeErrorResult(ctx, err)
		}
		return
	}

	b, err := x.marshal()
	if err != nil {
		executeErrorResult(ctx, err)
		return
	}

	ctx.ContentType(ContentTypeXml)
	ctx.Write(b)
}

func (j *XmlResult) ContentType() string {
	return ContentTypeXml
}

// Read reads marshaled data
func (x *XmlResult) Read(p []byte) (n int, err error) {
	if x.buffer == nil {
		b, err := x.marshal()
		if err != nil {
			return 0, err
		}
		x.buffer = bytes.NewBuffer(b)
	}

	return x.buffer.Read(p)
}

// WriteTo writes marshaled data to w
func (x *XmlResult) WriteTo(w io.Writer) (n int64, err error) {
	b, err := x.marshal()
	if err != nil {
		return 0, err
	}
	buffer := bytes.NewBuffer(b)
	return buffer.WriteTo(w)
}
