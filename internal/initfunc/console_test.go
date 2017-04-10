// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/logs/writers"
)

func TestConsole(t *testing.T) {
	a := assert.New(t)
	args := map[string]string{}

	// 可以接受空参数，Console 的 args 都有默认值
	w, err := Console(args)
	a.NotError(err).NotNil(w)

	// 无效的 output
	args["output"] = "file"
	w, err = Console(args)
	a.Error(err).Nil(w)
	args["output"] = "stderr"

	// 无效的 foreground
	args["foreground"] = "red1"
	w, err = Console(args)
	a.Error(err).Nil(w)
	args["foreground"] = "red"

	// 无效的 background
	args["background"] = "red1"
	w, err = Console(args)
	a.Error(err).Nil(w)

	args["background"] = "blue"
	w, err = Console(args)
	a.NotError(err).NotNil(w)

	_, ok := w.(*writers.Console)
	a.True(ok)
}
