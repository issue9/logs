// SPDX-License-Identifier: MIT

// Package writers 提供了一组实现 [io.Writer] 接口的结构
package writers

import "io"

type (
	ws []io.Writer

	WriteFunc func([]byte) (int, error)
)

func (f WriteFunc) Write(data []byte) (int, error) { return f(data) }

func (w ws) Write(data []byte) (n int, err error) {
	for _, writer := range w {
		if n, err = writer.Write(data); err != nil {
			return n, err
		}
	}
	return n, err
}

// New 将 [1,n] 个 [io.Writer] 合并成一个
func New(w ...io.Writer) io.Writer {
	switch len(w) {
	case 0:
		panic("参数 w 不能为空")
	case 1:
		return w[0]
	default:
		return ws(w)
	}
}
