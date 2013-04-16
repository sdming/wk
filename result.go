// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import ()

// HttpResult is a interface that define how to write server reply to response
type HttpResult interface {
	// Execute 
	Execute(ctx *HttpContext)
}

// ContentType is a interface that return ContentType of response 
type ContentType interface {
	ContentType() string
}

// HeadSetter is a interface that set http response header 
type HeadSetter interface {
	Set(ctx *HttpContext) string
}
