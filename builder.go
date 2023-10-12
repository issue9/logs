// SPDX-License-Identifier: MIT

package logs

import "sync"

const buffersPoolMaxSize = 1 << 10

var buffersPool = &sync.Pool{New: func() any {
	b := make([]byte, 0, 1024)
	return (*Buffer)(&b)
}}

type Buffer []byte

func newBuffer() *Buffer { return buffersPool.Get().(*Buffer) }

func (w HandlerFunc) Handle(e *Record) { w(e) }

func (w *Buffer) WString(s string) *Buffer {
	*w = append(*w, s...)
	return w
}

func (w *Buffer) WBytes(b ...byte) *Buffer {
	*w = append(*w, b...)
	return w
}

func (w *Buffer) Reset() *Buffer {
	*w = (*w)[:0]
	return w
}

func (w *Buffer) Bytes() []byte { return []byte(*w) }
