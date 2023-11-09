// SPDX-License-Identifier: MIT

package logs

import (
	"runtime"
	"sync"
	"time"

	"github.com/issue9/localeutil"
	"golang.org/x/xerrors"
)

var (
	recordPool       = &sync.Pool{New: func() any { return &Record{} }}
	withRecordPool   = &sync.Pool{New: func() any { return &withRecorder{} }}
	disabledRecorder = &disableRecorder{}
)

type (
	// Recorder 日志的输出接口
	Recorder interface {
		// With 创建带有指定属性的 [Recorder] 对象
		//
		// 返回对象与当前对象未必是同一个，由实现者决定。
		// 且返回对象是一次性的，在调用 Error、String 等输出之后即被回收。
		With(name string, val any) Recorder

		// Error 将一条错误信息作为一条日志输出
		//
		// 这是 Print 的特化版本，在已知类型为 error 时，
		// 采用此方法会比 Print(err) 有更好的性能。
		//
		// 如果 err 实现了 [xerrors.FormatError] 接口，同时也会打印调用信息。
		Error(err error)

		// String 将字符串作为一条日志输出
		//
		// 这是 Print 的特化版本，在已知类型为字符串时，
		// 采用此方法会比 Print(s) 有更好的性能。
		String(s string)

		// 输出一条日志信息
		Print(v ...any)
		Println(v ...any)
		Printf(format string, v ...any)
	}

	// Record 单条日志输出时产生的数据
	Record struct {
		// AppendCreated 添加字符串类型的日志创建时间
		//
		// 可能为空，根据 [Logs.CreatedFormat] 是否为空决定。
		AppendCreated AppendFunc

		// AppendMessage 向日志中添加字符串类型的日志消息
		//
		// 这是每一条日志的主消息，不会为空。
		AppendMessage AppendFunc

		// AppendLocation 添加字符串类型的日志触发位置信息
		//
		// 可能为空，根据 [Logs.HasLocation] 决定。
		AppendLocation AppendFunc

		// 额外的数据，比如由 [Recorder.With] 添加的数据。
		Attrs []Attr
	}

	Attr struct {
		K string
		V any
	}

	withRecorder struct {
		h    Handler
		logs *Logs
		r    *Record
	}

	disableRecorder struct{}
)

func NewRecord() *Record {
	e := recordPool.Get().(*Record)

	if e.Attrs != nil {
		e.Attrs = e.Attrs[:0]
	}
	e.AppendLocation = nil
	e.AppendMessage = nil
	e.AppendCreated = nil

	return e
}

// depth 表示调用，1 表示调用此方法的位置；
func (e *Record) initLocationCreated(logs *Logs, depth int) *Record {
	if logs.HasLocation() {
		_, p, l, _ := runtime.Caller(depth)
		e.AppendLocation = func(b *Buffer) {
			b.AppendString(p).AppendBytes(':').AppendInt(int64(l), 10)
		}
	}

	if logs.createdFormat != "" {
		t := time.Now() // 必须是当前时间，而不是放在 AppendCreated 中获取的时间。
		e.AppendCreated = func(b *Buffer) { b.AppendTime(t, logs.createdFormat) }
	}

	return e
}

func (e *Record) with(logs *Logs, name string, val any) *Record {
	if ls, ok := val.(localeutil.Stringer); ok && logs.printer != nil {
		val = ls.LocaleString(logs.printer)
	}
	e.Attrs = append(e.Attrs, Attr{K: name, V: val})
	return e
}

// depthError 输出 error 类型的内容到日志
//
// depth 表示调用，2 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) depthError(logs *Logs, h Handler, depth int, err error) {
	if err == nil {
		panic("参数 err 不能为空")
	}

	switch ee := err.(type) {
	case xerrors.Formatter:
		e.AppendMessage = func(b *Buffer) { appendError(logs.printer, b, ee) }
	case localeutil.Stringer:
		if pp := logs.printer; pp != nil {
			e.AppendMessage = func(b *Buffer) { b.AppendString(ee.LocaleString(pp)) }
		} else { // e2 必然是实现了 error 接口的
			e.AppendMessage = func(b *Buffer) { b.AppendString(ee.(error).Error()) }
		}
	default:
		e.AppendMessage = func(b *Buffer) { b.AppendString(err.Error()) }
	}

	e.initLocationCreated(logs, depth).output(logs.detail, h)
}

func appendError(p *localeutil.Printer, b *Buffer, ef xerrors.Formatter) {
	err := ef.FormatError(b)
	for err != nil {
		switch e2 := err.(type) {
		case xerrors.Formatter:
			err = e2.FormatError(b)
		case localeutil.Stringer:
			if p != nil {
				b.AppendString(e2.LocaleString(p))
			} else { // e2 必然是实现了 error 接口的
				b.AppendString(e2.(error).Error())
			}
			return
		default:
			b.AppendString(e2.Error())
			return
		}
	}
}

// depthString 输出字符串类型的内容到日志
//
// depth 表示调用，2 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) depthString(logs *Logs, h Handler, depth int, s string) {
	e.AppendMessage = func(b *Buffer) { b.AppendString(s) }
	e.initLocationCreated(logs, depth).output(logs.detail, h)
}

// depthPrint 输出任意类型的内容到日志
//
// depth 表示调用，2 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) depthPrint(logs *Logs, h Handler, depth int, v ...any) {
	replaceLocaleString(logs.printer, v)
	e.AppendMessage = func(b *Buffer) { b.Append(v...) }
	e.initLocationCreated(logs, depth).output(logs.detail, h)
}

// depthPrintf 输出任意类型的内容到日志
//
// depth 表示调用，2 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) depthPrintf(logs *Logs, h Handler, depth int, format string, v ...any) {
	replaceLocaleString(logs.printer, v)
	e.AppendMessage = func(b *Buffer) { b.Appendf(format, v...) }
	e.initLocationCreated(logs, depth).output(logs.detail, h)
}

// depthPrintln 输出任意类型的内容到日志
//
// depth 表示调用，2 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) depthPrintln(logs *Logs, h Handler, depth int, v ...any) {
	replaceLocaleString(logs.printer, v)
	e.AppendMessage = func(b *Buffer) { b.Appendln(v...) }
	e.initLocationCreated(logs, depth).output(logs.detail, h)
}

func (e *Record) output(detail bool, h Handler) {
	const poolMaxAttrs = 100
	h.Handle(e)
	if len(e.Attrs) < poolMaxAttrs {
		recordPool.Put(e)
	}
}

func replaceLocaleString(p *localeutil.Printer, v []any) {
	if p == nil {
		return
	}

	for i, val := range v {
		if ls, ok := val.(localeutil.Stringer); ok {
			v[i] = ls.LocaleString(p)
		}
	}
}

func (e *withRecorder) With(name string, val any) Recorder {
	e.r.with(e.logs, name, val)
	return e
}

func (e *withRecorder) Error(err error) {
	e.r.depthError(e.logs, e.h, 3, err)
	e.free()
}

func (e *withRecorder) String(s string) {
	e.r.depthString(e.logs, e.h, 3, s)
	e.free()
}

func (e *withRecorder) Print(v ...any) {
	e.r.depthPrint(e.logs, e.h, 3, v...)
	e.free()
}

func (e *withRecorder) Printf(format string, v ...any) {
	e.r.depthPrintf(e.logs, e.h, 3, format, v...)
	e.free()
}

func (e *withRecorder) Println(v ...any) {
	e.r.depthPrintln(e.logs, e.h, 3, v...)
	e.free()
}

func (e *withRecorder) free() { withRecordPool.Put(e) }

func (l *disableRecorder) With(string, any) Recorder { return l }

func (l *disableRecorder) Error(error) {}

func (l *disableRecorder) String(string) {}

func (l *disableRecorder) Print(...any) {}

func (l *disableRecorder) Printf(string, ...any) {}

func (l *disableRecorder) Println(...any) {}
