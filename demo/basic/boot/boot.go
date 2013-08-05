// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
boot
*/
package boot

import (
	"github.com/sdming/wk"
	"sync"
)

type ServerInitFunc func(*wk.HttpServer)

var Inits []ServerInitFunc = make([]ServerInitFunc, 0)
var lock sync.Mutex

func Boot(fn ServerInitFunc) {
	lock.Lock()
	defer lock.Unlock()

	Inits = append(Inits, fn)
}

func init() {

}
