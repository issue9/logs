// SPDX-License-Identifier: MIT

package logs

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/issue9/errwrap"
	"github.com/issue9/term/v3/colors"

	"github.com/issue9/logs/v5/writers"
)

// 常用的日志时间格式
const (
	DateMilliLayout = "2006-01-02T15:04:05.000"
	DateMicroLayout = "2006-01-02T15:04:05.000000"
	DateNanoLayout  = "2006-01-02T15:04:05.000000000"

	MilliLayout = "15:04:05.000"
	MicroLayout = "15:04:05.000000"
	NanoLayout  = "15:04:05.000000000"
)

var nop = &nopHandler{}

var buffersPool = &sync.Pool{New: func() any { return &errwrap.Buffer{} }}

var defaultTermColors = map[Level]colors.Color{
	LevelInfo:  colors.Green,
	LevelDebug: colors.Yellow,
	LevelTrace: colors.Yellow,
	LevelWarn:  colors.Yellow,
	LevelError: colors.Red,
	LevelFatal: colors.Red,
}

type (
	// Handler [Record] 的处理接口
	Handler interface {
		// Handle 将 [Record] 写入日志
		//
		// [Record] 中各个字段的名称由处理器自行决定。
		//
		// NOTE: 此方法应该保证输出内容是以换行符作为结尾。
		Handle(*Record)
	}

	HandleFunc func(*Record)

	nopHandler struct{}
)

func (w HandleFunc) Handle(e *Record) { w(e) }

// NewTextHandler 返回将 [Record] 以普通文本的形式写入 w 的对象
//
// NOTE: 如果向 w 输出内容时出错，会将错误信息输出到终端作为最后的处理方式。
func NewTextHandler(timeLayout string, w ...io.Writer) Handler {
	ww := writers.New(w...)
	mux := &sync.Mutex{} // 防止多个函数同时调用 HandleFunc 方法。

	return HandleFunc(func(e *Record) {
		b := buffersPool.Get().(*errwrap.Buffer)
		b.Reset()

		b.WByte('[').WString(e.Level.String()).WByte(']')

		var indent byte = ' '
		if e.Logs().HasCreated() {
			b.WByte(' ').WString(e.Created.Format(timeLayout))
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
				if bs, err := m.MarshalText(); err != nil {
					b.WString("Err(").WString(err.Error()).WByte(')')
				} else {
					b.WBytes(bs)
				}
			} else {
				b.Print(p.V)
			}
		}

		b.WByte('\n')

		mux.Lock()
		defer mux.Unlock()
		if _, err := ww.Write(b.Bytes()); err != nil { // 一次性写入，性能更好一些。
			fmt.Fprintf(os.Stderr, "NewTextHandler.Handle:%v\n", err)
		}
		buffersPool.Put(b)
	})
}

// NewJSONHandler 返回将 [Record] 以 JSON 的形式写入 w 的对象
//
// NOTE: 如果向 w 输出内容时出错，会将错误信息输出到终端作为最后的处理方式。
func NewJSONHandler(timeLayout string, w ...io.Writer) Handler {
	ww := writers.New(w...)
	mux := &sync.Mutex{}

	return HandleFunc(func(e *Record) {
		b := buffersPool.Get().(*errwrap.Buffer)
		b.Reset()

		b.WByte('{')

		b.WString(`"level":"`).WString(e.Level.String()).WString(`",`).
			WString(`"message":"`).WString(e.Message).WByte('"')

		if e.Logs().HasCreated() {
			b.WString(`,"created":"`).WString(e.Created.Format(timeLayout)).WByte('"')
		}

		if e.Logs().HasCaller() {
			b.WString(`,"path":"`).WString(e.Path).WString(`",`).
				WString(`"line":`).WString(strconv.Itoa(e.Line))
		}

		if len(e.Params) > 0 {
			b.WString(`,"params":[`)

			for i, p := range e.Params {
				val, err := json.Marshal(p.V) // TODO 基本类型直接处理，是不是会更快一些？
				if err != nil {
					val = []byte("\"Err(" + err.Error() + ")\"")
				}

				if i > 0 {
					b.WByte(',')
				}
				b.WString(`{"`).WString(p.K).WString(`":`).WBytes(val).WByte('}')
			}

			b.WByte(']')
		}

		b.WByte('}')

		mux.Lock()
		defer mux.Unlock()
		if _, err := ww.Write(b.Bytes()); err != nil {
			fmt.Fprintf(os.Stderr, "NewJSONHandler.Handle:%v\n", err)
		}
		buffersPool.Put(b)
	})
}

// NewTermHandler 返回将 [Record] 写入终端的对象
//
// timeLayout 表示输出的时间格式，遵守 time.Format 的参数要求；
// w 表示终端的接口，可以是 [os.Stderr] 或是 [os.Stdout]，
// 如果是其它的实现者则会带控制字符一起输出；
// foreColors 表示各类别信息的字符颜色，背景始终是默认色，未指定的颜色会从 [defaultTermColors] 获取；
//
// NOTE: 如果向 w 输出内容时出错，将会导致 panic。
func NewTermHandler(timeLayout string, w io.Writer, foreColors map[Level]colors.Color) Handler {
	cs := make(map[Level]colors.Color, len(defaultTermColors))
	for l, cc := range defaultTermColors {
		if c, found := foreColors[l]; found {
			cs[l] = c
		} else {
			cs[l] = cc
		}
	}

	mux := &sync.Mutex{}

	return HandleFunc(func(e *Record) {
		b := buffersPool.Get().(*errwrap.Buffer)
		b.Reset()
		ww := colors.New(b)

		fc := cs[e.Level]
		ww.WByte('[').Color(colors.Normal, fc, colors.Default).WString(e.Level.String()).Reset().WByte(']') // [WARN]

		var indent byte = ' '
		if e.Logs().HasCreated() {
			ww.WByte(' ').WString(e.Created.Format(timeLayout))
			indent = '\t'
		}

		if e.Logs().HasCaller() {
			ww.WByte(' ').WString(e.Path).WByte(':').WString(strconv.Itoa(e.Line))
			indent = '\t'
		}

		ww.WByte(indent).WString(e.Message)

		for _, p := range e.Params {
			ww.WByte(' ').WString(p.K).WByte('=').WString(fmt.Sprint(p.V))
		}

		ww.WByte('\n')

		mux.Lock()
		defer mux.Unlock()
		if _, err := w.Write(b.Bytes()); err != nil {
			// 大概率是写入终端失败，直接 panic。
			panic(fmt.Sprintf("NewTermHandler.Handle:%v\n", err))
		}
		buffersPool.Put(b)
	})
}

// NewDispatchHandler 根据 [Level] 派发到不同的 [Handler] 对象
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
