// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var buffersPool = &sync.Pool{New: func() any {
	return &Buffer{data: make([]byte, 0, 1024)}
}}

// Buffer []byte 复用对象池
//
// 同时实现了 [xerrors.Printer] 接口。
type Buffer struct {
	data   []byte
	detail bool
}

type AppendFunc = func(*Buffer)

// NewBuffer 声明 [Buffer] 对象
//
// detail 是否打印错误信息的调用堆栈；
func NewBuffer(detail bool) *Buffer { return buffersPool.Get().(*Buffer).Reset(detail) }

func (w *Buffer) Reset(detail bool) *Buffer {
	w.data = w.data[:0]
	w.detail = detail
	return w
}

func (w *Buffer) Write(b []byte) (int, error) {
	w.AppendBytes(b...)
	return len(b), nil
}

func (w *Buffer) AppendFunc(f AppendFunc) *Buffer {
	f(w)
	return w
}

func (w *Buffer) AppendString(s string) *Buffer {
	w.data = append(w.data, s...)
	return w
}

func (w *Buffer) AppendBytes(b ...byte) *Buffer {
	w.data = append(w.data, b...)
	return w
}

func (w *Buffer) AppendFloat(n float64, fmt byte, prec, bitSize int) *Buffer {
	w.data = strconv.AppendFloat(w.data, n, fmt, prec, bitSize)
	return w
}

func (w *Buffer) AppendInt(n int64, base int) *Buffer {
	w.data = strconv.AppendInt(w.data, n, 10)
	return w
}

func (w *Buffer) AppendUint(n uint64, base int) *Buffer {
	w.data = strconv.AppendUint(w.data, n, 10)
	return w
}

func (w *Buffer) AppendTime(t time.Time, layout string) *Buffer {
	w.data = t.AppendFormat(w.data, layout)
	return w
}

func (w *Buffer) Append(v ...any) *Buffer {
	w.data = fmt.Append(w.data, v...)
	return w
}

func (w *Buffer) Appendf(format string, v ...any) *Buffer {
	w.data = fmt.Appendf(w.data, format, v...)
	return w
}

func (w *Buffer) Appendln(v ...any) *Buffer {
	w.data = fmt.Appendln(w.data, v...)
	return w
}

func (w *Buffer) AppendBuffer(f func(b *Buffer)) *Buffer {
	bb := NewBuffer(w.detail)
	defer bb.Free()
	f(bb)

	return w.AppendBytes(bb.data...)
}

func (w *Buffer) Print(v ...any) { w.data = fmt.Append(w.data, v...) }

func (w *Buffer) Printf(f string, v ...any) { w.data = fmt.Appendf(w.data, f, v...) }

func (w *Buffer) Println(v ...any) { w.data = fmt.Appendln(w.data, v...) }

func (w *Buffer) Detail() bool { return w.detail }

func (w *Buffer) Bytes() []byte { return w.data }

func (w *Buffer) Free() {
	const buffersPoolMaxSize = 1 << 10
	if len(w.data) < buffersPoolMaxSize {
		buffersPool.Put(w)
	}
}
