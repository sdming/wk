package wk_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
)

type User struct {
	Name  string
	Age   int
	Web   string
	Email string
}

func defaultUser() User {
	return User{
		Name:  "Gopher",
		Age:   3,
		Web:   "http://golang.org",
		Email: "gopher@golang.org",
	}
}

func readFile(file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Sprintf("read %s error; %s", file, err.Error())
	}
	return string(b)
}

func toxml(x interface{}) string {
	b, err := xml.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func tojson(x interface{}) string {
	b, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func success(t *testing.T, message string, err error) {
	t.Log(message, err)

	if err != nil {
		t.Error(message, "error", err)
		return
	}
}

func successAndEq(t *testing.T, message string, err error, expect, actual interface{}) {
	if err != nil {
		success(t, message, err)
	} else {
		equal(t, message, expect, actual)
	}
}

func removeSpace(text string) string {
	re := regexp.MustCompile(`\s`)
	return re.ReplaceAllString(text, "")
}

func equal(t *testing.T, message string, expect, actual interface{}) {
	if expect != actual {
		t.Errorf("%s Equal fail, expect %v, actual %v ", message, expect, actual)
	}
	return
}
