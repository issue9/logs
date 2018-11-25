// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/issue9/logs/v2/config"
	"github.com/issue9/logs/v2/writers/rotate"
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
	cfg := &config.Config{Attrs: map[string]string{}}

	w, err := Rotate(cfg)
	a.Error(err).Nil(w)

	// 缺少 size
	cfg.Attrs["dir"] = "./testdata"
	w, err = Rotate(cfg)
	a.Error(err).Nil(w)

	// 错误的 size 参数
	cfg.Attrs["size"] = "12P"
	w, err = Rotate(cfg)
	a.Error(err).Nil(w)

	// 正常
	cfg.Attrs["size"] = "12"
	cfg.Attrs["filename"] = "%Y%i.log"
	w, err = Rotate(cfg)
	a.NotError(err).NotNil(w)

	_, ok := w.(*rotate.Rotate)
	a.True(ok)
}
