// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/logs/writers"
)

func TestToByte(t *testing.T) {
	a := assert.New(t)

	eq := func(str string, val int64) {
		size, err := toByte(str)
		a.NotError(err).Equal(size, val)
	}

	e := func(str string) {
		size, err := toByte(str)
		a.Error(err).Equal(size, -1)
	}

	eq("1m", 1024*1024)
	eq("100G", 100*1024*1024*1024)
	eq("10.2k", 10*1024)
	eq("10.9K", 10*1024)

	e("")
	e("M")
	e("-1M")
	e("-1.0G")
	e("1P")
	e("10MB")
}

func TestRotate(t *testing.T) {
	a := assert.New(t)
	args := map[string]string{}

	w, err := Rotate(args)
	a.Error(err).Nil(w)

	// 缺少 size
	args["dir"] = "./testdata"
	w, err = Rotate(args)
	a.Error(err).Nil(w)

	// 错误的 size 参数
	args["size"] = "12P"
	w, err = Rotate(args)
	a.Error(err).Nil(w)

	// 正常
	args["size"] = "12"
	w, err = Rotate(args)
	a.NotError(err).NotNil(w)

	_, ok := w.(*writers.Rotate)
	a.True(ok)
}
