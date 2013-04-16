// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"path"
)

// StaticProcessor handle request of static file
type StaticProcessor struct {
	server *HttpServer
}

// newStaticProcessor return default  *StaticProcessor
func newStaticProcessor() *StaticProcessor {
	return &StaticProcessor{}
}

// Register
func (p *StaticProcessor) Register(server *HttpServer) {
	p.server = server
}

// Execute set FileResult if request file does exist
func (p *StaticProcessor) Execute(ctx *HttpContext) {
	if ctx.Result != nil || ctx.Error != nil {
		return
	}

	physicalPath := path.Join(p.server.Config.PublicDir, ctx.RequestPath)
	if !isFileExists(physicalPath) {
		return
	}

	ctx.PhysicalPath = physicalPath
	ctx.Result = File(physicalPath)
}
