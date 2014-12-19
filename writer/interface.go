// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writer

import (
	"io"
)

// 缓存接口
//
// go中并没有提供析构函数的机制，所以想要在对象销毁时自动输出缓存中的内容，
// 只能定义一个接口。
type Flusher interface {

	// 将缓存的内容输出
	Flush() (size int, err error)
}

// 通过io.Writer.Write()写入缓存；
// 通过Flusher.Flush()输出缓存内容。
// io.Writer.Write()在缓存满时，也应该能调用Flusher.Flush()自动输出缓存的内容。
type WriteFlusher interface {
	Flusher
	io.Writer
}

// io.Writer的容器。
type Adder interface {
	// 向容器添加一个io.Writer实例
	Add(io.Writer) error
}

// 通过io.Writer.Write()输入的内容，会自动各容器中所有的io.Writer实例输出。
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
