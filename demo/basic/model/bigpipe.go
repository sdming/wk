// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*

*/
package model

import (
	"fmt"
	"github.com/sdming/wk"
	"math/rand"
	"time"
)

func RegisterBigPipeRoute(server *wk.HttpServer) {
	// url: get /bugpipe/test.html
	server.RouteTable.Get("/bugpipe/test.html").To(BigPipe)

}

func BigPipe(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	start := `<!DOCTYPE html>
    <head>       
    </head>
    <body>
       
        <div id="panel">

        </div>
        <script>
        function p(s) {
            var panel = document.getElementById('panel');
            var n = document.createElement('div');
            n.innerHTML = s + " create at " + new Date();
            panel.appendChild(n);
        }
        </script>	
    `

	end := `
        </body>
	</html>`

	l := 5

	r := &wk.ChanResult{
		Len:   l,
		CType: "text/html",
		Chan:  make(chan string, l),
		Start: []byte(start),
		End:   []byte(end),
	}

	for i := 0; i < l; i++ {
		go bigpipeOutput(r.Chan)
	}

	return r, nil

}

func bigpipeOutput(c chan string) {
	d := 5 * (2 + rand.Intn(8))
	time.Sleep(time.Duration(d) * time.Second)
	c <- fmt.Sprintf(`
		<script>
            p("delay %d")
        </script>
        `, d)
}
