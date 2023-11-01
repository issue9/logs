// SPDX-License-Identifier: MIT

// Package writers 提供了一组实现 [io.Writer] 接口的结构
package writers

import "io"

type WriteFunc func([]byte) (int, error)

func (f WriteFunc) Write(data []byte) (int, error) { return f(data) }

// New 将 [1,n] 个 [io.Writer] 合并成一个
//
// Deprecated: 请使用 [io.MultiWriter] 代替
func New(w ...io.Writer) io.Writer { return io.MultiWriter(w...) }
