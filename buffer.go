// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"sync"
)

var buffersPool = &sync.Pool{New: func() any {
	b := make([]byte, 0, 1024)
	return (*Buffer)(&b)
}}

type Buffer []byte

func NewBuffer() *Buffer { return buffersPool.Get().(*Buffer) }

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

func (w *Buffer) Write(b []byte) (int, error) {
	w.WBytes(b...)
	return len(b), nil
}

func (w *Buffer) Print(v ...any) { *w = fmt.Append(*w, v...) }

func (w *Buffer) Printf(f string, v ...any) { *w = fmt.Appendf(*w, f, v...) }

func (w *Buffer) Println(v ...any) { *w = fmt.Appendln(*w, v...) }

func (w *Buffer) Detail() bool { return true }

func (w *Buffer) Bytes() []byte { return []byte(*w) }

func (w *Buffer) Free() {
	const buffersPoolMaxSize = 1 << 10
	if len(*w) < buffersPoolMaxSize {
		buffersPool.Put(w)
	}
}
