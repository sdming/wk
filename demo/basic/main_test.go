package main_test

import (
	"encoding/json"
	"encoding/xml"
	"github.com/sdming/kiss/kson"
	"github.com/sdming/wk/demo/basic/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var defaultAccept = "text/html, application/xhtml+xml, */*"
var baseUrl = "http://localhost:8080"

func curl(method, url string, form url.Values) (string, error) {
	var resp *http.Response
	var err error

	url = baseUrl + url

	if method == "GET" {
		resp, err = http.Get(url)
	} else if method == "POST" {
		resp, err = http.PostForm(url, form)
	} else if method == "DELETE" {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return "", err
		}
		client := &http.Client{}
		resp, err = client.Do(req)
	}

	if err != nil {
		resp.Body.Close()
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func dotest(t *testing.T, method, url string, form url.Values, expect string) {
	t.Log(method, url)

	actual, err := curl(method, url, form)
	if err != nil {
		t.Error("error: curl", method, url, err)
		return
	}

	//t.Log("expect:", expect)
	t.Log("actual:", actual)

	if strings.TrimSpace(expect) != strings.TrimSpace(actual) {
		t.Error("error:", method, url, "expect:", expect, "actual:", actual)
		return
	}
}

func Xml(x interface{}) string {
	b, err := xml.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func Json(x interface{}) string {
	b, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func Kson(x interface{}) string {
	b, err := kson.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func BenchmarkHandle(b *testing.B) {
	url := "/data/top/10"
	btestGet(b, url)
}

func BenchmarkDataByInt(b *testing.B) {

	url := "/data/1"
	btestGet(b, url)
}

func BenchmarkDataByIntJson(b *testing.B) {

	url := "/data/1/json"
	btestGet(b, url)
}

func btestGet(b *testing.B, url string) {
	for i := 0; i < b.N; i++ {
		_, err := curl("GET", url, nil)
		if err != nil {
			b.Error("error: curl", url, err)
			break
		}
	}
}

func TestRoute(t *testing.T) {
	dotest(t, "GET", "/data/top/10", nil, Json(model.DataTop(10)))
	dotest(t, "GET", "/data/int/1", nil, Json(model.DataByInt(1)))
	dotest(t, "GET", "/data/range/1-9", nil, Json(model.DataByIntRange(1, 9)))
	dotest(t, "GET", "/data/int/1/xml", nil, Xml(model.DataByInt(1)))
	dotest(t, "GET", "/data/int/1/json", nil, Json(model.DataByInt(1)))
	dotest(t, "GET", "/data/int/1/kson", nil, Kson(model.DataByInt(1)))
	dotest(t, "GET", "/data/name/1", nil, Json(model.DataByInt(1)))
	dotest(t, "GET", "/data/namerange/1-9", nil, Json(model.DataByIntRange(1, 9)))
	dotest(t, "GET", "/data/namerange/?start=1&end=9", nil, Json(model.DataByIntRange(1, 9)))

	dotest(t, "POST", "/data/post?",
		url.Values{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}},
		model.DataSet("string", 1024, 32, 1.1, 64).String())

	dotest(t, "POST", "/data/postptr?",
		url.Values{"str": {"string"}, "uint": {"1024"}, "int": {"32"}, "float": {"1.1"}, "byte": {"64"}},
		model.DataSet("string", 1024, 32, 1.1, 64).String())

	dotest(t, "DELETE", "/data/delete/1", nil, model.DataDelete(1))

	dotest(t, "GET", "/data/set?str=string&uint=1024&int=32&float=3.14&byte=64", nil,
		Json(model.DataSet("string", 1024, 32, 3.14, 64)))

}
