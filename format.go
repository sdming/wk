// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import ()

// FormatList
type FormatList []FormatFunc

// FormatFunc return formatted HttpResult or return (nil, false) if doesn't format it
type FormatFunc func(*HttpContext, interface{}) (HttpResult, bool)

// Append register a FormatFunc 
func (f *FormatList) Append(fn FormatFunc) {
	*f = append(*f, fn)
}

// Remove unregister FormatFunc
func (f *FormatList) Remove(fn FormatFunc) {
	panic("TODO:")
}

var Formatters FormatList

//
func init() {
	Formatters = make([]FormatFunc, 0)
}
