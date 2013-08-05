// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*

*/
package model

import (
	"fmt"
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/boot"
	"math/rand"
	"sync"
	"time"
)

func init() {
	boot.Boot(RegisterBigPipeRoute)
}

func RegisterBigPipeRoute(server *wk.HttpServer) {
	// url: get /bigpipe/test.html
	server.RouteTable.Get("/bigpipe/test.html").To(BigPipe)

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
	            n.innerHTML = s + " -- at " + new Date().toTimeString();
	            panel.appendChild(n);
	        }
        </script>	
        
`

	end := `

	<script>
    	p("end")
    </script>

    </body>
</html>
`

	l := 5

	cr := &wk.ChanResult{
		Wait:        sync.WaitGroup{},
		Chan:        make(chan string, l),
		ContentType: "text/html",
		Start:       []byte(start),
		End:         []byte(end),
	}

	for i := 0; i < l; i++ {
		cr.Wait.Add(1)
		go func(index int) {
			defer cr.Wait.Done()

			d := 3 + rand.Intn(10)
			time.Sleep(time.Duration(d) * time.Second)

			cr.Chan <- fmt.Sprintf(`
		<script>
	        p("goroutine %d delay %d")
	    </script>
`, index, d)
		}(i)
	}

	return cr, nil

}

func bigpipeOutput(r *wk.ChanResult) {
	defer r.Wait.Done()

	d := 5 * (2 + rand.Intn(8))
	time.Sleep(time.Duration(d) * time.Second)
	r.Chan <- fmt.Sprintf(`
		<script>
            p("delay %d")
        </script>
        `, d)
}
