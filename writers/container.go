// SPDX-License-Identifier: MIT

package writers

import (
	"errors"
	"io"
)

// Container 为 io.Writer 的容器
type Container struct {
	ws []io.Writer
}

// NewContainer 构造 Container 实例
func NewContainer() *Container {
	return &Container{ws: make([]io.Writer, 0, 1)}
}

// 当某一项出错时，将直接返回其信息，后续的都将中断。
// 若容器为空时，则相当于不作任何动作。
func (c *Container) Write(bs []byte) (size int, err error) {
	for _, w := range c.ws {
		if size, err = w.Write(bs); err != nil {
			return
		}
	}

	return
}

// Add 添加一个 io.Writer 实例
func (c *Container) Add(w io.Writer) error {
	if w == nil {
		return errors.New("参数 w 不能为一个空值")
	}

	c.ws = append(c.ws, w)
	return nil
}

// Flush 调用所有子项的 Flush 函数
func (c *Container) Flush() error {
	for _, w := range c.ws {
		if b, ok := w.(Flusher); ok {
			if err := b.Flush(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Container) Close() error {
	if err := c.Flush(); err != nil {
		return err
	}

	for _, w := range c.ws {
		if b, ok := w.(io.Closer); ok {
			if err := b.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

// Len 包含的元素
func (c *Container) Len() int { return len(c.ws) }

// Clear 清除所有的 writer
func (c *Container) Clear() { c.ws = c.ws[:0] }
