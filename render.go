// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import ()

// RenderProcessor render http result
type RenderProcessor struct {
	server *HttpServer
}

// newRenderProcessor
func newRenderProcessor() *RenderProcessor {
	return &RenderProcessor{}
}

// Register
func (p *RenderProcessor) Register(server *HttpServer) {
	p.server = server
}

// Execute render http result to response
func (p *RenderProcessor) Execute(ctx *HttpContext) {

	if ctx.Result == nil {
		if ctx.Error != nil {
			ctx.Result = &ErrorResult{
				Err: ctx.Error,
			}
		}
	}

	if ctx.Result == nil {
		ctx.Result = resultNotFound
	}

	p.server.Fire(_render, _eventStartResultExecute, p, nil, ctx)
	ctx.Result.Execute(ctx)
	p.server.Fire(_render, _eventEndResultExecute, p, nil, ctx)
}
