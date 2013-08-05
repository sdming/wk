// +build appengine

package gwk

import (
	"github.com/sdming/wk"
	"github.com/sdming/wk/demo/basic/boot"
	_ "github.com/sdming/wk/demo/basic/controller"
	_ "github.com/sdming/wk/demo/basic/model"
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

	for _, fn := range boot.Inits {
		fn(server)
	}

	server.Setup()

	http.Handle("/", server)
}
