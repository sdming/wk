// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package wk

import (
	"fmt"
)

type SubscriberFunc func(*EventContext)

// The SubscriberFunc type is an adapter to Subscriber  
func (f SubscriberFunc) On(e *EventContext) {
	f(e)
}

// Subscriber calls f(e).
type Subscriber interface {
	On(e *EventContext)
}

// SubscriberItem is 
type SubscriberItem struct {
	Moudle  string
	Name    string
	Handler Subscriber
}

// EventContext is 
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
	fmt.Println("Subscribe Event On", moudle, name)

	sub := &SubscriberItem{
		Moudle:  moudle,
		Name:    name,
		Handler: handler,
	}
	Subscribers = append(Subscribers, sub)
}

func OnFunc(moudle, name string, handler func(*EventContext)) {
	On(moudle, name, SubscriberFunc(handler))
}
