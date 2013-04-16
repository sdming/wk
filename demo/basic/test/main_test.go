package test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/sdming/kiss/kson"
	"github.com/sdming/wk/demo/basic/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

var defaultAccept = "text/html, application/xhtml+xml, */*"
var baseUrl = "http://localhost:8080"

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

	//t.Log("expect:", expect)
	//t.Log("actual:", actual)

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
	url := "/data/top/1"
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

func BenchmarkControllerRangeCount(b *testing.B) {
	url := "/demo/rangecount/?start=1&end=99"
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

func demoData() *model.Data {
	return &model.Data{
		Int:   32,
		Uint:  1024,
		Str:   "string",
		Float: 1.1,
		Byte:  64,
	}
}

func dump(d *model.Data) {
	fmt.Printf("%#v \n", d)
}

func TestController(t *testing.T) {

	data := make([]*model.Data, 0)

	dotest(t, "GET", "/demo/clear/", nil, "0")

	dotest(t, "GET", "/demo/all/", nil, Json(data))

	dotest(t, "GET", "/demo/add/?int=32&str=string&uint=1024&float=1.1&byte=64", nil, demoData().String())
	data = append(data, demoData())
	dotest(t, "GET", "/demo/add/?int=32&str=string&uint=1024&float=1.1&byte=64", nil, demoData().String())
	data = append(data, demoData())

	dotest(t, "GET", "/demo/all/", nil, Json(data))

	dotest(t, "GET", "/demo/int/32", nil, Json(data))

	dotest(t, "GET", "/demo/rangecount/?start=1&end=99", nil, strconv.Itoa(len(data)))

	dotest(t, "GET", "/demo/put/32?str=s&uint=64&float=3.14&byte=8", nil, strconv.Itoa(len(data)))

	dotest(t, "GET", "/demo/delete/32", nil, strconv.Itoa(len(data)))

	dotest(t, "POST", "/demo/post/", []byte(Json(demoData())), "true")
	dotest(t, "POST", "/demo/post/", []byte(Json(demoData())), "true")

	dotest(t, "GET", "/demo/all/", nil, Json(data))
}
