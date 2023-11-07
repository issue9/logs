// SPDX-License-Identifier: MIT

package logs

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/issue9/term/v3/colors"

	"github.com/issue9/logs/v7/writers"
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
	// Handler 日志后端的处理接口
	Handler interface {
		// Handle 将 [Record] 写入日志
		//
		// [Record] 中各个字段的名称由处理器自行决定。
		//
		// NOTE: 此方法应该保证输出内容是以换行符作为结尾。
		Handle(*Record)

		// WithAttrs 根据参数生成新的 [Handler] 对象
		//
		// 新对象继承旧对象的属性，并添加了参数中的新属性。
		//
		// 即便参数的长度为零，也应该返回新的对象。
		WithAttrs([]Attr) Handler
	}

	textHandler struct {
		w     io.Writer
		mux   sync.Mutex
		attrs []byte
	}

	jsonHandler struct {
		w     io.Writer
		mux   sync.Mutex
		attrs []byte
	}

	termHandler struct {
		w          io.Writer
		foreColors map[Level]colors.Color
		mux        sync.Mutex
		attrs      []Attr
	}

	dispatchHandler struct {
		handlers map[Level]Handler
	}

	mergeHandler struct {
		handlers []Handler
	}

	nopHandler struct{}
)

// NewTextHandler 返回将 [Record] 以普通文本的形式写入 w 的对象
//
// NOTE: 如果向 w 输出内容时出错，会将错误信息输出到终端作为最后的处理方式。
func NewTextHandler(w ...io.Writer) Handler {
	return &textHandler{w: writers.New(w...)}
}

func (h *textHandler) Handle(e *Record) {
	b := NewBuffer(e.Logs().Detail())
	defer b.Free()

	b.AppendBytes('[').AppendString(e.Level.String()).AppendBytes(']')

	var indent byte = ' '
	if e.AppendCreated != nil {
		b.AppendBytes(' ').AppendFunc(e.AppendCreated)
		indent = '\t'
	}

	if e.AppendLocation != nil {
		b.AppendBytes(' ').AppendFunc(e.AppendLocation)
		indent = '\t'
	}

	b.AppendBytes(indent).AppendFunc(e.AppendMessage)

	b.AppendBytes(h.attrs...)

	h.buildAttrs(b, e.Attrs)

	b.AppendBytes('\n')

	h.mux.Lock()
	defer h.mux.Unlock()
	if _, err := h.w.Write(b.Bytes()); err != nil { // 一次性写入，性能更好一些。
		fmt.Fprintf(os.Stderr, "NewTextHandler.Handle:%v\n", err)
	}
}

func (h *textHandler) WithAttrs(attrs []Attr) Handler {
	b := NewBuffer(false)
	defer b.Free()

	h.buildAttrs(b, attrs)

	data := make([]byte, len(h.attrs), b.Len()+len(h.attrs))
	copy(data, h.attrs)

	return &textHandler{
		w:     h.w,
		attrs: append(data, b.Bytes()...),
	}
}

func (h *textHandler) buildAttrs(b *Buffer, attrs []Attr) {
	for _, p := range attrs {
		b.AppendBytes(' ').AppendString(p.K).AppendBytes('=')
		switch v := p.V.(type) {
		case string:
			b.AppendString(v)
		case int:
			b.AppendInt(int64(v), 10)
		case int64:
			b.AppendInt(v, 10)
		case int32:
			b.AppendInt(int64(v), 10)
		case int16:
			b.AppendInt(int64(v), 10)
		case int8:
			b.AppendInt(int64(v), 10)
		case uint:
			b.AppendUint(uint64(v), 10)
		case uint64:
			b.AppendUint(v, 10)
		case uint32:
			b.AppendUint(uint64(v), 10)
		case uint16:
			b.AppendUint(uint64(v), 10)
		case uint8:
			b.AppendUint(uint64(v), 10)
		case float32:
			b.AppendFloat(float64(v), 'f', 5, 32)
		case float64:
			b.AppendFloat(v, 'f', 5, 64)
		case encoding.TextMarshaler:
			if bs, err := v.MarshalText(); err != nil {
				b.AppendString("Err(").AppendString(err.Error()).AppendBytes(')')
			} else {
				b.AppendBytes(bs...)
			}
		default:
			b.Append(p.V)
		}
	}
}

// NewJSONHandler 返回将 [Record] 以 JSON 的形式写入 w 的对象
//
// NOTE: 如果向 w 输出内容时出错，会将错误信息输出到终端作为最后的处理方式。
func NewJSONHandler(w ...io.Writer) Handler {
	return &jsonHandler{w: writers.New(w...)}
}

func (h *jsonHandler) Handle(e *Record) {
	b := NewBuffer(e.Logs().Detail())
	defer b.Free()

	b.AppendBytes('{')

	b.AppendString(`"level":"`).AppendString(e.Level.String()).AppendString(`",`).
		AppendString(`"message":"`).AppendFunc(e.AppendMessage).AppendBytes('"')

	if e.AppendCreated != nil {
		b.AppendString(`,"created":"`).AppendFunc(e.AppendCreated).AppendBytes('"')
	}

	if e.AppendLocation != nil {
		b.AppendString(`,"path":"`).AppendFunc(e.AppendLocation).AppendBytes('"')
	}

	if len(e.Attrs) > 0 {
		b.AppendString(`,"attrs":[`)

		b.AppendBytes(h.attrs...)

		h.buildAttr(b, e.Attrs)

		b.AppendBytes(']')
	}

	b.AppendBytes('}')

	h.mux.Lock()
	defer h.mux.Unlock()
	if _, err := h.w.Write(b.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "NewJSONHandler.Handle:%v\n", err)
	}
}

func (h *jsonHandler) WithAttrs(attrs []Attr) Handler {
	b := NewBuffer(false)
	h.buildAttr(b, attrs)
	data := make([]byte, len(h.attrs), b.Len()+len(h.attrs)+1)
	copy(data, h.attrs)

	if len(data) > 0 {
		data = append(data, ',')
	}

	return &jsonHandler{
		w:     h.w,
		attrs: append(data, b.Bytes()...),
	}
}

func (h *jsonHandler) buildAttr(b *Buffer, attrs []Attr) {
	for i, p := range attrs {
		if i > 0 {
			b.AppendBytes(',')
		}
		b.AppendString(`{"`).AppendString(p.K).AppendString(`":`)

		switch v := p.V.(type) {
		case string:
			b.AppendBytes('"').AppendString(v).AppendBytes('"')
		case int:
			b.AppendInt(int64(v), 10)
		case int64:
			b.AppendInt(v, 10)
		case int32:
			b.AppendInt(int64(v), 10)
		case int16:
			b.AppendInt(int64(v), 10)
		case int8:
			b.AppendInt(int64(v), 10)
		case uint:
			b.AppendUint(uint64(v), 10)
		case uint64:
			b.AppendUint(v, 10)
		case uint32:
			b.AppendUint(uint64(v), 10)
		case uint16:
			b.AppendUint(uint64(v), 10)
		case uint8:
			b.AppendUint(uint64(v), 10)
		case float32:
			b.AppendFloat(float64(v), 'f', 5, 32)
		case float64:
			b.AppendFloat(v, 'f', 5, 64)
		default:
			val, err := json.Marshal(p.V)
			if err != nil {
				val = []byte("\"Err(" + err.Error() + ")\"")
			}
			b.AppendBytes(val...)
		}

		b.AppendBytes('}')
	}
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

	return &termHandler{w: w, foreColors: cs}
}

func (h *termHandler) Handle(e *Record) {
	b := NewBuffer(e.Logs().Detail())
	defer b.Free()

	ww := colors.New(b)
	fc := h.foreColors[e.Level]
	ww.WByte('[').Color(colors.Normal, fc, colors.Default).WString(e.Level.String()).Reset().WByte(']') // [WARN]

	var indent byte = ' '
	if e.AppendCreated != nil {
		b := NewBuffer(e.Logs().Detail())
		defer b.Free()
		e.AppendCreated(b)
		ww.WByte(' ').WBytes(b.data)
		indent = '\t'
	}

	if e.AppendLocation != nil {
		b := NewBuffer(e.Logs().Detail())
		defer b.Free()
		e.AppendLocation(b)
		ww.WByte(' ').WBytes(b.data)
		indent = '\t'
	}

	bb := NewBuffer(e.Logs().Detail())
	defer bb.Free()
	e.AppendMessage(bb)
	ww.WByte(indent).WBytes(bb.data)

	for _, p := range e.Attrs {
		ww.WByte(' ').WString(p.K).WByte('=').WString(fmt.Sprint(p.V))
	}

	ww.WByte('\n')

	h.mux.Lock()
	defer h.mux.Unlock()
	if _, err := h.w.Write(b.Bytes()); err != nil {
		// 大概率是写入终端失败，直接 panic。
		panic(fmt.Sprintf("NewTermHandler.Handle:%v\n", err))
	}
}

func (h *termHandler) WithAttrs(attrs []Attr) Handler {
	as := make([]Attr, len(h.attrs), len(h.attrs)+len(attrs))
	copy(as, h.attrs)

	return &termHandler{
		w:     h.w,
		attrs: append(as, attrs...),
	}
}

// NewDispatchHandler 根据 [Level] 派发到不同的 [Handler] 对象
func NewDispatchHandler(d map[Level]Handler) Handler {
	return &dispatchHandler{handlers: d}
}

func (h *dispatchHandler) Handle(e *Record) { h.handlers[e.Level].Handle(e) }

func (h *dispatchHandler) WithAttrs(attrs []Attr) Handler {
	m := make(map[Level]Handler, len(h.handlers))
	for l, hh := range h.handlers {
		m[l] = hh.WithAttrs(attrs)
	}
	return NewDispatchHandler(m)
}

// MergeHandler 将多个 Handler 合并成一个 Handler 接口对象
func MergeHandler(w ...Handler) Handler { return &mergeHandler{handlers: w} }

func (h *mergeHandler) Handle(e *Record) {
	for _, hh := range h.handlers {
		hh.Handle(e)
	}
}

func (h *mergeHandler) WithAttrs(attrs []Attr) Handler {
	slices := make([]Handler, 0, len(h.handlers))
	for _, hh := range h.handlers {
		slices = append(slices, hh.WithAttrs(attrs))
	}
	return MergeHandler(slices...)
}

// NewNopHandler 空的 Handler 接口实现
func NewNopHandler() Handler { return nop }

func (h *nopHandler) Handle(_ *Record) {}

func (h *nopHandler) WithAttrs([]Attr) Handler { return h }
