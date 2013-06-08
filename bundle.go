// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// BundleResult is bundle of Files
// TODO: cache?
type BundleResult struct {
	Files []string
}

// Execute read from BundleResult.Files and write to HttpContext.Response
func (b *BundleResult) Execute(ctx *HttpContext) error {
	if len(b.Files) == 0 {
		return errors.New("bundle files is invalid")
	}

	var modtime time.Time

	for _, file := range b.Files {
		info, err := os.Stat(file)

		if err != nil {
			return err
		}
		if info.IsDir() {
			return errors.New("bundle files is invalid")
		}

		if info.ModTime().After(modtime) {
			modtime = info.ModTime()
		}
	}

	if checkLastModified(ctx.Resonse, ctx.Request, modtime) {
		return nil
	}

	ctx.ContentType(b.Type())

	buffer := &bytes.Buffer{}
	for _, file := range b.Files {
		err := copyFromFile(file, buffer)

		if err != nil {
			return err
		}
	}
	http.ServeContent(ctx.Resonse, ctx.Request, b.Files[0], modtime, bytes.NewReader(buffer.Bytes()))

	return nil
}

func copyFromFile(file string, w io.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}

// ContentType return mime type of BundleResult.Files[0]
func (b *BundleResult) Type() string {
	if len(b.Files) == 0 {
		return ""
	}
	return mime.TypeByExtension(filepath.Ext(b.Files[0]))
}
