// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

// RedirectResult
type RedirectResult struct {
	// Permanent mean redirect permanent or not
	Permanent bool

	// Path it the url to redirect
	Path string
}

// Execute
func (r *RedirectResult) Execute(ctx *HttpContext) {
	panic("TODO:RedirectResult")
}
