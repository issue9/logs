// SPDX-License-Identifier: MIT

// Package initfunc 实现了 github.com/issue9/writers 下的 WriterInitializer 接口
package initfunc

import "fmt"

func argNotFoundErr(wname, argName string) error {
	return fmt.Errorf("[%s]配置文件中未指定参数:[%s]", wname, argName)
}
