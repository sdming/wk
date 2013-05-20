// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"encoding/xml"
	"io"
	"net/http"
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
}

// Xml return *XmlResult
func Xml(a interface{}) *XmlResult {
	return &XmlResult{
		Data:       a,
		NeedIndent: false,
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
func (x *XmlResult) Execute(ctx *HttpContext) error {
	ctx.ContentType(ContentTypeXml)

	if !x.NeedIndent {
		encoder := xml.NewEncoder(ctx.Resonse)
		return encoder.Encode(x.Data)
	}

	b, err := x.marshal()
	if err != nil {
		return err
	}

	_, err = ctx.Write(b)
	return err
}

func (x *XmlResult) Type() string {
	return ContentTypeXml
}

// WriteTo writes marshaled data to w
func (x *XmlResult) Write(header http.Header, body io.Writer) error {
	header.Set(HeaderContentType, x.Type())

	b, err := x.marshal()
	if err != nil {
		return err
	}

	_, err = body.Write(b)
	return err
}

// // Read reads marshaled data
// func (x *XmlResult) Read(p []byte) (n int, err error) {
// 	if x.buffer == nil {
// 		b, err := x.marshal()
// 		if err != nil {
// 			return 0, err
// 		}
// 		x.buffer = bytes.NewBuffer(b)
// 	}

// 	return x.buffer.Read(p)
// }

// // WriteTo writes marshaled data to w
// func (x *XmlResult) WriteTo(w io.Writer) (n int64, err error) {
// 	b, err := x.marshal()
// 	if err != nil {
// 		return 0, err
// 	}
// 	buffer := bytes.NewBuffer(b)
// 	return buffer.WriteTo(w)
// }
