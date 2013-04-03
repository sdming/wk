// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"net/http"
)

// NotModifiedResult
type NotModifiedResult struct {
}

// Execute
func (r *NotModifiedResult) Execute(ctx *HttpContext) {
	ctx.Status(http.StatusNotModified)
}
