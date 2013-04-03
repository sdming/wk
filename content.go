// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

// // ContentResult is raw content
// type ContentResult struct {
// 	// modify time
// 	ModifyTime time.Time

// 	// content type
// 	ContentType string

// 	// content
// 	Data interface{}
// }

// // Execute render raw content
// func (c *ContentResult) Execute(ctx *HttpContext) {
// 	if r, ok := c.Data.(io.ReadSeeker); ok {
// 		http.ServeContent(ctx.Resonse, ctx.Request, c.ContentType, c.ModifyTime, r)
// 		return
// 	}
// 	fmt.Fprintln(ctx.Resonse, c.Data)
// }
