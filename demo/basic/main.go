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
	"github.com/sdming/wk/demo/basic/boot"
	_ "github.com/sdming/wk/demo/basic/controller"
	_ "github.com/sdming/wk/demo/basic/model"
)

func main() {

	server, err := wk.NewDefaultServer()

	if err != nil {
		fmt.Println("NewDefaultServer error", err)
		return
	}

	for _, fn := range boot.Inits {
		fn(server)
	}

	server.Start()

}
