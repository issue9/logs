// SPDX-License-Identifier: MIT

package writers

import "io"

// Buffer 实现对输出内容的缓存
//
// 只有输出数量达到指定的值才会真正地向指定的 io.Writer 输出。
type Buffer struct {
	c      *Container
	size   int      // 最大的缓存数量
	buffer [][]byte // 缓存的内容
}

// NewBuffer 新建一个 Buffer
//
// 当 size<=1 时，所有的内容都不会被缓存。
func NewBuffer(size int) *Buffer {
	return &Buffer{size: size,
		c:      NewContainer(),
		buffer: make([][]byte, 0, size),
	}
}

// Add 添加一个 io.Writer 实例
func (b *Buffer) Add(w io.Writer) error {
	return b.c.Add(w)
}

func (b *Buffer) Write(bs []byte) (int, error) {
	if b.size < 2 {
		return b.c.Write(bs)
	}

	// 参数 bs 来源于 log.Logger.buf，该变量会被 log.Logger 不断
	// 重复使用，所以此处应该复制一份 bs 的内容再保存。
	cp := make([]byte, len(bs))
	copy(cp, bs)
	b.buffer = append(b.buffer, cp)

	if len(b.buffer) < b.size {
		return len(bs), nil
	}

	return len(bs), b.Flush()
}

// Flush 实现了 Flusher.Flush()
// 若容器为空时，则相当于不作任何动作。
func (b *Buffer) Flush() error {
	for _, buf := range b.buffer {
		if _, err := b.c.Write(buf); err != nil {
			return err
		}
	}
	b.buffer = b.buffer[:0]

	return b.c.Flush()
}

// SetSize 设置缓存的大小
//
// 若值小于 2，则所有的输出都不会被缓存。
func (b *Buffer) SetSize(size int) {
	b.size = size
}
