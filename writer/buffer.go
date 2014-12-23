// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writer

import (
	"errors"
	"io"
)

// Buffer 实现对输出内容的缓存，只有输出数量达到指定的值
// 才会真正地向指定的io.Writer输出。
type Buffer struct {
	size   int       // 最大的缓存数量
	buffer [][]byte  // 缓存的内容
	w      io.Writer // 输出的io.Writer
}

// 新建一个Buffer。
// w最终输出的方向；当size<=1时，所有的内容都不会缓存，直接向w输出。
func NewBuffer(w io.Writer, size int) *Buffer {
	return &Buffer{size: size, w: w, buffer: make([][]byte, 0, size)}
}

// Adder.Add()
func (b *Buffer) Add(w io.Writer) error {
	if b.w == nil {
		b.w = w
		return nil
	}

	if ws, ok := b.w.(WriteFlushAdder); ok {
		return ws.Add(w)
	}

	if ws, ok := w.(WriteFlushAdder); ok {
		ws.Add(b.w)
		b.w = ws
		return nil
	}

	b.w = NewContainer(b.w, w)
	return nil
}

// io.Writer.Write()
func (b *Buffer) Write(bs []byte) (int, error) {
	if b.size <= 1 {
		return b.w.Write(bs)
	}

	b.buffer = append(b.buffer, bs)

	if len(b.buffer) < b.size {
		return len(bs), nil
	}

	return b.Flush()
}

// Flusher.Flush()
func (b *Buffer) Flush() (size int, err error) {
	if b.w == nil {
		return 0, errors.New("并未指定输出环境，b.w指向空值")
	}

	for _, buf := range b.buffer {
		if size, err = b.w.Write(buf); err != nil {
			return
		}
	}

	b.buffer = b.buffer[:0]
	return
}
