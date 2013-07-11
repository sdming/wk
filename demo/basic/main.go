// +build !appengine

// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic demo`

*/
package main

import (
	"fmt"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/controller"
	"github.com/sdming/wk/demo/basic/model"
)

func main() {

	server, err := wk.NewDefaultServer()

	if err != nil {
		fmt.Println("NewDefaultServer error", err)
		return
	}

	controller.RegisterBasicRoute(server)
	controller.RegisterUserRoute(server)
	controller.RegisterDocRoute(server)
	controller.RegisterConfigRoute(server)

	model.RegisterDataRoute(server)

	//demo, show to define custome httpresult
	if enableQrCode := true; enableQrCode {
		model.RegisterQrRoute(server)
	}

	if enableEventTrace := true; enableEventTrace {
		model.RegisterEventTrace(server)
	}

	if enableCompress := true; enableCompress {
		server.Processes.InsertBefore("_render", wk.NewCompressProcess("compress_test", "*", "/compress/"))
	}

	if enableFile := true; enableFile {
		model.RegisterFileRoute(server)
	}

	if enableBigpipe := true; enableBigpipe {
		model.RegisterBigPipeRoute(server)
	}

	if debugSession := true; debugSession {
		controller.RegisterSessionRoute(server)
	}

	controller.RegisterHomeRoute(server)

	server.Start()

}
