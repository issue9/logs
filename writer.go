// SPDX-License-Identifier: MIT

package logs

import (
	"io"

	"github.com/issue9/errwrap"
	"github.com/issue9/term/v3/colors"
)

type (
	// FormatFunc 格式化输出数据
	//
	// 如果格式化出错，应该返回一个基本的信息。
	FormatFunc func(*Entry) []byte

	// Writer 将 Entry 转换成文本并输出的功能
	Writer interface {
		WriteEntry(*Entry)
	}

	singleWriter struct {
		format FormatFunc
		w      io.Writer
	}

	multiWriter struct {
		format FormatFunc
		ws     []io.Writer
	}

	writers struct {
		ws []Writer
	}

	termWriter struct {
		fore colors.Color
		w    *colors.Colorize
	}
)

func TextFormat(e *Entry) []byte {
	b := errwrap.StringBuilder{}
	b.WByte('[').WString(e.Level.String()).WString("] ")

	// TODO Time ...

	return []byte(b.String())
}

func JSONFormat(e *Entry) []byte {
	//

	return nil
}

func NewWriter(f FormatFunc, w ...io.Writer) Writer {
	switch len(w) {
	case 0:
		panic("参数 w 不能为空")
	case 1:
		return &singleWriter{format: f, w: w[0]}
	default:
		return &multiWriter{format: f, ws: w}
	}
}

func (w *singleWriter) WriteEntry(e *Entry) {
	data := w.format(e)
	w.w.Write(data)
}

func (w *multiWriter) WriteEntry(e *Entry) {
	data := w.format(e)
	for _, ww := range w.ws {
		ww.Write(data)
	}
}

// MergeWriter 将多个 Writer 合并成一个 Writer 接口对象
func MergeWriter(w ...Writer) Writer { return &writers{ws: w} }

func (w *writers) WriteEntry(e *Entry) {
	for _, ww := range w.ws {
		ww.WriteEntry(e)
	}
}

func NewTermWriter(fore colors.Color, w io.Writer) Writer {
	return &termWriter{fore: fore, w: colors.New(w)}
}

func (w *termWriter) WriteEntry(e *Entry) {
	w.w.WByte('[').Color(colors.Normal, w.fore, colors.Default).WString(e.Level.String()).Reset().WString("] ") // [WARN]

	if !e.Created.IsZero() {
		w.w.Write([]byte(e.Created.Format("15:16:0000")))
	}

	// caller

	w.w.Write([]byte(e.Message))

	// pairs

}
