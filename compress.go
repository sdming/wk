// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"

	"strings"
)

// CompressProcessor compress http response with gzip or deflate
// TODO: copy header from original result
// TODO: configurable
// TODO: filte by MimeType
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
func (p *CompressProcessor) Register(server *HttpServer) {
	//TODO: Read config file
	p.Enable = true
	p.Level = flate.BestSpeed
}

// Execute convert result to  CompressResult if can 
func (p *CompressProcessor) Execute(ctx *HttpContext) {

	if ctx.Result == nil || ctx.Error != nil {
		return
	}

	data, ok := ctx.Result.(io.Reader)
	if !ok {
		return
	}

	accept := ctx.Header("Accept-Encoding")
	if accept == "" {
		return
	}

	var contenType string

	if ctx.Resonse.Header().Get("Content-Type") == "" {
		if typ, ok := ctx.Result.(ContentType); ok {
			contenType = typ.ContentType()
		} else {
			//contenType = http.DetectContentType(b)
		}
	}

	encodings := strings.Split(accept, ",")
	for _, encoder := range encodings {
		if encoder == "deflate" {
			ctx.Result = &CompressResult{
				Level:       p.Level,
				Data:        data,
				Format:      "deflate",
				ContentType: contenType,
			}
		} else if encoder == "gzip" {
			ctx.Result = &CompressResult{
				Level:       p.Level,
				Data:        data,
				Format:      "gzip",
				ContentType: contenType,
			}
		}
		return
	}

	return

}

// CompressResult compress http response
type CompressResult struct {
	Data        io.Reader
	Level       int
	Format      string
	ContentType string
}

// Execute write compressed Data
// TODO: handle err
func (c *CompressResult) Execute(ctx *HttpContext) {

	if c.Data == nil {
		return
	}

	var writer io.WriteCloser

	if c.Format == "gzip" {
		ctx.SetHeader("Content-Encoding", "gzip")
		writer, _ = gzip.NewWriterLevel(ctx.Resonse, c.Level)
	} else if c.Format == "deflate" {
		ctx.SetHeader("Content-Encoding", "deflate")
		writer, _ = zlib.NewWriterLevel(ctx.Resonse, c.Level)
	} else {
		panic("Not Implemented")
	}

	defer writer.Close()

	if ctx.Resonse.Header().Get("Content-Type") == "" && c.ContentType != "" {
		ctx.ContentType(c.ContentType)
	}

	var err error

	if w, ok := c.Data.(io.WriterTo); ok {
		w.WriteTo(writer)
	} else {
		io.Copy(writer, c.Data)
	}

	//TODO: handle err
	if err != nil {
		Logger.Println("debug compress error", err)
	}
}
