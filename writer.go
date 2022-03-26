// SPDX-License-Identifier: MIT

package logs

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

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
		// WriteEntry 将 Entry 写入日志通道
		//
		// NOTE: 此方法应该保证以换行符结尾。
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
		timeLayout string
		fore       colors.Color
		w          *colors.Colorize
	}
)

func TextFormat(timeLayout string) FormatFunc {
	return func(e *Entry) []byte {
		b := errwrap.StringBuilder{}
		b.WByte('[').WString(e.Level.String()).WString("] ")

		if timeLayout != "" {
			b.WString(e.Created.Format(timeLayout)).WByte(' ')
		}

		b.WString(e.Message).WByte('\t')

		b.WString(e.Path).WByte(':').WString(strconv.Itoa(e.Line))

		for _, p := range e.Pairs {
			b.WByte(' ').WString(p.K).WByte('=').WString(fmt.Sprint(p.V))
		}

		return []byte(b.String())
	}
}

// JSONFormat 格式化 JSON 输出
//
// 如果无法转换成 JSON，那么将退化成 TextFormat 输出。
func JSONFormat(e *Entry) []byte {
	m := make(map[string]any, len(e.Pairs)+3)

	m["message"] = e.Message
	m["level"] = e.Level
	m["created"] = e.Created
	m["location"] = e.Path + ":" + strconv.Itoa(e.Line)
	for _, p := range e.Pairs {
		m[p.K] = p.V
	}

	data, err := json.Marshal(m)
	if err != nil {
		return TextFormat(time.RFC3339)(e)
	}
	return data
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
	data := append(w.format(e), '\n')
	w.w.Write(data)
}

func (w *multiWriter) WriteEntry(e *Entry) {
	data := append(w.format(e), '\n')
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
	w.w.WByte('[').Color(colors.Normal, w.fore, colors.Default).WString(e.Level.String()).Reset().WString("] ") // [WARN]

	if w.timeLayout != "" {
		w.w.WString(e.Created.Format(w.timeLayout)).WByte(' ')
	}

	w.w.WString(e.Message).WByte('\t')

	w.w.WString(e.Path).WByte(':').WString(strconv.Itoa(e.Line))

	for _, p := range e.Pairs {
		w.w.WByte(' ').WString(p.K).WByte('=').WString(fmt.Sprint(p.V))
	}

	w.w.WByte('\n')
}
