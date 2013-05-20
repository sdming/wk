// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"net/http"
)

// RedirectResult is wrap of http.StatusMovedPermanently or http.StatusFound
type RedirectResult struct {
	// Permanent mean redirect permanent or not
	Permanent bool

	// UrlStr is the url to redirect
	UrlStr string
}

func Redirect(urlStr string, permanent bool) *RedirectResult {
	return &RedirectResult{
		UrlStr:    urlStr,
		Permanent: permanent,
	}
}

// Execute write status code http.StatusMovedPermanently or http.StatusFound
func (r *RedirectResult) Execute(ctx *HttpContext) error {
	var code int

	if r.Permanent {
		code = http.StatusMovedPermanently
	} else {
		code = http.StatusFound
	}

	http.Redirect(ctx.Resonse, ctx.Request, r.UrlStr, code)
	return nil
}
