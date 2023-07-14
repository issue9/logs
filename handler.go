// SPDX-License-Identifier: MIT

package logs

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/issue9/errwrap"
	"github.com/issue9/term/v3/colors"
)

const (
	MilliLayout = "15:04:05.000"
	MicroLayout = "15:04:05.000000"
	NanoLayout  = "15:04:05.000000000"
)

var nop = &nopHandler{}

type (
	// Handler 将 [Record] 转换成文本并输出的功能
	Handler interface {
		// Handle 将 [Record] 写入日志通道
		//
		// NOTE: 此方法应该保证以换行符结尾。
		Handle(*Record)
	}

	HandleFunc func(*Record)

	textHandler struct {
		timeLayout string
		w          io.Writer
	}

	jsonHandler struct {
		timeLayout string
		w          io.Writer
	}

	termHandler struct {
		mux        sync.Mutex
		timeLayout string
		fore       colors.Color
		w          *colors.Colorize
	}

	nopHandler struct{}

	ws []io.Writer
)

func (w HandleFunc) Handle(e *Record) { w(e) }

func NewTextHandler(timeLayout string, w ...io.Writer) Handler {
	var ww io.Writer
	switch len(w) {
	case 0:
		panic("参数 w 不能为空")
	case 1:
		ww = w[0]
	default:
		ww = ws(w)
	}
	return &textHandler{timeLayout: timeLayout, w: ww}
}

func (w *textHandler) Handle(e *Record) {
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
		b.WByte(' ').WString(p.K).WByte('=')
		if m, ok := p.V.(encoding.TextMarshaler); ok {
			bs, err := m.MarshalText()
			if err != nil {
				log.Panicln("TextHandler.Handle:", err)
			} else {
				b.WBytes(bs)
			}
		}
		b.Print(p.V)
	}

	b.WByte('\n')

	// 一次性写入，性能更好一些。
	// NOTE: 单次写入整条记录，否则需要用锁确保写入数据的完整性。
	if _, err := w.w.Write([]byte(b.String())); err != nil {
		log.Println("JSONHandler.Handle:", err)
	}
}

// NewJSONHandler 声明 JSON 格式的输出
func NewJSONHandler(timeLayout string, w ...io.Writer) Handler {
	var ww io.Writer
	switch len(w) {
	case 0:
		panic("参数 w 不能为空")
	case 1:
		ww = w[0]
	default:
		ww = ws(w)
	}

	return &jsonHandler{timeLayout: timeLayout, w: ww}
}

func (w *jsonHandler) Handle(e *Record) {
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
			if err != nil {
				log.Println("JSONHandler.Handle:", err)
			}

			if i > 0 {
				b.WByte(',')
			}
			b.WString(`{"`).WString(p.K).WString(`":`).WBytes(val).WByte('}')
		}

		b.WByte(']')
	}

	b.WByte('}')

	// NOTE: 单次写入整条记录，否则需要用锁确保写入数据的完整性。
	if _, err := w.w.Write([]byte(b.String())); err != nil {
		log.Println("JSONHandler.Handle:", err)
	}
}

// NewTermHandler 带颜色的终端输出通道
//
// timeLayout 表示输出的时间格式，遵守 time.Format 的参数要求，
// 如果为空，则不输出时间信息；
// fore 表示终端信息的字符颜色，背景始终是默认色；
// w 表示终端的接口，可以是 [os.Stderr] 或是 [os.Stdout]，
// 如果是其它的实现者则会带控制字符一起输出；
func NewTermHandler(timeLayout string, fore colors.Color, w io.Writer) Handler {
	return &termHandler{
		timeLayout: timeLayout,
		fore:       fore,
		w:          colors.New(w),
	}
}

func (w *termHandler) Handle(e *Record) {
	w.mux.Lock()
	defer w.mux.Unlock()

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

// NewDispatchHandler 根据 Level 派发到不同的 Writer 对象
func NewDispatchHandler(d map[Level]Handler) Handler {
	return HandleFunc(func(e *Record) { d[e.Level].Handle(e) })
}

// NewNopHandler 空的 Handler 接口实现
func NewNopHandler() Handler { return nop }

func (w *nopHandler) Handle(_ *Record) {}

// MergeHandler 将多个 Handler 合并成一个 Handler 接口对象
func MergeHandler(w ...Handler) Handler {
	return HandleFunc(func(e *Record) {
		for _, ww := range w {
			ww.Handle(e)
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
