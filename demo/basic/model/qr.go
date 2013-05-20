package model

import (
	"code.google.com/p/rsc/qr"
	"errors"
	"github.com/sdming/wk"
	"log"
)

// http://research.swtch.com/qart
type QrCodeResult struct {
	Text string
}

//
func (qrcode *QrCodeResult) Execute(ctx *wk.HttpContext) error {

	c, err := qr.Encode(qrcode.Text, qr.M)
	if err != nil {
		log.Println("QrCodeResult Execute Error", err)
		return err
	}
	png := c.PNG()
	//ioutil.WriteFile(qrcode.Text+"_demo.png", png, 0666)
	ctx.ContentType("image/png")
	ctx.Write(png)

	return nil
}

func RegisterQrRoute(server *wk.HttpServer) {
	// url: get /qr/show/?text=hello
	// show qrcode png image
	server.RouteTable.Get("/qr/show/").To(showQrCode)
}

func showQrCode(ctx *wk.HttpContext) (result wk.HttpResult, err error) {
	text := ctx.FV("text")
	if text == "" {
		err = errors.New("text is invalid")
	} else {
		result = &QrCodeResult{
			Text: text,
		}
	}
	return
}
