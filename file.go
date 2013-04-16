// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

// FileResult is wrap of static file 
type FileResult struct {
	Path string
	file http.File
}

// String
func (file *FileResult) String() string {
	if file == nil {
		return "<nil>"
	}
	return file.Path
}

// File return *FileResult
func File(path string) *FileResult {
	return &FileResult{
		Path: path,
	}
}

// Execute replies to the request with the contents of FileResult.Path
func (file *FileResult) Execute(ctx *HttpContext) {
	http.ServeFile(ctx.Resonse, ctx.Request, file.Path)
}

// ContentType return mime type of file
func (file *FileResult) ContentType() string {
	ctype := mime.TypeByExtension(filepath.Ext(file.Path))

	// var buf [1024]byte
	// n, _ := io.ReadFull(content, buf[:])
	// b := buf[:n]
	// ctype = DetectContentType(b)
	// _, err := content.Seek(0, os.SEEK_SET)

	return ctype
}

// Read reads data from file
func (file *FileResult) Read(p []byte) (n int, err error) {
	if file.file != nil {
		return file.file.Read(p)
	}

	dir, name := filepath.Split(file.Path)

	file.file, err = http.Dir(dir).Open(name)
	if err != nil {
		return
	}

	Logger.Println("len", len(p))
	n, err = file.file.Read(p)
	Logger.Println("file.file", n, err)
	return n, err
}

// FileStreamResult 
type FileStreamResult struct {
	// ContentType
	CType string

	DownloadName string

	ModifyTime time.Time

	// Data ?Data io.ReadSeeker
	Data io.Reader
}

// FileStream return *FileStream
func FileStream(contentType, downloadName string, reader io.Reader, modtime time.Time) *FileStreamResult {
	return &FileStreamResult{
		CType:        contentType,
		DownloadName: downloadName,
		Data:         reader,
		ModifyTime:   modtime,
	}
}

func (file *FileStreamResult) Execute(ctx *HttpContext) {
	//Set ContentType = "application/octet-stream"; 

	if file.DownloadName != "" {
		ctx.SetHeader("Content-Disposition", "attachment; filename=\""+file.DownloadName+"\";")
	}

	if ctype := ctx.Resonse.Header().Get("Content-Type"); ctype == "" {
		ctype = file.ContentType()
		if ctype != "" {
			ctx.ContentType(ctype)
		}
	}

	if rs, ok := file.Data.(io.ReadSeeker); ok {
		http.ServeContent(ctx.Resonse, ctx.Request, file.DownloadName, file.ModifyTime, rs)
		return
	}

	io.Copy(ctx.Resonse, file.Data)

}

// ContentType return mime type of file
func (file *FileStreamResult) ContentType() string {
	if file.CType != "" {
		return file.CType
	}

	if file.DownloadName != "" {
		return mime.TypeByExtension(filepath.Ext(file.DownloadName))
	}

	return "application/octet-stream"

	// var buf [1024]byte
	// n, _ := io.ReadFull(content, buf[:])
	// b := buf[:n]
	// ctype = DetectContentType(b)
	// _, err := content.Seek(0, os.SEEK_SET)
}

// Read reads data from file
func (file *FileStreamResult) Read(p []byte) (n int, err error) {
	return file.Read(p)
}
