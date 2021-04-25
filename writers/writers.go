// SPDX-License-Identifier: MIT

// Package writers 提供了一组实现 io.Writer 接口的结构
package writers

import "io"

type Adder interface {
	// Add 向容器添加一个 io.Writer 实例
	Add(io.Writer) error
}

type Flusher interface {
	// Flush 将缓存的内容输出
	Flush() (err error)
}
