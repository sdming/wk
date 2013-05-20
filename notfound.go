// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"net/http"
)

// NotFoundResult is wrap of http 404
type NotFoundResult struct {
}

// Execute
// TODO: display custome 404 page
func (r *NotFoundResult) Execute(ctx *HttpContext) error {
	http.Error(ctx.Resonse, msgNotFound, http.StatusNotFound)
	return nil
}
