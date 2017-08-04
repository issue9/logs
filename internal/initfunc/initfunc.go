// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package initfunc 实现了 github.com/issue9/writers 下的 WriterInitializer 接口。
package initfunc

import (
	"fmt"
)

func argNotFoundErr(wname, argName string) error {
	return fmt.Errorf("[%s]配置文件中未指定参数:[%s]", wname, argName)
}
