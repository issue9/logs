// SPDX-License-Identifier: MIT

package writers

import (
	"io"
)

// Adder 为 io.Writer 的容器。
type Adder interface {
	// 向容器添加一个 io.Writer 实例
	Add(io.Writer) error
}

// Flusher 缓存接口
//
// go 中并没有提供析构函数的机制，所以想要在对象销毁时自动输出缓存中的内容，
// 只能定义一个接口。
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
