// SPDX-FileCopyrightText: 2014-2024 caixw
//
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
	// Recorder 定义了输出一条日志记录的各种方法
	Recorder interface {
		// With 创建带有指定属性的 [Recorder] 对象
		//
		// 返回对象与当前对象未必是同一个，由实现者决定。
		// 且返回对象是一次性的，在调用 Error、String 等输出之后即被回收。
		//
		// 如果 val 实现了 [localeutil.Stringer] 或是 [Marshaler] 接口，
		// 将被转换成字符串保存。
		With(name string, val any) Recorder

		// Error 将一条错误信息作为一条日志输出
		//
		// 这是 Print 的特化版本，在已知类型为 error 时，
		// 采用此方法会比 Print(err) 有更好的性能。
		//
		// 如果 err 实现了 [xerrors.FormatError] 接口，同时也会打印调用信息。
		//
		// NOTE: 此操作之后，当前对象不再可用！
		Error(err error)

		// LocaleString 输出一条本地化的信息
		//
		// 这是 Print 的特化版本，在已知类型为 [localeutil.Stringer] 时，
		// 采用此方法会比 Print(s) 有更好的性能。
		//
		// NOTE: 此操作之后，当前对象不再可用！
		LocaleString(localeutil.Stringer)

		// String 将字符串作为一条日志输出
		//
		// 这是 Print 的特化版本，在已知类型为字符串时，
		// 采用此方法会比 Print(s) 有更好的性能。
		//
		// NOTE: 此操作之后，当前对象不再可用！
		String(s string)

		// 输出一条日志信息
		//
		// NOTE: 此操作之后，当前对象不再可用！
		Print(v ...any)

		// 输出一条日志信息
		//
		// NOTE: 此操作之后，当前对象不再可用！
		//
		// NOTE: 不会对内容进行翻译，否则对于字符串类型的枚举类型可能会输出意想不到的内容，
		// 如果需要翻译内容，可以调用 [Recorder.LocaleString]。
		Printf(format string, v ...any)
	}

	// Record 单条日志输出时产生的数据
	//
	// NOTE: 该对象只能由 [Logs.NewRecord] 生成。
	Record struct {
		logs *Logs

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
		l *Logger
		r *Record
	}

	disableRecorder struct{}
)

func (logs *Logs) NewRecord() *Record {
	e := recordPool.Get().(*Record)

	if e.Attrs != nil {
		e.Attrs = e.Attrs[:0]
	}
	e.AppendLocation = nil
	e.AppendMessage = nil
	e.AppendCreated = nil
	e.logs = logs

	return e
}

// depth 表示调用，1 表示调用此方法的位置；
func (e *Record) initLocationCreated(depth int) *Record {
	if e.logs.HasLocation() {
		_, p, l, _ := runtime.Caller(depth)
		e.AppendLocation = func(b *Buffer) {
			b.AppendString(p).AppendBytes(':').AppendInt(int64(l), 10)
		}
	}

	if e.logs.createdFormat != "" {
		t := time.Now() // 必须是当前时间，而不是放在 AppendCreated 中获取的时间。
		e.AppendCreated = func(b *Buffer) { b.AppendTime(t, e.logs.createdFormat) }
	}

	return e
}

func (e *Record) with(name string, val any) *Record {
	switch v := val.(type) {
	case localeutil.Stringer:
		e.Attrs = append(e.Attrs, Attr{K: name, V: v.LocaleString(e.logs.printer)})
	case Marshaler:
		e.Attrs = append(e.Attrs, Attr{K: name, V: v.MarshalLog()})
	default:
		e.Attrs = append(e.Attrs, Attr{K: name, V: v})
	}
	return e
}

// DepthError 输出 error 类型的内容到日志
//
// depth 表示调用，2 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthError(depth int, err error) *Record {
	if err == nil {
		panic("参数 err 不能为空")
	}

	switch ee := err.(type) {
	case xerrors.Formatter:
		e.AppendMessage = func(b *Buffer) { appendError(e.logs.printer, b, ee) }
	case localeutil.Stringer:
		if pp := e.logs.printer; pp != nil {
			e.AppendMessage = func(b *Buffer) { b.AppendString(ee.LocaleString(pp)) }
		} else { // e2 必然是实现了 error 接口的
			e.AppendMessage = func(b *Buffer) { b.AppendString(err.Error()) }
		}
	default:
		e.AppendMessage = func(b *Buffer) { b.AppendString(err.Error()) }
	}

	return e.initLocationCreated(depth)
}

// DepthLocaleString 输出 [localeutil.Stringer]  类型的内容到日志
//
// depth 表示调用，2 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthLocaleString(depth int, s localeutil.Stringer) *Record {
	e.AppendMessage = func(b *Buffer) { b.AppendString(s.LocaleString(e.logs.printer)) }
	return e.initLocationCreated(depth)
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

// DepthString 输出字符串类型的内容到日志
//
// depth 表示调用，2 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthString(depth int, s string) *Record {
	e.AppendMessage = func(b *Buffer) { b.AppendString(s) }
	return e.initLocationCreated(depth)
}

// DepthPrint 输出任意类型的内容到日志
//
// depth 表示调用，2 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthPrint(depth int, v ...any) *Record {
	replaceLocaleString(e.logs.printer, v)
	e.AppendMessage = func(b *Buffer) { b.Append(v...) }
	return e.initLocationCreated(depth)
}

// DepthPrintf 输出任意类型的内容到日志
//
// depth 表示调用，2 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
//
// NOTE: 不会对内容进行翻译，否则对于字符串类型的枚举类型可能会输出意想不到的内容，
// 如果需要翻译内容，可以调用 [Record.DepthLocaleString]。
func (e *Record) DepthPrintf(depth int, format string, v ...any) *Record {
	replaceLocaleString(e.logs.printer, v)
	e.AppendMessage = func(b *Buffer) { b.Appendf(format, v...) }
	return e.initLocationCreated(depth)
}

// Output 输出当前记录到日志
func (e *Record) Output(l *Logger) {
	const poolMaxAttrs = 100
	l.Handler().Handle(e)
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
	e.r.with(name, val)
	return e
}

func (e *withRecorder) Error(err error) {
	e.r.DepthError(3, err).Output(e.l)
	e.free()
}

func (e *withRecorder) String(s string) {
	e.r.DepthString(3, s).Output(e.l)
	e.free()
}

func (e *withRecorder) LocaleString(s localeutil.Stringer) {
	e.r.DepthLocaleString(3, s).Output(e.l)
	e.free()
}

func (e *withRecorder) Print(v ...any) {
	e.r.DepthPrint(3, v...).Output(e.l)
	e.free()
}

func (e *withRecorder) Printf(format string, v ...any) {
	e.r.DepthPrintf(3, format, v...).Output(e.l)
	e.free()
}

func (e *withRecorder) free() { withRecordPool.Put(e) }

func (l *disableRecorder) With(string, any) Recorder { return l }

func (l *disableRecorder) Error(error) {}

func (l *disableRecorder) String(string) {}

func (l *disableRecorder) LocaleString(localeutil.Stringer) {}

func (l *disableRecorder) Print(...any) {}

func (l *disableRecorder) Printf(string, ...any) {}
