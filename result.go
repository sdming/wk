// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import ()

// HttpResult interface
type HttpResult interface {
	// Execute 
	Execute(ctx *HttpContext)
}
