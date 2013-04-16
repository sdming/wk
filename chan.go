// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"time"
)

var (
	chanTimeout   = 60 * time.Second
	chanBufferLen = 9
)

// ChanResult
// TODO: enhance it
type ChanResult struct {
	Len   int
	Chan  chan string
	CType string
	Start []byte
	End   []byte
}

// Execute read string from chan and write to response
// TODO: enhance it
func (c *ChanResult) Execute(ctx *HttpContext) {
	ctx.ContentType(c.CType)

	ctx.Write(c.Start)
	ctx.Flush()

	for i := 0; i < c.Len; i++ {
		s := <-c.Chan
		ctx.Write([]byte(s))
		ctx.Flush()
	}
	ctx.Write(c.End)
	ctx.Flush()
}
