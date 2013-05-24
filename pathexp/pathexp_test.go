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
		{"/any", "/any.html?&a=1", "true", ""},
		{"/any", "/any/1", "true", ""},
		{"/any", "/any/?1", "true", ""},
		{"/any", "/any?1", "true", ""},
		{"/any/", "/any", "false", ""},
		{"/any/", "/any/", "true", ""},
		{"/any/", "/any/1", "true", ""},
		{"/any/", "/any?1", "false", ""},
		{"/any/", "/any.html", "false", ""},
		{"/any/1/2/3", "/any/1/2/3", "true", ""},
		{"/any/1/2/3/", "/any/1/2/3", "false", ""},
		{"/any/", "/anyhoho/", "false", ""},
		//{"/any", "/anyhoho", "false", ""},     /*-------------------------------------*/

		// /*  /query.html     */
		{"/query.html", "/query", "false", ""},
		{"/query.html", "/query/", "false", ""},
		{"/query.html", "/query/any", "false", ""},
		{"/query.html", "/query/?any", "false", ""},
		{"/query.html", "/query.html", "true", ""},
		{"/query.html", "/query.html?any", "true", ""},
		{"/query.html", "/query.html?&any=1", "true", ""},
		{"/query.html", "/queryany.html", "false", ""},

		// /*  /{query}      */
		{"/{query}", "/query", "true", "query"},
		{"/{query}", "/query/", "true", "query"},
		{"/{query}", "/query/?any", "true", "query"},
		{"/{query}", "/query.html", "true", "query"},
		{"/{query}", "/query.html?any", "true", "query"},
		{"/{query}", "/query.html?&any=1", "true", "query"},

		// /*  /{query}/      */
		{"/{query}/", "/query", "false", ""},
		{"/{query}/", "/query/", "true", "query"},
		{"/{query}/", "/query/?any", "true", "query"},
		{"/{query}/", "/query.html", "false", ""},
		{"/{query}/", "/query.html?any", "false", ""},
		{"/{query}/", "/query.html?&any=1", "false", ""},

		// /*  /query/hoho     */
		{"/query/hoho", "/query", "false", ""},
		{"/query/hoho", "/query/", "false", ""},
		{"/query/hoho", "/query/any", "false", ""},
		{"/query/hoho", "/query/hoho", "true", ""},
		{"/query/hoho", "/query/hoho/?any", "true", ""},
		{"/query/hoho", "/query/hoho?any", "true", ""},
		{"/query/hoho", "/query/hoho.html", "true", ""},
		{"/query/hoho", "/query/hoho.html?any", "true", ""},
		{"/query/hoho", "/query/hoho.html?&any=1", "true", ""},
		//{"/query/hoho", "/query/hohoany", "false", ""}, /*-------------------------------------*/
		{"/query/hoho", "/queryany/hohoany", "false", ""},

		// /*  /query/{type}     */
		{"/query/{type}", "/query", "false", ""},
		{"/query/{type}", "/query/", "true", ""},
		{"/query/{type}", "/query/type", "true", "type"},
		{"/query/{type}", "/query/type/", "true", "type"},
		{"/query/{type}", "/query/type/any", "true", "type"},
		{"/query/{type}", "/query/type/?", "true", "type"},
		{"/query/{type}", "/query/type?", "true", "type"},
		{"/query/{type}", "/query/type.html", "true", "type"},
		{"/query/{type}", "/query/type.html?&any=1", "true", "type"},
		{"/query/{type}", "/query/?any", "true", ""},
		{"/query/{type}", "/query?any", "false", ""},
		{"/query/{type}", "/query.html", "false", ""},
		{"/query/{type}", "/query.html?any", "false", ""},
		{"/query/{type}", "/query.html?&any=1", "false", ""},
		{"/query/{type}", "/queryhoho", "false", ""},
		{"/query/{type}", "/queryhoho/type", "false", ""},

		// /*  /query/{type}/     */
		{"/query/{type}/", "/query", "false", ""},
		{"/query/{type}/", "/query/", "false", ""},
		{"/query/{type}/", "/query/type", "false", ""},
		{"/query/{type}/", "/query/type/", "true", "type"},
		{"/query/{type}/", "/query/type/any", "true", "type"},
		{"/query/{type}/", "/query/type/?", "true", "type"},
		{"/query/{type}/", "/query/type?", "false", ""},
		{"/query/{type}/", "/query/type.html", "false", ""},
		{"/query/{type}/", "/query/type.html?&any=1", "false", ""},
		{"/query/{type}/", "/query/?any", "false", ""},
		{"/query/{type}/", "/query?any", "false", ""},
		{"/query/{type}/", "/query.html", "false", ""},
		{"/query/{type}/", "/query.html?hoho", "false", ""},
		{"/query/{type}/", "/query.html?&hoho=1", "false", ""},
		{"/query/{type}/", "/queryhoho/", "false", ""},
		{"/query/{type}/", "/queryhoho/type/", "false", ""},

		// /*  /query/{type}/{year}/{month}/{day}     */
		{"/query/{type}/{year}/{month}/{day}", "/query", "false", ""},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year", "false", ""},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month", "false", ""},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day/", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day/any", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day/?", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day?", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day.html?hoho", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/query/type/year/month/day.html?&hoho=1", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}", "/queryhoho/type/year/month/day", "false", ""},
		{"/query/{type}/{year}/{month}/{day}", "/queryhoho/type", "false", ""},

		// /*  /query/{type}/{year}/{month}/{day}/     */
		{"/query/{type}/{year}/{month}/{day}/", "/query", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day/", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day/any", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day/?", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day?", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day.html", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day.html?hoho", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/type/year/month/day.html?&hoho=1", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/queryhoho/type/year/month/day", "false", ""},
		{"/query/{type}/{year}/{month}/{day}", "/queryhoho/type", "false", ""},

		// /*  /query/{type}/{year}-{month}-{day}    */
		{"/query/{type}/{year}-{month}-{day}", "/query", "false", ""},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year", "false", ""},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month", "false", ""},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day/", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day/any", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day/?", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day?", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day.html?hoho", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}", "/query/type/year-month-day.html?&hoho=1", "true", "type,year,month,day"},
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
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year-month-day/?", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year-month-day?", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year-month-day.html", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year-month-day.html?hoho", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/", "/query/type/year-month-day.html?&hoho=1", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/", "/query/typehoho/", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/", "/query/typehoho/year", "false", ""},

		// /*  /query/{type}/{year}/{month}/{day}.html    */
		{"/query/{type}/{year}/{month}/{day}.html", "/query", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day/", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day/any", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day/?", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day?", "false", ""},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day.html?hoho", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}.html", "/query/type/year/month/day.html?&hoho=1", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}.html", "/queryhoho/type/year/month/day.html", "false", ""},

		// /*  /query/{type}/{year}/{month}/{day}/detail.html    */
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day/", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day/any", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day/?", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day?", "false", ""},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day/detail.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day/detail.html?hoho", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/query/type/year/month/day/detail.html?&hoho=1", "true", "type,year,month,day"},
		{"/query/{type}/{year}/{month}/{day}/detail.html", "/queryhoho/type/year/month/day/detail.html", "false", ""},

		// /*  /query/{type}/{year}-{month}-{day}.html    */
		{"/query/{type}/{year}-{month}-{day}.html", "/query", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day/", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day/any", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day/?", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day?", "false", ""},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day.html?hoho", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}.html", "/query/type/year-month-day.html?&hoho=1", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}.html", "/queryhoho/type/year-month-day.html", "false", ""},

		// /*  /query/{type}/{year}-{month}-{day}/detail.html    */
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month-day", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month/day-", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month-day/any", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month-day/?", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month-day?", "false", ""},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month-day/detail.html", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month-day/detail.html?hoho", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}/detail.html", "/query/type/year-month-day/detail.html?&hoho=1", "true", "type,year,month,day"},
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
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}", "/query/type/year-month-day-1-orderby/?", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}/", "/query/type/year-month-day-1-orderby", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}/", "/query/type/year-month-day-1-orderby/", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}/", "/query/type/year-month-day-1-orderby.html", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}/", "/query/type/year-month-day-1-orderby/?any", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}.html", "/query/type/year-month-day-1-orderby", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1-{orderby}.html", "/query/type/year-month-day-1-orderby.html?any", "true", "type,year,month,day,orderby"},

		{"/query/{type}/{year}-{month}-{day}-1/{orderby}", "/query/type/year-month-day-1", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}", "/query/type/year-month-day-1/", "true", "type,year,month,day"},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}", "/query/type/year-month-day-1/orderby", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}", "/query/type/year-month-day-1/orderby.html", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}", "/query/type/year-month-day-1/orderby/?", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}/", "/query/type/year-month-day-1/orderby", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}/", "/query/type/year-month-day-1/orderby/", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}/", "/query/type/year-month-day-1/orderby.html", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}/", "/query/type/year-month-day-1/orderby/?any", "true", "type,year,month,day,orderby"},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}.html", "/query/type/year-month-day-1/orderby", "false", ""},
		{"/query/{type}/{year}-{month}-{day}-1/{orderby}.html", "/query/type/year-month-day-1/orderby.html?any", "true", "type,year,month,day,orderby"},

		{"/", "/", "true", ""}}

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

		parameters := re.FindAllStringSubmatch(input)
		matched := parameters != nil
		if (!matched && expect == "true") || (matched && expect == "false") {
			t.Errorf("FindAllStringSubmatch fail: expect:%s; actual=%v; pattern=%s; input=%s \n", expect, parameters != nil, pattern, input)
			continue
		}

		if ((!matched || len(parameters) == 0) && names != "") || ((matched && len(parameters) > 0) && names == "") {
			t.Errorf("FindAllStringSubmatch sub names config fail: expect:%v; actual=%v; pattern=%s; input=%s \n", names == "", matched && len(parameters) > 0, pattern, input)
			continue
		}

		if names == "" {
			continue
		}

		namesSlice := strings.Split(names, ",")
		if len(namesSlice) != len(parameters) {
			t.Errorf("FindAllStringSubmatch sub names length fail: expect:%v; actual=%v; pattern=%s; input=%s \n", len(namesSlice), len(parameters), pattern, input)
			continue
		}

		if len(namesSlice) == 0 {
			continue
		}

		for i, p := range parameters {
			if p[0] != namesSlice[i] {
				t.Errorf("FindAllStringSubmatch sub names check name fail: expect:%s; actual=%s; pattern=%s; input=%s \n", namesSlice[0], p[0], pattern, input)
			} else if p[0] != p[1] {
				t.Errorf("FindAllStringSubmatch sub names check value fail: expect:%s; actual=%s; pattern=%s; input=%s \n", p[0], p[1], pattern, input)
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
		// "/query/{type}/{*values}", TODO
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

		//"/query/{type}/{year:[0-9]+}-{month:[0-9]+}-{day:[0-9]+}", //TODO
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
		"/{query@#}",
		"/{?query}",
		"/query/{type}.html?{month}",
		"/query/{type}.{month}html",
		""}

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

func BenchmarkNoSubName(b *testing.B) {

	pattern := "/query/hoho/{type}/{year}-{month}-{day}/"
	input := "/query/hoho/type/year-month-day/"

	re, err := pathexp.Compile(pattern)
	if err != nil {
		b.Errorf("pattern compile fail: %s; pattern=%s \n", err, pattern)
		return
	}

	for i := 0; i < b.N; i++ {
		_ = re.FindAllStringSubmatch(input)
	}

}
