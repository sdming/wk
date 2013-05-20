package wk_test

import (
	"bytes"
	"github.com/sdming/wk"
	"mime"
	"net/http"
	"path/filepath"
	"testing"
)

func TestDataResult(t *testing.T) {
	a := &wk.DataResult{}
	var result wk.HttpResult = a
	t.Log("NotImplemented", result)
}

func TestContentResult(t *testing.T) {
	var result wk.HttpResult = &wk.ContentResult{}
	t.Log("NotImplemented", result)
}

func TestVoidResult(t *testing.T) {
	var result wk.HttpResult = &wk.VoidResult{}
	t.Log("NotImplemented", result)
}

func TestRedirectResult(t *testing.T) {
	var result wk.HttpResult = wk.Redirect(`http://golanr.org`, true)
	t.Log("NotImplemented", result)
}

func TestErrorResult(t *testing.T) {
	var result wk.HttpResult = &wk.ErrorResult{}
	t.Log("NotImplemented", result)
}

func TestNotModifiedResult(t *testing.T) {
	var result wk.HttpResult = &wk.NotModifiedResult{}
	t.Log("NotImplemented", result)
}

func TestNotFoundResult(t *testing.T) {
	var result wk.HttpResult = &wk.NotFoundResult{}
	t.Log("NotImplemented", result)
}

func TestFileResult(t *testing.T) {
	f := &wk.FileResult{}
	var result wk.HttpResult = f
	t.Log("NotImplemented", result)

	var ctype wk.ContentTyper = f
	t.Log("content:", ctype.Type())
}

func TestFileStreamResult(t *testing.T) {
	f := &wk.FileResult{}
	var result wk.HttpResult = f
	t.Log("NotImplemented", result)

	var ctype wk.ContentTyper = f
	t.Log("content:", ctype.Type())
}

func TestJsonResult(t *testing.T) {
	data := defaultUser()
	j := wk.Json(data)
	var result wk.HttpResult = j
	t.Log("JsonResult", result)

	header := make(http.Header)
	body := &bytes.Buffer{}
	var render wk.Render = j

	err := render.Write(header, body)
	success(t, "json write", err)
	equal(t, "json header", wk.ContentTypeJson, header.Get(wk.HeaderContentType))
	equal(t, "json body", tojson(data), body.String())
}

func TestXmlResult(t *testing.T) {
	data := defaultUser()
	x := wk.Xml(data)
	var result wk.HttpResult = x
	t.Log("XmlResult", result)

	header := make(http.Header)
	body := &bytes.Buffer{}
	var render wk.Render = x

	err := render.Write(header, body)
	success(t, "xml write", err)
	equal(t, "xml header", wk.ContentTypeXml, header.Get(wk.HeaderContentType))
	equal(t, "xml body", toxml(data), body.String())
}

func TestViewResult(t *testing.T) {
	data := make(wk.ViewData)
	data["user"] = defaultUser()

	var view *wk.ViewResult = wk.View("basic.html", data)
	var result wk.HttpResult = view
	t.Log("ViewResult", result)

	wk.DefaultViewEngine = defaultViewEngine()
	t.Log("ViewEngine", wk.DefaultViewEngine)
	expect := defaultRootPath() + "views/" + "output-basic.html"

	header := make(http.Header)
	body := &bytes.Buffer{}
	var render wk.Render = view

	err := render.Write(header, body)
	success(t, "view write", err)
	equal(t, "view header", mime.TypeByExtension(filepath.Ext(view.File)), header.Get(wk.HeaderContentType))
	//t.Log(readFile(expect))
	//t.Log(body.String())
	equal(t, "view body", removeSpace(readFile(expect)), removeSpace(body.String()))
}
