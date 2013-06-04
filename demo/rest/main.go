// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
rest web api demo

*/
package main

import (
	"fmt"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/rest/controller"
)

func main() {

	server, err := wk.NewDefaultServer()
	if err != nil {
		fmt.Println("NewDefaultServer error", err)
		return
	}

	server.Processes.Remove("_static")
	server.Config.ViewEnable = false
	server.Config.SessionEnable = false

	controller := controller.NewBasicController()

	// url: /demo/xxx/xxx
	// route to controller
	server.RouteTable.Path("/basic/{action}/{id}").ToController(controller)

	server.Start()

}
