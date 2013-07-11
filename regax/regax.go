// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package regax

import (
	"regexp"
	"strings"
)

// MutiRegAx 
type MutiRegAx struct {
	Pattern string
	res     []*RegAx
}

func NewMutiRegAx(pattern string) *MutiRegAx {
	pattern = strings.TrimSpace(pattern)
	ps := strings.Split(pattern, ";")
	res := make([]*RegAx, 0)

	for _, p := range ps {
		res = append(res, Compile(strings.TrimSpace(p)))
	}
	return &MutiRegAx{
		Pattern: pattern,
		res:     res,
	}
}

func (mre *MutiRegAx) Match(text string) bool {
	if mre.res == nil {
		return false
	}

	count := len(mre.res)
	for i := 0; i < count; i++ {
		re := mre.res[i]
		if re.Match(text) {
			return true
		}
	}
	return false
}

// RegAx is a very simple regexp match, '*' mean match anything, other mean match themself
type RegAx struct {
	Pattern string
	m       matcher
}

func (re RegAx) String() string {
	return re.Pattern
}

func (re *RegAx) Match(text string) bool {
	return re.m.match(text)
}

type matcher interface {
	match(text string) bool
}

type anyMatcher struct {
}

func (m *anyMatcher) match(text string) bool {
	return true
}

type equalMatcher struct {
	text string
}

func (m *equalMatcher) match(text string) bool {
	return text == m.text
}

type prefixMatcher struct {
	prefix string
}

func (m *prefixMatcher) match(text string) bool {
	return strings.HasPrefix(text, m.prefix)
}

type suffixMatcher struct {
	suffix string
}

func (m *suffixMatcher) match(text string) bool {
	return strings.HasSuffix(text, m.suffix)
}

type htMatcher struct {
	suffix string
	prefix string
	length int
}

func (m *htMatcher) match(text string) bool {
	return len(text) >= m.length && strings.HasPrefix(text, m.prefix) && strings.HasSuffix(text, m.suffix)
}

type containsMatcher struct {
	text string
}

func (m *containsMatcher) match(text string) bool {
	return strings.Contains(text, m.text)
}

type regexpMatcher struct {
	pattern string
	re      *regexp.Regexp
}

func (m *regexpMatcher) match(text string) bool {
	return m.re.MatchString(text)
}

func Compile(s string) *RegAx {
	pattern := strings.TrimSpace(s)

	if pattern == "" || pattern == "*" || strings.Replace(pattern, "*", "", -1) == "" {
		return &RegAx{
			Pattern: pattern,
			m:       &anyMatcher{},
		}
	}

	count := strings.Count(s, "*")
	if count == 0 {
		return &RegAx{
			Pattern: pattern,
			m:       &equalMatcher{text: pattern},
		}
	}

	if count == 1 {
		if strings.HasSuffix(s, "*") {
			return &RegAx{
				Pattern: pattern,
				m:       &prefixMatcher{prefix: strings.Replace(pattern, "*", "", -1)},
			}
		} else if strings.HasPrefix(s, "*") {
			return &RegAx{
				Pattern: pattern,
				m:       &suffixMatcher{suffix: strings.Replace(pattern, "*", "", -1)},
			}
		} else {
			index := strings.Index(pattern, "*")
			return &RegAx{
				Pattern: pattern,
				m: &htMatcher{
					prefix: pattern[0:index],
					suffix: pattern[index+1:],
					length: len(pattern) - 1,
				},
			}
		}
	}

	if count == 2 && strings.HasSuffix(s, "*") && strings.HasPrefix(s, "*") {
		return &RegAx{
			Pattern: pattern,
			m: &containsMatcher{
				text: strings.Replace(pattern, "*", "", -1),
			},
		}
	}

	return &RegAx{
		Pattern: pattern,
		m: &regexpMatcher{
			pattern: strings.Replace(pattern, "*", `[\S]*`, -1),
			re:      regexp.MustCompile(strings.Replace(pattern, "*", `[\S]*`, -1)),
		},
	}
}
