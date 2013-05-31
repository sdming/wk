package model

import (
	"bytes"
	"fmt"
	"github.com/sdming/wk"
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

//var re = regexp.MustCompile("^/doc/otherdemo")

func eventTraceFunc(e *wk.EventContext) {
	path := strings.ToLower(e.Context.RequestPath)
	if !strings.HasPrefix(path, "/doc/otherdemo") {
		return
	}

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
		printTrace(e.Context, events)
		return
	}

	e.Context.SetFlash(eventFlashKey, events)

}

func printTrace(ctx *wk.HttpContext, trace []EventTrace) {

	if len(trace) == 0 {
		return
	}

	var buffer bytes.Buffer
	indent := 0
	s := "    "
	offset := trace[0].Timestamp

	buffer.WriteString("\n")
	//buffer.WriteString(fmt.Sprintf("url:%v \n", ctx.Request.URL))
	//buffer.WriteString(strings.Repeat("-", 10))
	//buffer.WriteString("\n")

	buffer.WriteString("<script>\n")
	buffer.WriteString("var pageProfiler = \"")
	for _, t := range trace {
		if strings.HasPrefix(t.Name, "end_") && indent > 0 {
			indent--
		}

		buffer.WriteString(strings.Repeat(s, indent))
		//buffer.WriteString(fmt.Sprintf("%s\t %s\t %d \n", t.Module, t.Name, t.Timestamp-offset))
		buffer.WriteString(fmt.Sprintf("%s\t %s\t %d ns \\n", t.Module, t.Name, t.Timestamp-offset))

		if strings.HasPrefix(t.Name, "start_") {
			indent++
		}
	}
	//buffer.WriteString(strings.Repeat("-", 10))
	//buffer.WriteString("\n")
	buffer.WriteString("\";")
	buffer.WriteString("\n</script>\n")
	ctx.Write(buffer.Bytes())
}
