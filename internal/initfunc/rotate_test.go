// SPDX-License-Identifier: MIT

package initfunc

import (
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/issue9/logs/v3/config"
	"github.com/issue9/logs/v3/writers/rotate"
)

func TestToByte(t *testing.T) {
	a := assert.New(t, false)

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
	a := assert.New(t, false)
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
