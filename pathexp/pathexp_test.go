// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package pathexp_test

import (
	"github.com/sdming/wk/pathexp"
	"strings"
	"testing"
)

func TestMatch(t *testing.T) {
	var test [][4]string = [][4]string{

		// /*   /              */
		{"/", "/any", "true", ""},
		{"/", "/any.html", "true", ""},
		{"/any", "/any", "true", ""},
		{"/any", "/any/", "true", ""},
		{"/any", "/any.html", "true", ""},
		{"/any", "/any/1", "true", ""},
		{"/any/", "/any", "false", ""},
		{"/any/", "/any/", "true", ""},
		{"/any/", "/any/1", "true", ""},
		{"/any/", "/any.html", "false", ""},
		{"/any/1/2/3", "/any/1/2/3", "true", ""},
		{"/any/1/2/3/", "/any/1/2/3", "false", ""},
		{"/any/", "/anyhoho/", "false", ""},

		// /*  /query.html     */
		{"/query.html", "/query", "false", ""},
		{"/query.html", "/query/", "false", ""},
		{"/query.html", "/query/any", "false", ""},
		{"/query.html", "/query.html", "true", ""},
		{"/query.html", "/queryany.html", "false", ""},

		// /*  /{query}      */
		{"/{query}", "/query", "true", "query"},
		{"/{query}", "/query/", "true", "query"},
		{"/{query}", "/query.html", "true", "query"},

		// /*  /{query}/      */
		{"/{query}/", "/query", "false", ""},
		{"/{query}/", "/query/", "true", "query"},
		{"/{query}/", "/query.html", "false", ""},

		// /*  /query/hoho     */
		{"/query/hoho", "/query", "false", ""},
		{"/query/hoho", "/query/", "false", ""},
		{"/query/hoho", "/query/any", "false", ""},
		{"/query/hoho", "/query/hoho", "true", ""},
		{"/query/hoho", "/query/hoho.html", "true", ""},
		{"/query/hoho", "/queryany/hohoany", "false", ""},

		// /*  /query/{type}     */
		{"/query/{type}", "/query", "false", ""},
		{"/query/{type}", "/query/", "true", ""},
		{"/query/{type}", "/query/type", "true", "type"},
		{"/query/{type}", "/query/type/", "true", "type"},
		{"/query/{type}", "/query/type/any", "true", "type"},
		{"/query/{type}", "/query/type.html", "true", "type"},
		{"/query/{type}", "/query.html", "false", ""},
		{"/query/{type}", "/queryhoho", "false", ""},
		{"/query/{type}", "/queryhoho/type", "false", ""},

		// /*  /query/{type}/     */
		{"/query/{type}/", "/query", "false", ""},
		{"/query/{type}/", "/query/", "false", ""},
		{"/query/{type}/", "/query/type", "false", ""},
		{"/query/{type}/", "/query/type/", "true", "type"},
		{"/query/{type}/", "/query/type/any", "true", "type"},
		{"/query/{type}/", "/query/type.html", "false", ""},
		{"/query/{type}/", "/query.html", "false", ""},
		{"/query/{type}/", "/queryhoho/", "false", ""},
		{"/query/{type}/", "/queryhoho/type/", "false", ""},

		// /*  /query/{type}/{year}/{month}/{day}     */
		{"/query/{type}/{year}/{month}/{day}", "/query", "false", ""},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year", "false", ""},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month", "false", ""},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day/", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day/any", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/queryhoho/type/year/month/day", "false", ""},
		{"/query/{type}/{year}/{month}/{day}", "/queryhoho/type", "false", ""},

		// /*  /query/{type}/{year}/{month}/{day}/     */
		{"/query/{type}/{year}/{month}/{day}/", "/query", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day/", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day/any", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day.html", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/queryhoho/type/year/month/day", "false", ""},
		{"/query/{type}/{year}/{month}/{day}", "/queryhoho/type", "false", ""},

		// /*  /query/{type}/{year}-{month}-{day}    */
		{"/query/{type}/{year}-{month}-{day}", "/query", "false", ""},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year", "false", ""},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month", "false", ""},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day/", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day/any", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/queryhoho/type-year-month-day", "false", ""},
		{"/query/{type}/{year}-{month}-{day}", "/query/typehoho", "false", ""},
		{"/query/{type}/{year}-{month}-{day}", "/query/typehoho/year", "false", ""},

		// /*  /query/{type}/{year}-{month}-{day}/     */
		{"/query/{type}/{year}-{month}-{day}/", "/query", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year-month", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year-month-day", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year-month-day/", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year-month-day/any", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year-month-day.html", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/", "/query/typehoho/", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/typehoho/year", "false", ""},

		// /*  /query/{type}/{year}/{month}/{day}.html    */
		{"/query/{type}/{year}/{month}/{day}.html", "/query", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day/", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day/any", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}.html", "/queryhoho/type/year/month/day.html", "false", ""},

		// /*  /query/{type}/{year}/{month}/{day}/detail.html    */
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day/", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day/any", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day/detail.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/queryhoho/type/year/month/day/detail.html", "false", ""},

		// /*  /query/{type}/{year}-{month}-{day}.html    */
		{"/query/{type}/{year}-{month}-{day}.html", "/query", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day/", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day/any", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}.html", "/queryhoho/type/year-month-day.html", "false", ""},

		// /*  /query/{type}/{year}-{month}-{day}/detail.html    */
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month-day", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month/day-", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month-day/any", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month-day/detail.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/queryhoho/type/year-month-day/detail.html", "false", ""},

		/* page & orderby */
		{"/query/{type}/{year}-{month}-{day}-1", "/query/type/year-month-day", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1", "/query/type/year-month-day--any", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1", "/query/type/year-month-day-1", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}-1/", "/query/type/year-month-day-1", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1/", "/query/type/year-month-day-1/", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}-1.html", "/query/type/year-month-day-1.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}/1", "/query/type/year-month-day", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/1", "/query/type/year-month-day/any", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/1", "/query/type/year-month-day/1", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}/1/", "/query/type/year-month-day/1", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/1/", "/query/type/year-month-day/1/", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}/1.html", "/query/type/year-month-day/1.html", "true", "type,year,month,day"},

		{"/query/{type}/{year}-{month}-{day}-1-{orderby}", "/query/type/year-month-day-1", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}", "/query/type/year-month-day-1-", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}", "/query/type/year-month-day-1-orderby", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}", "/query/type/year-month-day-1-orderby.html", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}/", "/query/type/year-month-day-1-orderby", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}/", "/query/type/year-month-day-1-orderby/", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}/", "/query/type/year-month-day-1-orderby.html", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}.html", "/query/type/year-month-day-1-orderby", "false", ""},

		{"/query/{type}/{year}-{month}-{day}-1/{orderby}", "/query/type/year-month-day-1", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}", "/query/type/year-month-day-1/", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}", "/query/type/year-month-day-1/orderby", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}", "/query/type/year-month-day-1/orderby.html", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}/", "/query/type/year-month-day-1/orderby", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}/", "/query/type/year-month-day-1/orderby/", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}/", "/query/type/year-month-day-1/orderby.html", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}.html", "/query/type/year-month-day-1/orderby", "false", ""},

		// /*  /query#/{type}/{any}    */
		{"/query#/{type}/{any}", "/query", "true", ""},
		{"/query#/{type}/{any}", "/query/type", "true", "type"},
		{"/query#/{type}/{any}", "/query/type/", "true", "type"},
		{"/query#/{type}/{any}", "/query/type/any", "true", "type,any"},
		{"/query#/{type}/{any}", "/query/type.html", "true", "type"},
		{"/query#/{type}/{any}", "/query/type/any.html", "true", "type,any"},
		{"/query#/{type}/{any}/", "/query.html", "true", ""},
		{"/query#/{type}/{any}/", "/queryhoho/", "true", ""},
		{"/query#/{type}/{any}/", "/queryhoho/type/", "true", ""},

		{"/", "/", "true", ""},
	}

	// if len(test) > 0 {
	// 	return
	// }

	for _, each := range test {

		pattern := each[0]
		input := each[1]
		expect := each[2]
		names := each[3]

		re, err := pathexp.Compile(pattern)
		if err != nil {
			t.Errorf("pattern compile fail: %s; pattern=%s \n", err, pattern)
			continue
		}

		matched, data := re.Match(input)
		if (!matched && expect == "true") || (matched && expect == "false") {
			t.Errorf("Match fail: expect:%s; actual=%v; pattern=%s; input=%s \n", expect, matched, pattern, input)
			continue
		}

		if (!matched && names != "") || ((matched && len(data) > 0) && names == "") {
			t.Errorf("Match sub names config fail: expect:%v; actual=%v; pattern=%s; input=%s \n", names == "", matched && len(data) > 0, pattern, input)
			continue
		}

		if names == "" {
			continue
		}

		namesSlice := strings.Split(names, ",")
		if !strings.Contains(pattern, "#") && len(namesSlice) != len(data) {
			t.Errorf("Match sub names length fail: expect:%v; actual=%v; pattern=%s; input=%s \n", len(namesSlice), len(data), pattern, input)
			continue
		}

		if len(namesSlice) == 0 {
			continue
		}

		for i, p := range data {
			if p[0] != namesSlice[i] {
				t.Errorf("Match sub names check name fail: expect:%s; actual=%s; pattern=%s; input=%s \n", namesSlice[0], p[0], pattern, input)
			} else if p[0] != p[1] {
				t.Errorf("Match sub names check value fail: expect:%s; actual=%s; pattern=%s; input=%s \n", p[0], p[1], pattern, input)
			}
		}

	}

}

func TestCompileSuccess(t *testing.T) {

	var success = []string{
		"/{query}",
		"/{query}/",
		"/query",
		"/query/",
		"/query/{type}",
		"/query/{type}/",
		"/query/{type}/{year}/{month}/{day}",
		"/query/{type}/{year}/{month}/{day}/",
		"/query/{type}/{year}-{month}-{day}",
		"/query/{type}/{year}-{month}-{day}/",
		"/query/{type}/{year}/{month}/{day}.html",
		"/query/{type}/{year}/{month}/{day}/detail.html",
		"/query/{type}/{year}-{month}-{day}.html",
		"/query/{type}/{year}-{month}-{day}/detail.html",
		"/query/{type}/{year}-{month}-{day}-{page}",
		"/query/{type}/{year}-{month}-{day}/{page}",
		"/query/{type}/{year}-{month}-{day}-{page}-{orderby}",
		"/query/{type}/{year}-{month}-{day}/{page}/{orderby}",

		"/query/{type}/{year}-{month}-{day}-1",
		"/query/{type}/{year}-{month}-{day}-1/",
		"/query/{type}/{year}-{month}-{day}-1.html",
		"/query/{type}/{year}-{month}-{day}/1",
		"/query/{type}/{year}-{month}-{day}/1/",
		"/query/{type}/{year}-{month}-{day}/1.html",

		"/query/{type}/{year}-{month}-{day}-1-{orderby}",
		"/query/{type}/{year}-{month}-{day}-1-{orderby}/",
		"/query/{type}/{year}-{month}-{day}-1-{orderby}.html",
		"/query/{type}/{year}-{month}-{day}/1/{orderby}",
		"/query/{type}/{year}-{month}-{day}/1/{orderby}/",
		"/query/{type}/{year}-{month}-{day}/1/{orderby}.html",

		"/query/type?",
		"/query#/{type}",

		"/{_controller}/{_action}",
		"/{_controller}/{_action}/",
		"/{_controller}/{_action}.html",
		"/"}

	for _, s := range success {
		_, err := pathexp.Compile(s)
		if err != nil {
			t.Errorf("pattern compile fail: %s; pattern=%s \n", err, s)
		}
	}
}

func TestCompileFail(t *testing.T) {

	var fail = []string{
		".",
		"query",
		"/{query",
		"/query/type}",
		"/query/{type}/{year}{month}{day}",
		"/query/{{type}",
		"/query/{type}}",
		"/{query@}",
		"/{?query}",
		"/query/{type}.html?{month}",
		"/query/{type}.{month}html",
		"/query/?1",
		"/query#/?",
		"/query{type#}/",
		"",
	}

	for _, s := range fail {
		_, err := pathexp.Compile(s)
		if err == nil {
			t.Errorf("pattern should be invalid: pattern=%s \n", s)
		}
	}
}

func BenchmarkCompile(b *testing.B) {

	pattern := "/query/hoho/{type}/{year}-{month}-{day}/"

	for i := 0; i < b.N; i++ {
		re, err := pathexp.Compile(pattern)
		if err != nil {
			b.Errorf("pattern compile fail: %s; pattern=%s \n", err, pattern, re)
			break
		}
	}

}

func BenchmarkUser(b *testing.B) {
	b.StopTimer()
	pattern := "/user/{action}/{arg}/"
	input := "/user/view/1/"

	re, err := pathexp.Compile(pattern)
	if err != nil {
		b.Errorf("pattern compile fail: %s; pattern=%s \n", err, pattern)
		return
	}
	//matched, data := re.Match(input)
	// b.Log(matched, data)
	// if !matched || data["action"] != "view" || data["arg"] != "1" {
	// 	b.Error("Match fail", matched, data)
	// 	return
	// }
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(input)
	}

}

func BenchmarkUserRegexp(b *testing.B) {
	b.StopTimer()
	//pattern := "/user/{action}/{arg}/"
	//pattern := "^/user(/(?P<action>[[:alnum:]]+)/(?P<arg>[[:alnum:]]+))?"
	pattern := "^/user/(?P<action>[[:alnum:]]+)/(?P<arg>[[:alnum:]]+)"
	input := "/user/view/1"

	re, err := pathexp.RegexpCompile(pattern)
	if err != nil {
		b.Errorf("pattern compile fail: %s; pattern=%s \n", err, pattern)
		return
	}
	// matched, data := re.Match(input)
	// b.Log(matched, data)
	// if !matched || data["action"] != "view" || data["arg"] != "1" {
	// 	b.Error("Match fail", matched, data)
	// 	return
	// }

	//var m pathexp.Matcher = re

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(input)
	}

}
