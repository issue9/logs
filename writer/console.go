// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writer

import (
	"fmt"
	"io"

	"github.com/issue9/term"
)

// 带色彩输出的控制台。不支持windows系统。
type Console struct {
	color string
	w     io.Writer
}

// 新建Console实例
//
// color ansi的颜色控制符，有关颜色定义字符串在term包中已经定义。
// w 控制台实例，只能是os.Stderr,osStdout，其它将不会显示颜色。
func NewConsole(w io.Writer, color string) *Console {
	return &Console{
		color: color,
		w:     w,
	}
}

// 更改输出颜色
func (c *Console) SetColor(color string) {
	c.color = color
}

// io.Writer
func (c *Console) Write(b []byte) (size int, err error) {
	// 写入颜色
	fmt.Fprintf(c.w, "%v", c.color)

	size, err = c.w.Write(b)

	// 恢复默认值
	fmt.Fprintf(c.w, "%v", term.Reset)

	return
}
