// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import ()

// 
type FormatList []FormatFunc

//
type FormatFunc func(*HttpContext, interface{}) (HttpResult, bool)

// Append add a FormatFunc 
func (f *FormatList) Append(fn FormatFunc) {
	*f = append(*f, fn)
}

// Remove a FormatFunc
func (f *FormatList) Remove(fn FormatFunc) {
	panic("TODO:")
}

var Formatter FormatList

//
func init() {
	Formatter = make([]FormatFunc, 0)
}
