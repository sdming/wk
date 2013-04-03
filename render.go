// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

// RenderProcessor rende http result to client
type RenderProcessor struct {
	server *HttpServer
}

func newRenderProcessor() *RenderProcessor {
	return &RenderProcessor{}
}

func (p *RenderProcessor) Register(server *HttpServer) {
	p.server = server
}

// Execute rende result to client
func (p *RenderProcessor) Execute(ctx *HttpContext) {
	if ctx.Result == nil {
		if ctx.Error != nil {
			ctx.Result = &ErrorResult{
				Err: ctx.Error,
			}
		}
	}

	p.server.Fire(_render, _eventResultExecuting, p, nil, ctx)
	ctx.Result.Execute(ctx)
	p.server.Fire(_render, _eventResultExecuted, p, nil, ctx)
}
