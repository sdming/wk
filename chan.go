// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"sync"
	"time"
)

var (
	defaultChanResultTimeout = 60 * time.Second
)

// ChanResult
// TODO: just a demo, need to enhance
// TODO: doesn't work on chrome
type ChanResult struct {
	Wait        sync.WaitGroup
	Chan        chan string
	ContentType string
	Start       []byte
	End         []byte
	Timeout     time.Duration
}

// Execute read string from chan and write to response
// TODO: enhance it
func (c *ChanResult) Execute(ctx *HttpContext) error {
	ctx.ContentType(c.ContentType)

	ctx.Write(c.Start)
	ctx.Flush()

	if c.Timeout < time.Millisecond {
		c.Timeout = defaultChanResultTimeout
	}

	var waitchan chan bool = make(chan bool)

	go func() {
		for s := range c.Chan {
			ctx.Write([]byte(s))
			ctx.Flush()
		}
	}()

	go func() {
		c.Wait.Wait()
		close(c.Chan)
		waitchan <- true
		//close(waitchan)
	}()

	select {
	case <-waitchan:
	case <-time.After(c.Timeout):
	}

	close(waitchan)

	ctx.Write(c.End)
	//ctx.Flush()

	return nil
}

// http://dave.cheney.net/2013/04/30/curious-channels
func waitMany(a, b chan bool) {
	for a != nil || b != nil {
		select {
		case <-a:
			a = nil
		case <-b:
			b = nil
		}
	}
}
