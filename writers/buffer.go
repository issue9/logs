// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writers

import (
	"io"
)

// Buffer 实现对输出内容的缓存，只有输出数量达到指定的值
// 才会真正地向指定的io.Writer输出。
type Buffer struct {
	size   int         // 最大的缓存数量
	buffer [][]byte    // 缓存的内容
	ws     []io.Writer // 输出的io.Writer
}

// 新建一个Buffer。
// w最终输出的方向；当size<=1时，所有的内容都不会缓存，直接向w输出。
func NewBuffer(size int, ws ...io.Writer) *Buffer {
	return &Buffer{size: size, ws: ws, buffer: make([][]byte, 0, size)}
}

// Adder.Add()
func (b *Buffer) Add(w io.Writer) error {
	b.ws = append(b.ws, w)
	return nil
}

// io.Writer.Write()
// 若容器为空时，则相当于不作任何动作。
func (b *Buffer) Write(bs []byte) (int, error) {
	if b.size < 2 {
		return b.write(bs)
	}

	b.buffer = append(b.buffer, bs)

	if len(b.buffer) < b.size {
		return len(bs), nil
	}

	return b.Flush()
}

// Flusher.Flush()
// 若容器为空时，则相当于不作任何动作。
func (b *Buffer) Flush() (size int, err error) {
	for _, buf := range b.buffer {
		if size, err = b.write(buf); err != nil {
			return
		}
	}

	b.buffer = b.buffer[:0]
	return
}

// 设置缓存的大小，若值小于2，则所有的输出都不会被缓存。
func (b *Buffer) SetSize(size int) {
	b.size = size
}

func (b *Buffer) write(bs []byte) (size int, err error) {
	for _, w := range b.ws {
		if size, err = w.Write(bs); err != nil {
			return
		}
	}
	return
}
