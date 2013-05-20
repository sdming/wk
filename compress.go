// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"strings"
)

// CompressProcessor compress http response with gzip or deflate
// TODO: copy header from original result
// TODO: configurable
// TODO: filter by MimeType
type CompressProcessor struct {
	Enable   bool
	Level    int
	MimeType string
}

// NewCompressProcess return a *Process that wrap CompressProcessor
func NewCompressProcess(name, method, path string) *Process {
	return &Process{
		Path:   path,
		Method: method,
		Name:   name,
		Handler: &CompressProcessor{
			Enable: true,
			Level:  flate.BestSpeed,
		},
	}
}

// Register initialize CompressProcessor
func (cp *CompressProcessor) Register(server *HttpServer) {
	//TODO: Read config file
	cp.Enable = true
	cp.Level = flate.BestSpeed
}

// Execute convert result to  CompressResult
func (cp *CompressProcessor) Execute(ctx *HttpContext) {

	if ctx.Result == nil || ctx.Error != nil {
		return
	}

	if ctx.Method == HttpVerbsHead {
		return
	}

	accept := ctx.ReqHeader("Accept-Encoding")
	if accept == "" {
		return
	}

	var contenType string

	if ctx.ResHeader("Content-Type") == "" {
		if typ, ok := ctx.Result.(ContentTyper); ok {
			contenType = typ.Type()
		} else {
			//TODO: DetectContentType
			//contenType = http.DetectContentType(b)
		}
	}

	if contenType == "" {
		return
	}

	var writer io.Writer
	var err error
	var format string

	encodings := strings.Split(accept, ",")
	for _, encoder := range encodings {
		if encoder == "deflate" {
			format = encoder
			writer, err = zlib.NewWriterLevel(ctx.Resonse, cp.Level)
		} else if encoder == "gzip" {
			writer, err = gzip.NewWriterLevel(ctx.Resonse, cp.Level)
		}

		if format != "" {
			break
		}
	}

	if format == "" || err != nil {
		return
	}

	ctx.Resonse = &compressResponseWriter{
		ctx:         ctx,
		writer:      writer,
		contentType: contenType,
		format:      format,
	}

}

// compressResponseWriter wrap gzip/deflate and ResponseWriter
type compressResponseWriter struct {
	ctx           *HttpContext
	writer        io.Writer
	contentType   string
	format        string
	headerWritten bool
}

func (crw *compressResponseWriter) Header() http.Header {
	return crw.ctx.Resonse.Header()
}

func (crw *compressResponseWriter) WriteHeader(status int) {
	crw.ctx.Resonse.WriteHeader(status)
}

func (crw *compressResponseWriter) Write(p []byte) (int, error) {
	if !crw.headerWritten {
		crw.ctx.SetHeader("Content-Encoding", crw.format)
		if crw.ctx.ResHeader("Content-Type") == "" && crw.contentType != "" {
			crw.ctx.ContentType(crw.contentType)
		}
		crw.headerWritten = true
	}
	return crw.writer.Write(p)
}
