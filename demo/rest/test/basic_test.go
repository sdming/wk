package test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

var defaultAccept = "text/html, application/xhtml+xml, */*"
var baseUrl = "http://localhost:8080"

type Data struct {
	Str   string
	Uint  uint64
	Int   int
	Float float32
	Byte  byte
}

func (d Data) String() string {
	return fmt.Sprintf("string:%s;uint:%d;int:%d;float:%f;byte:%c",
		d.Str, d.Uint, d.Int, d.Float, d.Byte)
}

func curl(method, urlstr string, data interface{}) (string, error) {
	var resp *http.Response
	var err error
	client := &http.Client{}

	urlstr = baseUrl + urlstr

	if method == "GET" {
		resp, err = client.Get(urlstr)
	} else if method == "POST" {
		if form, ok := data.(url.Values); ok {
			resp, err = client.PostForm(urlstr, form)
		} else if buf, ok := data.([]byte); ok {
			resp, err = client.Post(urlstr, "text/plain", bytes.NewReader(buf))
		} else {
			err = errors.New("error post data ")
		}
	} else if method == "DELETE" {
		var req *http.Request
		req, err = http.NewRequest(method, urlstr, nil)
		if err != nil {
			return "", err
		}
		resp, err = client.Do(req)
	}

	if err != nil {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func dotest(t *testing.T, method, url string, data interface{}, expect string) {
	t.Log(method, url)

	actual, err := curl(method, url, data)
	if err != nil {
		t.Error("curl error:", err, method, url)
		return
	}

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

func BenchmarkControllerRangeCount(b *testing.B) {
	url := "/basic/rangecount/?start=1&end=99"
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

func demoData() Data {
	return Data{
		Int:   32,
		Uint:  1024,
		Str:   "string",
		Float: 1.1,
		Byte:  64,
	}
}

func TestController(t *testing.T) {

	data := make([]Data, 0)

	dotest(t, "GET", "/basic/clear/", nil, "0")

	dotest(t, "GET", "/basic/all/", nil, Json(data))

	dotest(t, "GET", "/basic/add/?int=32&str=string&uint=1024&float=1.1&byte=64", nil, demoData().String())
	data = append(data, demoData())
	dotest(t, "GET", "/basic/add/?int=32&str=string&uint=1024&float=1.1&byte=64", nil, demoData().String())
	data = append(data, demoData())

	dotest(t, "GET", "/basic/all/", nil, Json(data))

	dotest(t, "GET", "/basic/int/32", nil, Json(data))

	dotest(t, "GET", "/basic/rangecount/?start=1&end=99", nil, strconv.Itoa(len(data)))

	dotest(t, "GET", "/basic/set/32?str=s&uint=64&float=3.14&byte=8", nil, strconv.Itoa(len(data)))

	dotest(t, "GET", "/basic/delete/32", nil, strconv.Itoa(len(data)))

	dotest(t, "POST", "/basic/post/", []byte(Json(demoData())), "true")
	dotest(t, "POST", "/basic/post/", []byte(Json(demoData())), "true")

	dotest(t, "GET", "/basic/all/", nil, Json(data))
}
