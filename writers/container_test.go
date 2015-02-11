// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writers

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

var _ WriteFlushAdder = &Container{}

func TestContainer(t *testing.T) {
	a := assert.New(t)
	b1 := bytes.NewBufferString("")
	b2 := bytes.NewBufferString("")
	a.NotNil(b2).NotNil(b1)

	c := NewContainer()
	a.NotError(c.Add(b1))
	size, err := c.Write([]byte("hello"))
	a.NotError(err).
		True(size > 0).
		Equal(1, c.Len())

	// 传递错误的nil值。
	a.Error(c.Add(nil))

	// 只向c添加了b1，此时b1有内容，b2没内容
	a.Equal("hello", b1.String())
	a.NotEqual("hello", b2.String())

	c.Add(b2)
	size, err = c.Write([]byte(" world"))
	a.NotError(err).
		True(size > 0).
		Equal(2, c.Len())

	// b2后添加，此时b1有全部的内容，而b2只有后半部分。
	a.Equal("hello world", b1.String())
	a.Equal(" world", b2.String())

	// 清除
	c.Clear()
	c.Write([]byte("hello world"))
	a.Equal("hello world", b1.String())
	a.Equal(" world", b2.String())
	a.Equal(0, c.Len())

	// 只添加b2，b1应该保持不变
	c.Add(b2)
	c.Write([]byte("hello"))
	a.Equal("hello world", b1.String())
	a.Equal(" worldhello", b2.String())
	a.Equal(1, c.Len())
}
