// SPDX-License-Identifier: MIT

// Package writers 提供了一组实现 [io.Writer] 接口的结构
package writers

type WriteFunc func([]byte) (int, error)

func (f WriteFunc) Write(data []byte) (int, error) { return f(data) }
