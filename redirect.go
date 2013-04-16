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

// Execute write status code http.StatusMovedPermanently or http.StatusFound
func (r *RedirectResult) Execute(ctx *HttpContext) {
	var code int

	if r.Permanent {
		code = http.StatusMovedPermanently
		return
	} else {
		code = http.StatusFound
	}

	http.Redirect(ctx.Resonse, ctx.Request, r.UrlStr, code)
}
