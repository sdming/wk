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
	//file http.File
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
func (file *FileResult) Execute(ctx *HttpContext) error {
	http.ServeFile(ctx.Resonse, ctx.Request, file.Path)
	return nil
}

// ContentType return mime type of file
func (file *FileResult) Type() string {
	ctype := mime.TypeByExtension(filepath.Ext(file.Path))

	// var buf [1024]byte
	// n, _ := io.ReadFull(content, buf[:])
	// b := buf[:n]
	// ctype = DetectContentType(b)
	// _, err := content.Seek(0, os.SEEK_SET)

	return ctype
}

// // Read reads data from file
// func (file *FileResult) Read(p []byte) (n int, err error) {
// 	if file.file != nil {
// 		return file.file.Read(p)
// 	}

// 	dir, name := filepath.Split(file.Path)

// 	file.file, err = http.Dir(dir).Open(name)
// 	if err != nil {
// 		return
// 	}

// 	n, err = file.file.Read(p)
// 	return n, err
// }

// FileStreamResult
type FileStreamResult struct {
	// ContentType
	ContentType string

	DownloadName string

	ModifyTime time.Time

	// Data ?Data io.ReadSeeker
	Data io.Reader
}

// FileStream return *FileStream
func FileStream(contentType, downloadName string, reader io.Reader, modtime time.Time) *FileStreamResult {
	return &FileStreamResult{
		ContentType:  contentType,
		DownloadName: downloadName,
		Data:         reader,
		ModifyTime:   modtime,
	}
}

func (file *FileStreamResult) Execute(ctx *HttpContext) error {
	//Set ContentType = "application/octet-stream";

	if file.DownloadName != "" {
		ctx.SetHeader("Content-Disposition", "attachment; filename=\""+file.DownloadName+"\";")
	}

	if ctype := ctx.Resonse.Header().Get("Content-Type"); ctype == "" {
		ctype = file.Type()
		if ctype != "" {
			ctx.ContentType(ctype)
		}
	}

	if rs, ok := file.Data.(io.ReadSeeker); ok {
		http.ServeContent(ctx.Resonse, ctx.Request, file.DownloadName, file.ModifyTime, rs)
		return nil
	}

	if checkLastModified(ctx.Resonse, ctx.Request, file.ModifyTime) {
		return nil
	}

	io.Copy(ctx.Resonse, file.Data)
	return nil
}

// ContentType return mime type of file
func (file *FileStreamResult) Type() string {
	if file.ContentType != "" {
		return file.ContentType
	}

	if file.DownloadName != "" {
		return mime.TypeByExtension(filepath.Ext(file.DownloadName))
	}

	return ""
	//return "application/octet-stream"

	// var buf [1024]byte
	// n, _ := io.ReadFull(content, buf[:])
	// b := buf[:n]
	// ctype = DetectContentType(b)
	// _, err := content.Seek(0, os.SEEK_SET)
}
