// SPDX-License-Identifier: MIT

package writers

import (
	"bytes"
	"io"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ Adder          = &Container{}
	_ Flusher        = &Container{}
	_ io.WriteCloser = &Container{}
)

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

	// 只向 c 添加了 b1，此时 b1 有内容，b2 没内容
	a.Equal("hello", b1.String())
	a.NotEqual("hello", b2.String())

	a.NotError(c.Add(b2))
	size, err = c.Write([]byte(" world"))
	a.NotError(err).
		True(size > 0).
		Equal(2, c.Len())

	// b2 后添加，此时 b1 有全部的内容，而 b2 只有后半部分。
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

func TestContainer_Close(t *testing.T) {
	a := assert.New(t)
	c := NewContainer()

	c1 := &testContainer{}
	c2 := &testContainer{}
	a.NotError(c.Add(c1))
	a.NotError(c.Add(c2))

	a.False(c1.flushed).False(c1.closed)
	a.False(c2.flushed).False(c2.closed)

	a.NotError(c.Close())
	a.True(c1.flushed).True(c1.closed)
	a.True(c2.flushed).True(c2.closed)
}
