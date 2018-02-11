// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writers

import (
	"errors"
	"io"
)

// Container 为 io.Writer 的容器。
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
		return errors.New("参数w不能为一个空值")
	}

	c.ws = append(c.ws, w)
	return nil
}

// Flush 调用所有子项的 Flush 函数。
func (c *Container) Flush() (size int, err error) {
	for _, w := range c.ws {
		b, ok := w.(Flusher)
		if !ok {
			continue
		}

		if size, err = b.Flush(); err != nil {
			return size, err
		}
	}
	return size, err
}

// Len 包含的元素
func (c *Container) Len() int {
	return len(c.ws)
}

// Clear 清除所有的 writer
func (c *Container) Clear() {
	c.ws = c.ws[:0]
}
