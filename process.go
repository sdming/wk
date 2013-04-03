// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"strings"
)

// HttpProcessor
type HttpProcessor interface {
	Execute(ctx *HttpContext)
	Register(server *HttpServer)
}

// HttpHandler item
type Process struct {
	// name
	Name string

	// path to match (prefix, / or empty to match all)    
	Path string

	// method to match (* or empty to match all)
	Method string //http method

	// handler  
	Handler HttpProcessor
}

// match
func (p *Process) match(ctx *HttpContext) bool {
	if !strings.HasPrefix(ctx.RequestPath, p.Path) {
		return false
	}
	if p.Method == "*" || p.Method == "" || strings.Contains(p.Method, ctx.Method) {
		return true
	}
	return false
}

// ProcessTable
type ProcessTable []*Process

// Append add a process at end
func (pt ProcessTable) Append(p *Process) {
	pt = append(pt, p)
}

// Remove delet a process by name
func (pt ProcessTable) Remove(name string) {
	for i := 0; i < len(pt); i++ {
		if pt[i].Name == name {
			pt = append(pt[:i], pt[i+1:]...)
		}
	}
}

var Processes ProcessTable

// init ProcessTable
func init() {
	Processes = make([]*Process, 0)

	RegisterProcessor(_static, newStaticProcessor())
	RegisterProcessor(_route, newRouteProcessor())
	RegisterProcessor(_render, newRenderProcessor())

}

// RegisterProcessor
func RegisterProcessor(name string, p HttpProcessor) {
	process := &Process{
		Name:    name,
		Path:    _root,
		Method:  _any,
		Handler: p,
	}

	Processes.Append(process)
}
