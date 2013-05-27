// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

package pathexp

import (
	"errors"
)

var debug bool = false

func patternError(text string) error {
	return errors.New(text)
}

// Pathex is a parttern to match url path
type Pathex struct {
	Pattern  string
	SubNames []SubName
}

// string
func (re Pathex) String() string {
	return re.Pattern
}

func (re *Pathex) execute(chars string) (matched bool, data [][2]string) {
	if chars == "" {
		return false, nil
	}

	charLen, patternLen := len(chars), len(re.Pattern)
	var values [][2]string = make([][2]string, 0, len(re.SubNames))

	if charLen < 1 {
		return false, nil
	}

	var charIndex, patternIndex, nameIndex int

	// if debug {
	// 	fmt.Println("execute:", chars, charLen, patternLen)
	// }

	for {
		p := re.Pattern[patternIndex]
		var c byte

		if p == '{' {
			value := make([]byte, 0, 32)
			subName := re.SubNames[nameIndex]

			// if debug {
			// 	fmt.Println("subName", subName.Name, subName.end)
			// }

			for {
				if charIndex >= charLen {
					break
				}

				c = chars[charIndex]

				// if debug {
				// 	fmt.Printf("c: %d %d %c \n", charIndex, c, c)
				// }

				if c == '/' || c == '.' || c == '?' || (subName.end != byte(0) && c == subName.end) {
					// if debug {
					// 	fmt.Printf("sub name end: %d %c \n", charIndex, c)
					// }
					break
				} else {
					value = append(value, c)
				}

				charIndex++
			}

			if len(value) != 0 {
				values = append(values, [2]string{subName.Name, string(value)})
			}
			patternIndex = patternIndex + subName.length + 2
			nameIndex++
		} else if p == '?' {
			// if debug {
			// 	fmt.Println("?:", charIndex, charLen)
			// }
			if charIndex < charLen {
				return false, nil
			}
			return true, values
		} else {
			if charIndex >= charLen {
				return false, nil
			}

			c = chars[charIndex]

			if p != c {
				return false, nil
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

		if patternIndex >= patternLen {
			break
		}

	}

	return true, values
}

// MatchString returns whether the Pathex matches the string s
func (re *Pathex) MatchString(s string) bool {
	matched, _ := re.execute(s)
	return matched
}

// FindAllStringSubmatch returns a slice of all successive matches of the expression, nil indicates no match.
func (re *Pathex) FindAllStringSubmatch(s string) [][2]string {

	matched, data := re.execute(s)
	if !matched {
		return nil
	}
	return data
}

// subName
type SubName struct {
	pathIndex int
	Name      string
	end       byte
	length    int
	// pattern string not implemented
}

// string
func (self SubName) String() string {
	return self.Name
}

// Compile parses a expression and returns
func Compile(s string) (exp *Pathex, err error) {

	if s == "" {
		err = patternError("pattern can not be empty")
		return
	}

	if s[0] != '/' {
		err = patternError("pattern need start with /")
		return
	}

	text := []byte(s)
	subNames := make([]SubName, 0)

	var name []byte
	nameIndex := -1
	pathIndex := 0
	dotIndex := -1

	for i, c := range text {
		if (nameIndex > -1) && (c == '{' || c == '/' || c == '.') {
			if nameIndex > -1 {
				err = patternError("brackets match fail ")
				return
			}
		}
		if (dotIndex > -1) && (c == '{' || c == '/' || c == '.' || c == '}') {
			if nameIndex > -1 {
				err = patternError("pattern '.' match fail ")
				return
			}
		}

		switch c {
		case '{':
			if text[i-1] == '}' {
				err = patternError("need chars split betwwen bracket ")
				return
			}

			nameIndex = 0
			name = make([]byte, 0)
		case '}':
			if nameIndex < 0 {
				err = patternError("can not find a match for { ")
				return
			}
			if len(name) == 0 {
				err = patternError("subName can not be empty ")
				return
			}
			nameIndex = -1
			end := byte(0)
			if i < len(text)-1 {
				end = text[i+1]
			}
			subNames = append(subNames, SubName{
				pathIndex: pathIndex,
				Name:      string(name),
				length:    len(name),
				end:       end})
		case '/':
			pathIndex++
		case '?':
			if i != len(text)-1 {
				err = patternError("? should at end of pattern? ")
				return
			}
		case '.':
			if i < 1 {
				err = patternError("'.' is invalid ")
				return
			}
			dotIndex = i
		default:
			if nameIndex > -1 {
				if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || (c == '_') {
					name = append(name, c)
					nameIndex++
				} else {
					err = patternError("subName should be a-z, A-Z, 0-9 and _  ")
					return
				}
			}
		}

	}

	if nameIndex > -1 {
		err = patternError("brackets match fail ")
		return
	}

	return &Pathex{s, subNames}, nil
}

// MatchString returns whether the pattern matches the string s
func MatchString(pattern string, s string) (matched bool, error error) {
	re, err := Compile(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(s), nil
}
