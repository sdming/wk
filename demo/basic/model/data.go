// Copyright 2012 by sdm. All rights reserved.
// license that can be found in the LICENSE file.

/*
basic demo

*/
package model

import (
	"fmt"
	"strconv"
)

type Data struct {
	Str   string
	Uint  uint64
	Int   int
	Float float32
	Byte  byte
}

func newData(i int) *Data {
	return &Data{
		Str:   "string:" + strconv.Itoa(i),
		Uint:  uint64(i * 100),
		Int:   i * 10,
		Float: 0.1 + float32(i),
		Byte:  byte(i),
	}
}

func DataTop(count int) []*Data {
	if count <= 0 || count > 100 {
		panic("count is invalid")
	}

	data := make([]*Data, count)
	for i := 0; i < count; i++ {
		data[i] = newData(i)
	}

	return data
}

func DataByInt(i int) *Data {
	if i < 0 {
		i = 0
	}
	return newData(i)
}

func DataByIntRange(start, end int) []*Data {
	data := make([]*Data, 2)
	data[0] = newData(start)
	data[1] = newData(end)
	return data
}

func DataSet(s string, u uint64, i int, f float32, b byte) *Data {
	return &Data{
		Str:   s,
		Uint:  u,
		Int:   i,
		Float: f,
		Byte:  b,
	}
}

func DataDelete(i int) string {
	return fmt.Sprintf("delete %d", i)
}

func DataPost(data Data) string {
	return data.String()
}

func DataPostPtr(data *Data) string {
	return data.String()
}

func DataPostWithError(data Data) (int, error) {
	fmt.Println(data)
	return 0, nil
}

func (d *Data) String() string {
	if d == nil {
		return "<nil>"
	}
	return fmt.Sprintf("string:%s;uint:%d;int:%d;float:%f;byte:%c",
		d.Str, d.Uint, d.Int, d.Float, d.Byte)
}
