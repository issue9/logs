// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

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

// io.Writer + Flusher
type WriteFlusher interface {
	Flusher
	io.Writer
}

// io.Writer + Adder
type WriteAdder interface {
	Adder
	io.Writer
}

// Flusher + Adder
type FlushAdder interface {
	Flusher
	Adder
}

// io.Writer + Flusher + Adder
type WriteFlushAdder interface {
	io.Writer
	Flusher
	Adder
}
