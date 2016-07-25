// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writers

import (
	"os"

	"github.com/issue9/term/colors"
)

// 带色彩输出的控制台。
type Console struct {
	out *os.File
	c   colors.Colorize
}

// 新建Console实例
//
// out为输出方向，可以是colors.Stderr和colors.Stdout两个值。
// foreground,background 为输出文字的前景色和背景色。
func NewConsole(out *os.File, foreground, background colors.Color) *Console {
	return &Console{
		out: out,
		c:   colors.New(foreground, background),
	}
}

// 更改输出颜色
func (c *Console) SetColor(foreground, background colors.Color) {
	c.c.Foreground = foreground
	c.c.Background = background
}

// io.Writer
func (c *Console) Write(b []byte) (size int, err error) {
	return c.c.Fprint(c.out, string(b))
}
