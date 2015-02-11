// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writers

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/issue9/assert"
)

var _ WriteFlushAdder = &Buffer{}

func TestBuffer(t *testing.T) {
	a := assert.New(t)

	// 正确的NewBuffer()
	buf := NewBuffer(0)
	a.NotNil(buf)

	// 不能添加空值
	a.Error(buf.Add(nil))

	// 不缓存，直接输出，但又没指定输出方向，相当于直接扔掉！
	size, err := buf.Write([]byte("abc"))
	a.NotError(err).Equal(0, size)

	b1 := bytes.NewBufferString("")
	b2 := bytes.NewBufferString("")
	a.NotNil(b1).NotNil(b2)

	// 设置了缓存，不会直接输出内容
	buf.SetSize(10)
	size, err = buf.Write([]byte("0"))
	a.NotError(err).True(size > 0)

	// 内容应该还在缓存中，b1中还是空
	a.NotError(buf.Add(b1))
	a.Equal(b1.Len(), 0)

	// 补满10次。向b1输出内容。
	for i := 1; i < 10; i++ {
		size, err = buf.Write([]byte(strconv.Itoa(i)))
		a.NotError(err).True(size > 0)
	}
	a.Equal(b1.String(), "0123456789")

	// 添加B2
	buf.Add(b2)
	size, err = buf.Write([]byte("9"))
	a.NotError(err).True(size > 0)

	a.Equal(b2.Len(), 0).Equal(b1.Len(), 10)
	for i := 8; i >= 0; i-- {
		size, err = buf.Write([]byte(strconv.Itoa(i)))
		a.NotError(err).True(size > 0)
	}
	a.Equal(b1.Len(), 20).Equal(b1.String(), "01234567899876543210")
	a.Equal(b2.Len(), 10).Equal(b2.String(), "9876543210")
}
