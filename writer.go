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

const (
	MilliLayout = "15:04:05.000"
	MicroLayout = "15:04:05.000000"
	NanoLayout  = "15:04:05.000000000"
)

var nop = &nopWriter{}

type (
	// Writer 将 [Record] 转换成文本并输出的功能
	Writer interface {
		// WriteRecord 将 [Record] 写入日志通道
		//
		// NOTE: 此方法应该保证以换行符结尾。
		WriteRecord(*Record)
	}

	WriteRecord func(*Record)

	textWriter struct {
		timeLayout string
		w          io.Writer
	}

	jsonWriter struct {
		timeLayout string
		w          io.Writer
	}

	termWriter struct {
		timeLayout string
		fore       colors.Color
		w          *colors.Colorize
	}

	nopWriter struct{}

	ws []io.Writer
)

func (w WriteRecord) WriteRecord(e *Record) { w(e) }

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
	return &textWriter{timeLayout: timeLayout, w: ww}
}

func (w *textWriter) WriteRecord(e *Record) {
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

	// 一次性写入，性能更好一些。
	if _, err := w.w.Write([]byte(b.String())); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// NewJSONWriter 声明 JSON 格式的输出
func NewJSONWriter(timeLayout string, w ...io.Writer) Writer {
	var ww io.Writer
	switch len(w) {
	case 0:
		panic("参数 w 不能为空")
	case 1:
		ww = w[0]
	default:
		ww = ws(w)
	}

	return &jsonWriter{timeLayout: timeLayout, w: ww}
}

func (w *jsonWriter) WriteRecord(e *Record) {
	b := errwrap.StringBuilder{}
	b.WByte('{')

	b.WString(`"level":"`).WString(e.Level.String()).WString(`",`).
		WString(`"message":"`).WString(e.Message).WByte('"')

	if e.Logs().HasCreated() {
		b.WString(`,"created":"`).WString(e.Created.Format(w.timeLayout)).WByte('"')
	}

	if e.Logs().HasCaller() {
		b.WString(`,"path":"`).WString(e.Path).WString(`",`).
			WString(`"line":`).WString(strconv.Itoa(e.Line))
	}

	if len(e.Params) > 0 {
		b.WString(`,"params":[`)

		for i, p := range e.Params {
			val, err := json.Marshal(p.V)
			if err != nil { // 理论上不应该触发此错误，所以直接 panic
				panic(err)
			}

			if i > 0 {
				b.WByte(',')
			}
			b.WString(`{"`).WString(p.K).WString(`":`).WBytes(val).WByte('}')
		}

		b.WByte(']')
	}

	b.WByte('}')

	if _, err := w.w.Write([]byte(b.String())); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// NewTermWriter 带颜色的终端输出通道
//
// timeLayout 表示输出的时间格式，遵守 time.Format 的参数要求，
// 如果为空，则不输出时间信息；
// fore 表示终端信息的字符颜色，背景始终是默认色；
// w 表示终端的接口，可以是 [os.Stderr] 或是 [os.Stdout]，
// 如果是其它的实现者则会带控制字符一起输出；
func NewTermWriter(timeLayout string, fore colors.Color, w io.Writer) Writer {
	return &termWriter{
		timeLayout: timeLayout,
		fore:       fore,
		w:          colors.New(w),
	}
}

func (w *termWriter) WriteRecord(e *Record) {
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
func NewDispatchWriter(d map[Level]Writer) Writer {
	return WriteRecord(func(e *Record) { d[e.Level].WriteRecord(e) })
}

// NewNopWriter 空的 Writer 接口实现
func NewNopWriter() Writer { return nop }

func (w *nopWriter) WriteRecord(_ *Record) {}

// MergeWriter 将多个 Writer 合并成一个 Writer 接口对象
func MergeWriter(w ...Writer) Writer {
	return WriteRecord(func(e *Record) {
		for _, ww := range w {
			ww.WriteRecord(e)
		}
	})
}

func (w ws) Write(data []byte) (n int, err error) {
	for _, writer := range w {
		if n, err = writer.Write(data); err != nil {
			return n, err
		}
	}
	return n, err
}
