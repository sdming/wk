package model

import (
	"bytes"
	"fmt"
	"github.com/sdming/wk"
	"log"
	"net/url"
	"strings"
	"time"
)

type EventTrace struct {
	Module    string
	Name      string
	Timestamp int64
}

func RegisterEventTrace(server *wk.HttpServer) {
	wk.OnFunc("*", "*", eventTraceFunc)
}

const (
	eventFlashKey string = ":eventTrace"
)

func eventTraceFunc(e *wk.EventContext) {
	var events []EventTrace

	if v, ok := e.Context.GetFlash(eventFlashKey); ok {
		events = v.([]EventTrace)
	} else {
		events = make([]EventTrace, 0)
	}

	events = append(events, EventTrace{
		Module:    e.Moudle,
		Name:      e.Name,
		Timestamp: time.Now().UnixNano(),
	})

	if e.Name == "end_request" {
		printTrace(e.Context.Request.URL, events)
		return
	}

	e.Context.SetFlash(eventFlashKey, events)

}

func printTrace(request *url.URL, trace []EventTrace) {

	if len(trace) == 0 {
		log.Panicln("")
		return
	}

	var buffer bytes.Buffer
	indent := 0
	s := "    "
	offset := trace[0].Timestamp

	buffer.WriteString("\n")
	buffer.WriteString(fmt.Sprintf("url:%v \n", request))
	buffer.WriteString(strings.Repeat("-", 10))
	buffer.WriteString("\n")
	for _, t := range trace {
		if strings.HasPrefix(t.Name, "end_") && indent > 0 {
			indent--
		}

		buffer.WriteString(strings.Repeat(s, indent))
		buffer.WriteString(fmt.Sprintf("%s\t %s\t %d \n", t.Module, t.Name, t.Timestamp-offset))

		if strings.HasPrefix(t.Name, "start_") {
			indent++
		}
	}
	buffer.WriteString(strings.Repeat("-", 10))
	buffer.WriteString("\n")

	log.Println(buffer.String())
}
