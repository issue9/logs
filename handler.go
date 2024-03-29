// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package logs

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"sync"

	"github.com/issue9/term/v3/colors"

	"github.com/issue9/logs/v7/writers"
)

var defaultTermColors = map[Level]colors.Color{
	LevelInfo:  colors.Green,
	LevelDebug: colors.Yellow,
	LevelTrace: colors.Yellow,
	LevelWarn:  colors.Yellow,
	LevelError: colors.Red,
	LevelFatal: colors.Red,
}

var nop = &nopHandler{}

type (
	// Handler 日志后端的处理接口
	Handler interface {
		// Handle 将 [Record] 写入日志
		//
		// [Record] 中各个字段的名称由处理器自行决定；
		// detail 表示是否显示错误的堆栈信息；
		//
		// NOTE: 此方法应该保证输出内容是以换行符作为结尾。
		Handle(r *Record)

		// New 根据当前对象的参数派生新的 [Handler] 对象
		//
		// detail 表示是否需要显示错误的调用堆栈信息；
		// lv 表示输出的日志级别；
		// attrs 表示日志属性；
		// 这三个参数主要供 [Handler] 缓存这些数据以提升性能；
		//
		// 对于重名的问题并无规定，只要 [Handler] 自身能处理相应的情况即可。
		//
		// NOTE: 即便所有的参数均为零值，也应该返回一个新的对象。
		New(detail bool, lv Level, attrs []Attr) Handler
	}

	textHandler struct {
		w   io.Writer
		mux sync.Mutex

		attrs  []byte // 预编译的属性值
		level  []byte // 预处理的 level 内容
		detail bool
	}

	jsonHandler struct {
		w   io.Writer
		mux sync.Mutex

		attrs  []byte // 预编译的属性值
		level  []byte // 预处理的 level 内容
		detail bool
	}

	termHandler struct {
		textHandler
		foreColors map[Level]colors.Color
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
	if err := h.handle(e); err != nil {
		fmt.Fprintf(os.Stderr, "NewTextHandler.Handle:%v\n", err)
	}
}

func (h *textHandler) handle(e *Record) error {
	b := NewBuffer(h.detail)
	defer b.Free()

	b.AppendBytes(h.level...)

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
	// 必须要在 Buffer 回收之前将内容写入 h.w
	// 一次性写入，性能更好一些。
	_, err := h.w.Write(b.Bytes())
	return err
}

func (h *textHandler) New(detail bool, lv Level, attrs []Attr) Handler {
	b := NewBuffer(false)
	defer b.Free()
	h.buildAttrs(b, attrs)

	data := make([]byte, 0, b.Len()+len(h.attrs))
	data = append(data, h.attrs...)

	return &textHandler{
		w: h.w,

		attrs:  append(data, b.Bytes()...),
		level:  []byte("[" + lv.String() + "]"),
		detail: detail,
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
			b.AppendFloat(float64(v), 'f', -1, 32)
		case float64:
			b.AppendFloat(v, 'f', -1, 64)
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
	b := NewBuffer(h.detail)
	defer b.Free()

	b.AppendBytes('{')

	b.AppendBytes(h.level...)

	b.AppendString(`"message":"`).AppendFunc(e.AppendMessage).AppendBytes('"')

	if e.AppendCreated != nil {
		b.AppendString(`,"created":"`).AppendFunc(e.AppendCreated).AppendBytes('"')
	}

	if e.AppendLocation != nil {
		b.AppendString(`,"path":"`).AppendFunc(e.AppendLocation).AppendBytes('"')
	}

	if len(e.Attrs) > 0 || len(h.attrs) > 0 {
		b.AppendString(`,"attrs":[`)

		b.AppendBytes(h.attrs...)

		if len(e.Attrs) > 0 && len(h.attrs) > 0 {
			b.AppendBytes(',')
		}

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

func (h *jsonHandler) New(detail bool, lv Level, attrs []Attr) Handler {
	b := NewBuffer(false)
	defer b.Free()

	h.buildAttr(b, attrs)
	data := make([]byte, 0, b.Len()+len(h.attrs)+1)
	data = append(data, h.attrs...)
	if len(h.attrs) > 0 && len(attrs) > 0 {
		data = append(data, ',')
	}

	return &jsonHandler{
		w: h.w,

		attrs:  append(data, b.Bytes()...),
		level:  []byte(`"level":"` + lv.String() + `",`),
		detail: detail,
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
			b.AppendFloat(float64(v), 'f', -1, 32)
		case float64:
			b.AppendFloat(v, 'f', -1, 64)
		default:
			val, err := json.Marshal(p.V)
			if err != nil {
				val = []byte(`"Err(` + err.Error() + `)"`)
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
	if w == nil {
		panic("参数 w 不能为空")
	}

	cs := make(map[Level]colors.Color, len(defaultTermColors))
	for l, cc := range defaultTermColors {
		if c, found := foreColors[l]; found {
			cs[l] = c
		} else {
			cs[l] = cc
		}
	}

	return &termHandler{textHandler: textHandler{w: w}, foreColors: cs}
}

func (h *termHandler) Handle(e *Record) {
	if err := h.handle(e); err != nil {
		// 大概率是写入终端失败，直接 panic。
		panic(fmt.Sprintf("NewTermHandler.Handle:%v\n", err))
	}
}

func (h *termHandler) New(detail bool, lv Level, attrs []Attr) Handler {
	l := "[" + colors.Sprint(colors.Normal, h.foreColors[lv], colors.Default, lv.String()) + "]"

	b := NewBuffer(false)
	defer b.Free()
	h.buildAttrs(b, attrs)
	data := make([]byte, 0, b.Len()+len(h.attrs))
	data = append(data, h.attrs...)

	return &termHandler{
		textHandler: textHandler{
			w: h.w,

			attrs:  append(data, b.Bytes()...),
			level:  []byte(l),
			detail: detail,
		},
		foreColors: maps.Clone(h.foreColors),
	}
}

// NewDispatchHandler 根据 [Level] 派发到不同的 [Handler] 对象
//
// 返回对象的 [Handler.New] 方法会根据其传递的 Level 参数从 d 中选择一个相应的对象返回。
func NewDispatchHandler(d map[Level]Handler) Handler {
	if len(d) != len(levelStrings) {
		panic("NewDispatchHandler: 需指定所有 Level 对应的对象")
	}

	return &dispatchHandler{handlers: d}
}

func (h *dispatchHandler) Handle(e *Record) { panic("不支持该功能") }

func (h *dispatchHandler) New(detail bool, lv Level, attrs []Attr) Handler {
	if hh, found := h.handlers[lv]; found {
		return hh.New(detail, lv, attrs)
	}
	panic(fmt.Sprintf("无效的 lv 参数：%v", lv)) // 由 [NewDispatchHandler] 确保不会执行到此
}

// MergeHandler 将多个 [Handler] 合并成一个 [Handler] 接口对象
func MergeHandler(w ...Handler) Handler {
	handlers := make([]Handler, 0, len(w))
	for _, ww := range w {
		if h, ok := ww.(*mergeHandler); ok {
			handlers = append(handlers, h.handlers...)
		} else {
			handlers = append(handlers, ww)
		}
	}
	return &mergeHandler{handlers: handlers}
}

func (h *mergeHandler) Handle(e *Record) {
	for _, hh := range h.handlers {
		hh.Handle(e)
	}
}

func (h *mergeHandler) New(detail bool, lv Level, attrs []Attr) Handler {
	slices := make([]Handler, 0, len(h.handlers))
	for _, hh := range h.handlers {
		slices = append(slices, hh.New(detail, lv, attrs))
	}
	return MergeHandler(slices...)
}

func NewNopHandler() Handler { return nop }

func (h *nopHandler) Handle(*Record) {}

func (h *nopHandler) New(bool, Level, []Attr) Handler { return h }
