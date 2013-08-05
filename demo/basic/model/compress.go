// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*

*/
package model

import (
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/boot"
)

func init() {
	boot.Boot(RegisterCompress)
}

func RegisterCompress(server *wk.HttpServer) {
	server.Processes.InsertBefore("_render", wk.NewCompressProcess("compress_test", "*", "/compress/"))
}
