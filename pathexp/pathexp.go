// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package pathexp

import (
	//"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Matcher interface {
	Match(expr string) (matched bool, submatch [][2]string)
}

type RegexpMatch struct {
	pattern        string
	re             *regexp.Regexp
	subNames       []string
	subCount       int
	subNameedCount int
}

func RegexpCompile(pattern string) (*RegexpMatch, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	subNames := re.SubexpNames()
	count := 0
	for _, s := range subNames {
		if s != "" {
			count++
		}
	}
	return &RegexpMatch{
		pattern:        pattern,
		re:             re,
		subNames:       subNames,
		subCount:       len(subNames),
		subNameedCount: count,
	}, nil
}

func (re *RegexpMatch) Match(expr string) (matched bool, submatch [][2]string) {
	var data [][]string = re.re.FindAllStringSubmatch(expr, -1)

	if data == nil || len(data) != 1 {
		return false, nil
	}
	if len(data[0]) != re.subCount {
		return false, nil
	}
	if len(data[0][0]) != len(expr) {
		return false, nil
	}
	matched = true
	//submatch = make(map[string]string)
	submatch = make([][2]string, 0, re.subNameedCount)

	l := len(data[0])
	d := data[0]
	for i := 0; i < l; i++ {
		if re.subNames[i] != "" {
			//submatch[re.subNames[i]] = d[i]
			submatch = append(submatch, [2]string{re.subNames[i], d[i]})
		}
	}
	return
}

var debug bool = false

func patternError(text string, index int) error {
	return errors.New(text + fmt.Sprintf(" index %d", index))
}

// Pathex is a parttern to match url path
type Pathex struct {
	Pattern  string
	SubNames []SubName
	prefix   []byte
	ps       string
	p        []byte
	pi       int
}

// string
func (re Pathex) String() string {
	return re.Pattern
}

func (re *Pathex) Match(chars string) (matched bool, data [][2]string) {
	//chars := []byte(expr)

	charLen, patternLen := len(chars), len(re.Pattern)
	if charLen < re.pi {
		return false, nil
	}

	if !(charLen >= re.pi && chars[0:re.pi] == re.ps) {
		return false, nil
	}

	var charIndex, patternIndex, nameIndex int
	var c, p byte
	matched = false
	data = make([][2]string, 0, len(re.SubNames))

	charIndex = re.pi
	patternIndex = re.pi

	for {
		if patternIndex >= patternLen {
			break
		}

		p = re.p[patternIndex]
		if p == '{' {
			value := make([]byte, 0, 32)
			subName := re.SubNames[nameIndex]

			for {
				if charIndex >= charLen {
					break
				}

				c = chars[charIndex]
				if c == '/' || c == '.' || c == '?' || (subName.end != byte(0) && c == subName.end) {
					break
				} else {
					value = append(value, c)
				}
				charIndex++
			}

			if len(value) != 0 {
				data = append(data, [2]string{subName.Name, string(value)})
			}
			patternIndex = patternIndex + subName.length + 2
			nameIndex++
		} else if p == '?' {
			if charIndex < charLen {
				return false, nil
			}
			return true, data
		} else if p == '#' {
			matched = true
			patternIndex++
			//charIndex++
		} else {
			if charIndex >= charLen {
				if !matched {
					return false, nil
				} else {
					break
				}
			}

			c = chars[charIndex]
			if p != c {
				if !matched {
					return false, nil
				} else {
					break
				}
			}
			patternIndex++
			charIndex++

			if patternIndex == patternLen && p != '/' && charIndex < charLen {
				c = chars[charIndex]
				if c != '/' && c != '.' && c != '?' {
					return false, nil
				}
			}
		}
	}

	return true, data
}

// MatchString returns whether the Pathex matches the string expr
func (re *Pathex) MatchString(expr string) bool {
	matched, _ := re.Match(expr)
	return matched
}

// func (re *Pathex) Match(expr string) (matched bool, submatch map[string]string) {
// 	return re.execute(expr)
// }

// // FindAllStringSubmatch returns a slice of all successive matches of the expression, nil indicates no match.
// func (re *Pathex) FindAllStringSubmatch(s string) [][2]string {
// 	matched, data := re.execute(s)
// 	if !matched {
// 		return nil
// 	}
// 	return data
// }

// subName
type SubName struct {
	Name   string
	end    byte
	length int
}

// string
func (self SubName) String() string {
	return self.Name
}

// Compile parses a expression 
func Compile(s string) (exp *Pathex, err error) {
	s = strings.TrimSpace(s)
	if s == "" {
		err = patternError("pattern can not be empty", 0)
		return
	}

	if s[0] != '/' {
		err = patternError("pattern need start with /", 0)
		return
	}

	text := []byte(s)
	subNames := make([]SubName, 0)

	var name []byte
	nameIndex := -1
	dotIndex := -1
	numberIndex := -1
	pi := 0

	for i, c := range text {
		if (nameIndex > -1) && (c == '{' || c == '/' || c == '.' || c == '#' || c == '?') {
			err = patternError("brackets fail", i)
			return
		}
		if (dotIndex > -1) && (c == '{' || c == '/' || c == '.' || c == '}') {
			err = patternError("pattern '.' match fail", i)
			return
		}

		if c == '{' || c == '?' || c == '#' {
			if pi == 0 {
				pi = i
			}
		}

		switch c {
		case '{':
			if text[i-1] == '}' {
				err = patternError("need delimiter between bracket", i)
				return
			}

			nameIndex = 0
			name = make([]byte, 0)
		case '}':
			if nameIndex < 0 {
				err = patternError("can not find a match for {", i)
				return
			}
			if len(name) == 0 {
				err = patternError("subName can not be empty", i)
				return
			}
			nameIndex = -1
			end := byte(0)
			if i < len(text)-1 {
				end = text[i+1]
			}
			subNames = append(subNames, SubName{
				Name:   string(name),
				length: len(name),
				end:    end})
		case '/':
			//pathIndex++
		case '?':
			if i != len(text)-1 || numberIndex > 0 {
				err = patternError("? is invalid", i)
				return
			}
		case '.':
			if i < 1 || dotIndex > 0 {
				err = patternError("'.' is invalid", i)
				return
			}
			dotIndex = i
		case '#':
			if i < 1 || numberIndex > 0 {
				err = patternError("'#' is invalid", i)
				return
			}
			numberIndex = i
		default:
			if nameIndex > -1 {
				if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || (c == '_') {
					name = append(name, c)
					nameIndex++
				} else {
					err = patternError("subName should be a-z, A-Z, 0-9 or _ ", i)
					return
				}
			}
		}

	}

	if nameIndex > -1 {
		err = patternError("brackets match fail", len(text))
		return
	}

	if pi == 0 {
		pi = len(text)
	}

	return &Pathex{
		Pattern:  s,
		SubNames: subNames,
		p:        text,
		prefix:   text[:pi],
		ps:       string(text[:pi]),
		pi:       pi,
	}, nil
}

// MatchString returns whether the pattern matches the string s
func MatchString(pattern string, s string) (matched bool, error error) {
	re, err := Compile(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(s), nil
}
