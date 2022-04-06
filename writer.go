// SPDX-License-Identifier: MIT

package logs

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/issue9/errwrap"
	"github.com/issue9/term/v3/colors"
)

var nop = &nopWriter{}

type (
	// Writer 将 Entry 转换成文本并输出的功能
	Writer interface {
		// WriteEntry 将 Entry 写入日志通道
		//
		// NOTE: 此方法应该保证以换行符结尾。
		WriteEntry(*Entry)
	}

	textWriter struct {
		timeLayout string
		b          io.Writer
	}

	jsonWriter struct {
		enc *json.Encoder
	}

	termWriter struct {
		timeLayout string
		fore       colors.Color
		w          *colors.Colorize
	}

	dispatchWriter map[Level]Writer

	nopWriter struct{}

	writers struct {
		ws []Writer
	}

	ws []io.Writer
)

func NewTextWriter(timeLayout string, w ...io.Writer) Writer {
	var ww io.Writer
	switch len(w) {
	case 0:
		panic("参数 w 不能为空")
	case 1:
		ww = w[0]
	default:
		ww = ws(w)
	}
	return &textWriter{timeLayout: timeLayout, b: ww}
}

func (w *textWriter) WriteEntry(e *Entry) {
	b := errwrap.StringBuilder{}
	b.WByte('[').WString(e.Level.String()).WByte(']')

	var indent byte = ' '
	if e.Logs().HasCreated() {
		b.WByte(' ').WString(e.Created.Format(w.timeLayout))
		indent = '\t'
	}

	if e.Logs().HasCaller() {
		b.WByte(' ').WString(e.Path).WByte(':').WString(strconv.Itoa(e.Line))
		indent = '\t'
	}

	b.WByte(indent).WString(e.Message)

	for _, p := range e.Params {
		b.WByte(' ').WString(p.K).WByte('=').WString(fmt.Sprint(p.V))
	}

	b.WByte('\n')

	if _, err := w.b.Write([]byte(b.String())); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// NewJSONWriter 声明 JSON 格式的输出
//
// format 是否格式化 json，如果为 true，会以 \t 作为缩进；
func NewJSONWriter(format bool, w ...io.Writer) Writer {
	var ww io.Writer
	switch len(w) {
	case 0:
		panic("参数 w 不能为空")
	case 1:
		ww = w[0]
	default:
		ww = ws(w)
	}

	enc := json.NewEncoder(ww)
	if format {
		enc.SetIndent("", "\t")
	}
	return &jsonWriter{enc: enc}
}

func (w *jsonWriter) WriteEntry(e *Entry) {
	if err := w.enc.Encode(e); err != nil {
		fmt.Fprintln(os.Stderr, err) // 编码错误
	}
}

// NewTermWriter 带颜色的终端输出通道
//
// timeLayout 表示输出的时间格式，遵守 time.Format 的参数要求，如果为空，则不输出时间信息；
// fore 表示终端信息的字符颜色，背景始终是默认色；
// w 表示终端的接口，可以是 os.Stderr 或是 os.Stdout，如果是其它的实现者则会带控制字符一起输出；
func NewTermWriter(timeLayout string, fore colors.Color, w io.Writer) Writer {
	return &termWriter{
		timeLayout: timeLayout,
		fore:       fore,
		w:          colors.New(w),
	}
}

func (w *termWriter) WriteEntry(e *Entry) {
	w.w.WByte('[').Color(colors.Normal, w.fore, colors.Default).WString(e.Level.String()).Reset().WByte(']') // [WARN]

	var indent byte = ' '
	if e.Logs().HasCreated() {
		w.w.WByte(' ').WString(e.Created.Format(w.timeLayout))
		indent = '\t'
	}

	if e.Logs().HasCaller() {
		w.w.WByte(' ').WString(e.Path).WByte(':').WString(strconv.Itoa(e.Line))
		indent = '\t'
	}

	w.w.WByte(indent).WString(e.Message)

	for _, p := range e.Params {
		w.w.WByte(' ').WString(p.K).WByte('=').WString(fmt.Sprint(p.V))
	}

	w.w.WByte('\n')
}

// NewDispatchWriter 根据 Level 派发到不同的 Writer 对象
func NewDispatchWriter(d map[Level]Writer) Writer { return dispatchWriter(d) }

func (w dispatchWriter) WriteEntry(e *Entry) { w[e.Level].WriteEntry(e) }

// NewNopWriter 空的 Writer 接口实现
func NewNopWriter() Writer { return nop }

func (w *nopWriter) WriteEntry(_ *Entry) {}

// MergeWriter 将多个 Writer 合并成一个 Writer 接口对象
func MergeWriter(w ...Writer) Writer { return &writers{ws: w} }

func (w *writers) WriteEntry(e *Entry) {
	for _, ww := range w.ws {
		ww.WriteEntry(e)
	}
}

func (w ws) Write(data []byte) (n int, err error) {
	for _, writer := range w {
		if n, err = writer.Write(data); err != nil {
			return n, err
		}
	}
	return n, err
}
