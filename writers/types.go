// SPDX-License-Identifier: MIT

package writers

import "io"

// Adder 为 io.Writer 的容器
type Adder interface {
	// 向容器添加一个 io.Writer 实例
	Add(io.Writer) error
}

// Flusher 将缓存内容输出到日志
type Flusher interface {
	// 将缓存的内容输出
	Flush() (size int, err error)
}

type WriteFlusher interface {
	Flusher
	io.Writer
}

type WriteAdder interface {
	Adder
	io.Writer
}

type FlushAdder interface {
	Flusher
	Adder
}

type WriteFlushAdder interface {
	io.Writer
	Flusher
	Adder
}
