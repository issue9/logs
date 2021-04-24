// SPDX-License-Identifier: MIT

// Package writers 提供了一组实现 io.Writer 接口的结构
package writers

import "io"

// Adder 为 io.Writer 的容器
type Adder interface {
	// Add 向容器添加一个 io.Writer 实例
	Add(io.Writer) error
}

// Flusher 将缓存内容输出到日志
type Flusher interface {
	// Flush 将缓存的内容输出
	Flush() (err error)
}

type WriteFlushAdder interface {
	io.Writer
	Flusher
	Adder
}
