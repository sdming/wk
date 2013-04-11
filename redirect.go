// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"net/http"
)

// RedirectResult
type RedirectResult struct {
	// Permanent mean redirect permanent or not
	Permanent bool

	// UrlStr  it the url to redirect
	UrlStr string
}

// Execute
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
