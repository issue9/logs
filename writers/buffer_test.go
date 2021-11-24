// SPDX-License-Identifier: MIT

package writers

import (
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/issue9/assert/v2"
)

var (
	_ Adder          = &Buffer{}
	_ Flusher        = &Buffer{}
	_ io.WriteCloser = &Buffer{}
)

func TestBuffer(t *testing.T) {
	a := assert.New(t, false)

	// 正确的 NewBuffer()
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

func TestBuffer_Close(t *testing.T) {
	a := assert.New(t, false)
	buffer := NewBuffer(10)

	c1 := &testContainer{}
	c2 := &testContainer{}
	a.NotError(buffer.Add(c1))
	a.NotError(buffer.Add(c2))

	a.False(c1.flushed).False(c1.closed).Equal(0, c1.Len())
	a.False(c2.flushed).False(c2.closed).Equal(0, c2.Len())

	_, err := buffer.Write([]byte("123"))
	a.NotError(err)
	a.Equal(0, c1.Len())

	a.NotError(buffer.Close())
	a.True(c1.flushed).True(c1.closed).Equal(3, c1.Len())
	a.True(c2.flushed).True(c2.closed).Equal(3, c2.Len())
}
