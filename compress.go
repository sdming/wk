// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"net/http"
	"strings"
)

// compresser is interface of gzip/flate writer
type compresser interface {
	Flush() error
	Close() error
	Write(p []byte) (int, error)
}

// CompressProcessor compress http response with gzip/deflate
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

// Execute convert ctx.Response to compressResponseWriter
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

	var writer compresser
	var err error
	var format string

	if strings.Contains(accept, "deflate") {
		format = "deflate"
		writer, err = zlib.NewWriterLevel(ctx.Resonse, cp.Level)
	} else if strings.Contains(accept, "gzip") {
		format = "gzip"
		writer, err = gzip.NewWriterLevel(ctx.Resonse, cp.Level)
	}

	if format == "" || err != nil {
		return
	}

	ctx.Resonse.Header().Set("Content-Encoding", format)
	ctx.Resonse = &compressResponseWriter{
		rw:          ctx.Resonse,
		writer:      writer,
		contentType: contenType,
		format:      format,
	}

}

// compressResponseWriter wrap gzip/deflate and ResponseWriter
type compressResponseWriter struct {
	rw            http.ResponseWriter
	writer        compresser
	contentType   string
	format        string
	headerWritten bool
}

func (crw *compressResponseWriter) Header() http.Header {
	return crw.rw.Header()
}

func (crw *compressResponseWriter) WriteHeader(status int) {
	crw.rw.WriteHeader(status)
}

func (crw *compressResponseWriter) Write(p []byte) (int, error) {
	if !crw.headerWritten {
		if crw.rw.Header().Get(HeaderContentType) == "" && crw.contentType != "" {
			crw.rw.Header().Set(HeaderContentType, crw.contentType)
		}
		crw.headerWritten = true
	}
	n, err := crw.writer.Write(p)
	err = crw.writer.Flush()
	return n, err
}
