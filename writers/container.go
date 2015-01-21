// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writers

import (
	"io"
)

// 对WriterContainer的默认实现。
type Container struct {
	ws []io.Writer
}

// 构造Container实例
func NewContainer(writers ...io.Writer) *Container {
	return &Container{ws: writers}
}

// 当某一项出错时，将直接返回其信息，后续的都将中断。
func (c *Container) Write(bs []byte) (size int, err error) {
	for _, w := range c.ws {
		if size, err = w.Write(bs); err != nil {
			return
		}
	}

	return
}

// 添加一个io.Writer实例
func (c *Container) Add(w io.Writer) error {
	c.ws = append(c.ws, w)
	return nil
}

// 调用所有子项的Flush函数。
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
