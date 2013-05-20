// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"strings"
)

// HttpProcessor is interface that handle http request
type HttpProcessor interface {
	// Execute handle request
	Execute(ctx *HttpContext)

	// Register is called once when create a new HttpServer
	Register(server *HttpServer)
}

// Process is wrap of HttpProcessor
type Process struct {
	// Name
	Name string

	// Path is url to match, / or empty to match all
	// change to regex? containers muti ?
	Path string

	// Method is http method to match, * or empty to match all
	Method string //http method

	// Handler is the HttpProcessor
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

// ProcessTable is alias of []*Process
type ProcessTable []*Process

// Append add a *Process at end
func (pt *ProcessTable) Append(p *Process) {
	*pt = append(*pt, p)
}

// InsertBefore add a *Process before name
func (pt *ProcessTable) InsertBefore(name string, p *Process) {
	for i := 0; i < len(*pt); {
		if (*pt)[i].Name == name {
			*pt = append(*pt, p)
			copy((*pt)[i+1:], (*pt)[i:])
			(*pt)[i] = p
			return
		} else {
			i++
		}
	}
}

// InsertAfter add a *Process after name
func (pt *ProcessTable) InsertAfter(name string, p *Process) {
	for i := 0; i < len(*pt); {
		if (*pt)[i].Name == name {
			*pt = append(*pt, p)
			copy((*pt)[i+1:], (*pt)[i+1:])
			(*pt)[i+1] = p
			return
		} else {
			i++
		}
	}
}

// Remove delete a *Process from ProcessTable
func (pt *ProcessTable) Remove(name string) {
	for i := 0; i < len(*pt); {
		if (*pt)[i].Name == name {
			//*pt = append((*pt)[:i], (*pt)[i+1:]...)
			copy((*pt)[i:], (*pt)[i+1:])
			(*pt)[len((*pt))-1] = nil
			(*pt) = (*pt)[:len((*pt))-1]
		} else {
			i++
		}
	}
}

// Processes is global ProcessTable configration
var Processes ProcessTable

// init Processes
func init() {
	Processes = make([]*Process, 0, 11)
	RegisterProcessor(_static, newStaticProcessor())
	RegisterProcessor(_route, newRouteProcessor())
	RegisterProcessor(_render, newRenderProcessor())

}

// RegisterProcessor append a HttpProcessor to global ProcessTable
func RegisterProcessor(name string, p HttpProcessor) {
	process := &Process{
		Name:    name,
		Path:    _root,
		Method:  _any,
		Handler: p,
	}

	Processes.Append(process)
}
