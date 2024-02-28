// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package writers 提供了一组实现 [io.Writer] 接口的结构
package writers

import "io"

type WriteFunc func([]byte) (int, error)

func (f WriteFunc) Write(data []byte) (int, error) { return f(data) }

func New(w ...io.Writer) io.Writer {
	ws := make([]io.Writer, 0, len(w))
	for _, ww := range w {
		if ww != nil {
			ws = append(ws, ww)
		}
	}

	if len(ws) == 0 {
		panic("参数 w 为空或是所有值均为 nil")
	}

	return io.MultiWriter(ws...)
}
