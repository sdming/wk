// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic demo

*/
package main

import (
	"fmt"
	"github.com/sdming/wk"
)

func main() {
	server, err := wk.NewDefaultServer()

	if err != nil {
		fmt.Println("DefaultServer error", err)
		return
	}

	server.Start()

}
