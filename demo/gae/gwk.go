// +build appengine

package gwk

import (
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/controller"
	"github.com/sdming/wk/demo/basic/model"
	"io/ioutil"
	"log"
	"net/http"
)

var server *wk.HttpServer
var logger *log.Logger

func init() {

	logger = log.New(ioutil.Discard, "gwk", log.Ldate|log.Ltime)
	wk.Logger = logger
	server, _ := wk.NewHttpServer(wk.NewDefaultConfig())

	controller.RegisterBasicRoute(server)
	controller.RegisterUserRoute(server)
	controller.RegisterDocRoute(server)

	// data
	model.RegisterDataRoute(server)

	// event
	model.RegisterEventTrace(server)

	// compress
	server.Processes.InsertBefore("_render", wk.NewCompressProcess("compress_test", "*", "/compress/"))

	// file
	model.RegisterFileRoute(server)

	// pipe
	model.RegisterBigPipeRoute(server)

	// session
	controller.RegisterSessionRoute(server)

	// home
	controller.RegisterHomeRoute(server)

	server.Setup()

	http.Handle("/", server)
}
