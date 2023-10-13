// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/issue9/term/v3/colors"

	"github.com/issue9/logs/v6/writers"
)

var nop = &nopHandler{}

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

	HandlerFunc func(*Record)

	nopHandler struct{}
)

func (w HandlerFunc) Handle(e *Record) { w(e) }

// NewTextHandler 返回将 [Record] 以普通文本的形式写入 w 的对象
//
// NOTE: 如果向 w 输出内容时出错，会将错误信息输出到终端作为最后的处理方式。
func NewTextHandler(w ...io.Writer) Handler {
	ww := writers.New(w...)
	mux := &sync.Mutex{} // 防止多个函数同时调用 HandlerFunc 方法。

	return HandlerFunc(func(e *Record) {
		b := NewBuffer()
		b.Reset()

		b.WBytes('[').WString(e.Level.String()).WBytes(']')

		var indent byte = ' '
		if !e.Created.IsZero() {
			b.WBytes(' ').WString(e.Created.Format(e.Logs().createdFormat))
			indent = '\t'
		}

		if e.Path != "" {
			b.WBytes(' ').WString(e.Path)
			indent = '\t'
		}

		b.WBytes(indent)
		*b = e.AppendMessage(b.Bytes())

		for _, p := range e.Params {
			b.WBytes(' ').WString(p.K).WBytes('=')
			switch v := p.V.(type) {
			case string:
				b.WString(v)
			case int:
				b.WString(strconv.Itoa(v))
			case int64:
				b.WString(strconv.FormatInt(v, 10))
			case int32:
				b.WString(strconv.FormatInt(int64(v), 10))
			case int16:
				b.WString(strconv.FormatInt(int64(v), 10))
			case int8:
				b.WString(strconv.FormatInt(int64(v), 10))
			case uint:
				b.WString(strconv.FormatUint(uint64(v), 10))
			case uint64:
				b.WString(strconv.FormatUint(v, 10))
			case uint32:
				b.WString(strconv.FormatUint(uint64(v), 10))
			case uint16:
				b.WString(strconv.FormatUint(uint64(v), 10))
			case uint8:
				b.WString(strconv.FormatUint(uint64(v), 10))
			case float32:
				b.WString(strconv.FormatFloat(float64(v), 'f', 5, 32))
			case float64:
				b.WString(strconv.FormatFloat(float64(v), 'f', 5, 64))
			default:
				if m, ok := p.V.(encoding.TextMarshaler); ok {
					if bs, err := m.MarshalText(); err != nil {
						b.WString("Err(").WString(err.Error()).WBytes(')')
					} else {
						b.WBytes(bs...)
					}
				} else {
					*b = fmt.Append(b.Bytes(), p.V)
				}
			}
		}

		b.WBytes('\n')

		mux.Lock()
		defer mux.Unlock()
		if _, err := ww.Write(b.Bytes()); err != nil { // 一次性写入，性能更好一些。
			fmt.Fprintf(os.Stderr, "NewTextHandler.Handle:%v\n", err)
		}

		if len(*b) < buffersPoolMaxSize {
			buffersPool.Put(b)
		}
	})
}

// NewJSONHandler 返回将 [Record] 以 JSON 的形式写入 w 的对象
//
// NOTE: 如果向 w 输出内容时出错，会将错误信息输出到终端作为最后的处理方式。
func NewJSONHandler(w ...io.Writer) Handler {
	ww := writers.New(w...)
	mux := &sync.Mutex{}

	return HandlerFunc(func(e *Record) {
		b := NewBuffer()
		b.Reset()

		b.WBytes('{')

		b.WString(`"level":"`).WString(e.Level.String()).WString(`",`).
			WString(`"message":"`)
		*b = e.AppendMessage(b.Bytes())
		b.WBytes('"')

		if !e.Created.IsZero() {
			b.WString(`,"created":"`).WString(e.Created.Format(e.Logs().createdFormat)).WBytes('"')
		}

		if e.Path != "" {
			b.WString(`,"path":"`).WString(e.Path).WBytes('"')
		}

		if len(e.Params) > 0 {
			b.WString(`,"params":[`)

			for i, p := range e.Params {
				if i > 0 {
					b.WBytes(',')
				}
				b.WString(`{"`).WString(p.K).WString(`":`)

				switch v := p.V.(type) {
				case string:
					b.WBytes('"').WString(v).WBytes('"')
				case int:
					b.WString(strconv.Itoa(v))
				case int64:
					b.WString(strconv.FormatInt(v, 10))
				case int32:
					b.WString(strconv.FormatInt(int64(v), 10))
				case int16:
					b.WString(strconv.FormatInt(int64(v), 10))
				case int8:
					b.WString(strconv.FormatInt(int64(v), 10))
				case uint:
					b.WString(strconv.FormatUint(uint64(v), 10))
				case uint64:
					b.WString(strconv.FormatUint(v, 10))
				case uint32:
					b.WString(strconv.FormatUint(uint64(v), 10))
				case uint16:
					b.WString(strconv.FormatUint(uint64(v), 10))
				case uint8:
					b.WString(strconv.FormatUint(uint64(v), 10))
				case float32:
					b.WString(strconv.FormatFloat(float64(v), 'e', 5, 32))
				case float64:
					b.WString(strconv.FormatFloat(float64(v), 'e', 5, 64))
				default:
					val, err := json.Marshal(p.V)
					if err != nil {
						val = []byte("\"Err(" + err.Error() + ")\"")
					}
					b.WBytes(val...)
				}

				b.WBytes('}')
			}

			b.WBytes(']')
		}

		b.WBytes('}')

		mux.Lock()
		defer mux.Unlock()
		if _, err := ww.Write(b.Bytes()); err != nil {
			fmt.Fprintf(os.Stderr, "NewJSONHandler.Handle:%v\n", err)
		}

		if len(*b) < buffersPoolMaxSize {
			buffersPool.Put(b)
		}
	})
}

// NewTermHandler 返回将 [Record] 写入终端的对象
//
// w 表示终端的接口，可以是 [os.Stderr] 或是 [os.Stdout]，
// 如果是其它的实现者则会带控制字符一起输出；
// foreColors 表示各类别信息的字符颜色，背景始终是默认色，未指定的颜色会从 [defaultTermColors] 获取；
//
// NOTE: 如果向 w 输出内容时出错，将会导致 panic。
func NewTermHandler(w io.Writer, foreColors map[Level]colors.Color) Handler {
	cs := make(map[Level]colors.Color, len(defaultTermColors))
	for l, cc := range defaultTermColors {
		if c, found := foreColors[l]; found {
			cs[l] = c
		} else {
			cs[l] = cc
		}
	}

	mux := &sync.Mutex{}

	return HandlerFunc(func(e *Record) {
		buf := new(bytes.Buffer)
		ww := colors.New(buf)

		fc := cs[e.Level]
		ww.WByte('[').Color(colors.Normal, fc, colors.Default).WString(e.Level.String()).Reset().WByte(']') // [WARN]

		var indent byte = ' '
		if !e.Created.IsZero() {
			ww.WByte(' ').WString(e.Created.Format(e.Logs().createdFormat))
			indent = '\t'
		}

		if e.Path != "" {
			ww.WByte(' ').WString(e.Path)
			indent = '\t'
		}

		ww.WByte(indent).WBytes(e.AppendMessage([]byte{}))

		for _, p := range e.Params {
			ww.WByte(' ').WString(p.K).WByte('=').WString(fmt.Sprint(p.V))
		}

		ww.WByte('\n')

		mux.Lock()
		defer mux.Unlock()
		if _, err := w.Write(buf.Bytes()); err != nil {
			// 大概率是写入终端失败，直接 panic。
			panic(fmt.Sprintf("NewTermHandler.Handle:%v\n", err))
		}

		//if b.Len() < buffersPoolMaxSize {
		//buffersPool.Put(buf)
		//}
	})
}

// NewDispatchHandler 根据 [Level] 派发到不同的 [Handler] 对象
func NewDispatchHandler(d map[Level]Handler) Handler {
	return HandlerFunc(func(e *Record) { d[e.Level].Handle(e) })
}

// NewNopHandler 空的 Handler 接口实现
func NewNopHandler() Handler { return nop }

func (w *nopHandler) Handle(_ *Record) {}

// MergeHandler 将多个 Handler 合并成一个 Handler 接口对象
func MergeHandler(w ...Handler) Handler {
	return HandlerFunc(func(e *Record) {
		for _, ww := range w {
			ww.Handle(e)
		}
	})
}
