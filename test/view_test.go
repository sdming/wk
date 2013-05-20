package wk_test

import (
	"bytes"
	"github.com/sdming/wk"
	"os"
	"path"
	"strings"
	"testing"
)

func defaultRootPath() string {
	p := ""
	pwd, err := os.Getwd()
	if err == nil {
		p = pwd
	} else {
		p = path.Dir(os.Args[0])
	}

	if os.PathSeparator == '\\' {
		p = strings.Replace(p, `\`, `/`, -1)
	}
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}
	return p
}

func defaultViewEngine() wk.ViewEngine {
	veBase := defaultRootPath() + "views/"
	ve, _ := wk.NewGoHtml(veBase)
	return ve
}

func testTemplate(t *testing.T, template, expect string, data interface{}) {
	ve := defaultViewEngine()
	veBase := defaultRootPath() + "views/"
	expect = veBase + expect

	var err error

	buffer := new(bytes.Buffer)

	t.Log("template", template)
	t.Log("expect", expect)
	t.Log("data", data)

	err = ve.Execte(buffer, template, data)
	if err != nil {
		t.Error("execute error:", err)
		return
	}

	expectText := readFile(expect)
	if removeSpace(expectText) != removeSpace(buffer.String()) {
		t.Error("execute output error", "expect", expectText, "actual", buffer.String())
	}
}

func TestBasic(t *testing.T) {
	var user User = User{
		Name:  "Gopher",
		Age:   3,
		Web:   "http://golang.org",
		Email: "gopher@golang.org",
	}
	data := make(wk.ViewData)
	data["user"] = user

	testTemplate(t, `basic.html`, `output-basic.html`, data)

}
