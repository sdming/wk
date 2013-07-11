// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package regax_test

import (
	"github.com/sdming/wk/regax"
	"strings"
	"testing"
)

func BenchmarkSuffix(b *testing.B) {
	b.StopTimer()
	p := "*.jpg"
	s := "a.jpg"
	re := regax.Compile(p)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(s)
	}
}

func BenchmarkPrefix(b *testing.B) {
	b.StopTimer()
	p := "text/*"
	s := "text/html"
	re := regax.Compile(p)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(s)
	}
}

func BenchmarkHT(b *testing.B) {
	b.StopTimer()
	p := "*/*"
	s := "text/html"
	re := regax.Compile(p)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(s)
	}
}

func BenchmarkRegexp(b *testing.B) {
	b.StopTimer()
	p := "a*b*c"
	s := "axbxc"
	re := regax.Compile(p)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		re.Match(s)
	}
}

func TestMuti(t *testing.T) {
	p := "/admin*;*.jpg;text/html;application/*;*abc*;a*b*c;;"
	var tc []string = []string{
		"/adminx",
		"a.jpg",
		"text/html",
		"application/javascript",
		"xabcx",
		"axbxc",
	}

	mre := regax.NewMutiRegAx(p)
	for _, each := range tc {
		matched := mre.Match(each)

		if !matched {
			t.Error(matched, each)
		}
	}
}

func TestMatch(t *testing.T) {
	var tc [][2]string = [][2]string{
		{"*", "a;*;"},
		{"a.jpg", "a.jpg;"},
		{"*.jpg", "a.jpg;.jpg;"},
		{"/a*", "/axx;/a;"},
		{"a*b", "ab;axb;"},
		{"*x*", "axb;ax;xb;x;"},
		{"a*b*c", "axbxc"},
	}

	for _, each := range tc {
		p := each[0]
		a := each[1]

		re := regax.Compile(p)
		for _, s := range strings.Split(a, ";") {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}

			matched := re.Match(s)

			if !matched {
				t.Error(matched, p, s)
			}
		}
	}

}

func TestFail(t *testing.T) {
	var tc [][2]string = [][2]string{
		{"a.jpg", "x.jpg;xa.jpg;a.jpgx;"},
		{"*.jpg", "jpg;xjpg;a.jpgx"},
		{"/a*", "/x;/;"},
		{"a*b", "xab;abx;"},
		{"*x*", "ab;"},
		{"a*b*c", "bxaxc;axbx;ab;bc;"},
	}

	for _, each := range tc {
		p := each[0]
		a := each[1]

		re := regax.Compile(p)
		for _, s := range strings.Split(a, ";") {
			if s == "" {
				continue
			}

			matched := re.Match(s)
			if matched {
				t.Error(matched, p, s)
			}
		}
	}

}
