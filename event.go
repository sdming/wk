// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"fmt"
)

// SubscriberFunc 
type SubscriberFunc func(*EventContext)

// On implement Subscriber
func (f SubscriberFunc) On(e *EventContext) {
	f(e)
}

// Subscriber is interface of event listener
type Subscriber interface {
	On(e *EventContext)
}

// SubscriberItem is wrap of Subscriber
type SubscriberItem struct {
	Moudle  string
	Name    string
	Handler Subscriber
}

// EventContext is the data pass to Subscriber
type EventContext struct {
	Moudle  string
	Name    string
	Source  interface{}
	Data    interface{}
	Context *HttpContext
}

// String 
func (e *EventContext) String() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("event: %s, %s", e.Moudle, e.Name)
}

var (
	Subscribers []*SubscriberItem
)

func init() {
	Subscribers = make([]*SubscriberItem, 0)
}

// On register event lister 
func On(moudle, name string, handler Subscriber) {

	sub := &SubscriberItem{
		Moudle:  moudle,
		Name:    name,
		Handler: handler,
	}
	Subscribers = append(Subscribers, sub)
}

// On register event lister 
func OnFunc(moudle, name string, handler func(*EventContext)) {
	On(moudle, name, SubscriberFunc(handler))
}
