// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
rest web api demo

*/
package main

import (
	"fmt"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/controller"
)

func main() {

	server, err := wk.NewDefaultServer()
	if err != nil {
		fmt.Println("NewDefaultServer error", err)
		return
	}

	server.Processes.Remove("_static")

	controller := controller.NewDemoController()

	// url: /demo/xxx/xxx
	// route to controller
	server.RouteTable.Path("/demo/{action}/{id}").ToController(controller)

	server.Start()

}
