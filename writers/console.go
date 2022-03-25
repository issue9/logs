// SPDX-License-Identifier: MIT

package writers

import (
	"os"

	"github.com/issue9/term/v3/colors"
)

// Console 带色彩输出的控制台
type Console struct {
	out *os.File
	c   colors.Colorize
}

// NewConsole 新建 Console 实例
//
// out 为输出方向，可以是 colors.Stderr 和 colors.Stdout 两个值。
// foreground,background 为输出文字的前景色和背景色。
func NewConsole(out *os.File, foreground, background colors.Color) *Console {
	return &Console{
		out: out,
		c:   colors.New(colors.Normal, foreground, background),
	}
}

// SetColor 更改输出颜色
func (c *Console) SetColor(foreground, background colors.Color) {
	c.c.Foreground = foreground
	c.c.Background = background
}

func (c *Console) Write(b []byte) (size int, err error) {
	return c.c.Fprint(c.out, string(b))
}
