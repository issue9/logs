// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/issue9/logs/v2/config"
	"github.com/issue9/logs/v2/writers"
)

func TestConsole(t *testing.T) {
	a := assert.New(t)
	cfg := &config.Config{Attrs: map[string]string{}}

	// 可以接受空参数，Console 的 args 都有默认值
	w, err := Console(cfg)
	a.NotError(err).NotNil(w)

	// 无效的 output
	cfg.Attrs["output"] = "file"
	w, err = Console(cfg)
	a.Error(err).Nil(w)
	cfg.Attrs["output"] = "stderr"

	// 无效的 foreground
	cfg.Attrs["foreground"] = "red1"
	w, err = Console(cfg)
	a.Error(err).Nil(w)
	cfg.Attrs["foreground"] = "red"

	// 无效的 background
	cfg.Attrs["background"] = "red1"
	w, err = Console(cfg)
	a.Error(err).Nil(w)

	cfg.Attrs["background"] = "blue"
	w, err = Console(cfg)
	a.NotError(err).NotNil(w)

	_, ok := w.(*writers.Console)
	a.True(ok)
}
