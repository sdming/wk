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

// XmlResult ("application/xml")
type XmlResult struct {
	Data       interface{}
	NeedIndent bool
	Prefix     string
	Indent     string
}

func Xml(a interface{}) *XmlResult {
	return &XmlResult{
		Data:       a,
		NeedIndent: false,
		Prefix:     _xmlPrefix,
		Indent:     _xmlIndent,
	}
}

func (x *XmlResult) marshal() ([]byte, error) {
	if x.NeedIndent {
		return xml.MarshalIndent(x.Data, x.Prefix, x.Indent)
	}
	return xml.Marshal(x.Data)
}

// Execute marshal result as xml
func (x *XmlResult) Execute(ctx *HttpContext) {
	b, err := x.marshal()
	if err != nil {
		executeErrorResult(ctx, err)
		return
	}

	ctx.ContentType(ContentTypeXml)
	ctx.Write(b)
}

// Read reads marshaled data
func (x *XmlResult) Read(p []byte) (n int, err error) {
	b, err := x.marshal()
	if err != nil {
		return 0, err
	}
	buffer := bytes.NewBuffer(b)
	return buffer.Read(p)
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

// // Text return 
// func (x *XmlResult) Text() (text []byte, err error) {
// 	b, err := x.marshal()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return b, nil
// }
