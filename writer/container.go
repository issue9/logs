// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writer

import (
	"io"
)

// 对WriterContainer的默认实现。
type Container struct {
	ws []io.Writer
}

var _ WriteFlushAdder = &Container{}

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

func (c *Container) Add(w io.Writer) error {
	c.ws = append(c.ws, w)
	return nil
}

func (c *Container) Flush() (size int, err error) {
	for _, w := range c.ws {
		if b, ok := w.(Flusher); ok {
			size, err = b.Flush()
		}
	}
	return size, err
}
