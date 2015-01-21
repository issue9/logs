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

	b1 := bytes.NewBufferString("")
	b2 := bytes.NewBufferString("")
	a.NotNil(b1).NotNil(b2)

	buf := NewBuffer(nil, 10)

	size, err := buf.Write([]byte("0"))
	a.NotError(err).True(size > 0)

	err = buf.Add(b1)
	a.NotError(err)

	// 仅写入一次，应该没有向b1输出内容。
	a.Equal(b1.Len(), 0)

	// 补满10次。向b1输出内容。
	for i := 1; i < 10; i++ {
		size, err = buf.Write([]byte(strconv.Itoa(i)))
		a.NotError(err).True(size > 0)
	}
	a.Equal(b1.Len(), 10).Equal(b1.String(), "0123456789")

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
