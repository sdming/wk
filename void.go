// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

// VoidResult
type VoidResult struct {
}

// Execute
func (r *VoidResult) Execute(ctx *HttpContext) {
	ctx.Resonse.Write([]byte(``))
}
