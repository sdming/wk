// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

// JsonpResult
type JsonpResult struct {
	Data interface{}
}

// Execute
// Content-Type is application/javascript
func (j *JsonpResult) Execute(ctx *HttpContext) {
	panic("TODO:JsonpResult")
}
