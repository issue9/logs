// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"io"
	"testing"

	"github.com/issue9/assert"
)

func init1(a1 map[string]string) (io.Writer, error) {
	return nil, nil
}

func init2(a1 map[string]string) (io.Writer, error) {
	return nil, nil
}

func TestInit(t *testing.T) {
	a := assert.New(t)

	// 清空，包的init函数有可能会初始化一些数据。
	clearInitializer()

	a.True(Register("init1", init1)).
		True(IsRegisted("init1")).
		Equal(Registed(), []string{"init1"})

	a.True(Register("init2", init2)).
		True(IsRegisted("init2")).
		True(IsRegisted("init1")).
		Equal(Registed(), []string{"init1", "init2"})

	a.False(IsRegisted("init3"))

	a.False(Register("init1", init2)) // 重复注册
	a.True(IsRegisted("init1"))

	clearInitializer()
	a.Equal(0, len(regNames)).Equal(0, len(regInitializer))
}
