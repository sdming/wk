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
func (c *ChanResult) Execute(ctx *HttpContext) error {
	ctx.ContentType(c.ContentType)

	ctx.Write(c.Start)
	ctx.Flush()

	if c.Timeout < time.Millisecond {
		c.Timeout = defaultChanResultTimeout
	}

	waitchan := make(chan bool)
	donechan := make(chan bool)

	go func() {
		for s := range c.Chan {
			ctx.Write([]byte(s))
			ctx.Flush()
		}
		donechan <- true
	}()

	go func() {
		c.Wait.Wait()
		close(c.Chan)
		waitchan <- true
	}()

	select {
	case <-waitchan:
	case <-time.After(c.Timeout):
	}

	<-donechan
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
