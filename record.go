// SPDX-License-Identifier: MIT

package logs

import (
	"io"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/issue9/localeutil"
	"golang.org/x/xerrors"

	"github.com/issue9/logs/v6/writers"
)

const poolMaxParams = 100

var recordPool = &sync.Pool{New: func() any { return &Record{} }}

type (
	// Record 单条日志产生的数据
	Record struct {
		logs *Logs

		Level Level

		// AppendCreated 添加字符串类型的日志创建时间
		//
		// 可能为空，根据 [Logs.CreatedFormat] 是否为空决定。
		AppendCreated AppendFunc

		// AppendMessage 向日志中添加字符串类型的日志消息
		//
		// 这是每一条日志的主消息，不会为空。
		// 内容是根据 Depth* 系列方法生成的。
		AppendMessage AppendFunc

		// AppendLocation 添加字符串类型的日志触发位置信息
		//
		// 可能为空，根据 [Logs.HasLocation] 决定。
		AppendLocation AppendFunc

		// 额外的数据，比如由 [Logger.With] 添加的数据。
		Params []Pair
	}

	Pair struct {
		K string
		V any
	}
)

func (logs *Logs) NewRecord(lv Level) *Record {
	e := recordPool.Get().(*Record)

	e.logs = logs
	if e.Params != nil {
		e.Params = e.Params[:0]
	}
	e.AppendLocation = nil
	e.AppendMessage = nil
	e.AppendCreated = nil
	e.Level = lv

	return e
}

// 转换成 io.Writer
//
// 仅供内部使用，因为 depth 值的关系，只有固定的调用层级关系才能正常显示行号。
func (e *Record) asWriter() io.Writer {
	return writers.WriteFunc(func(data []byte) (int, error) {
		e.DepthString(5, string(data))
		return len(data), nil
	})
}

func (e *Record) Logs() *Logs { return e.logs }

// depth 表示调用，1 表示调用 Location 的位置；
func (e *Record) initLocationCreated(depth int) *Record {
	if e.Logs().HasLocation() {
		_, p, l, _ := runtime.Caller(depth)
		e.AppendLocation = func(b *Buffer) {
			b.AppendString(p).AppendBytes(':').AppendInt(int64(l), 10)
		}
	}

	if e.Logs().createdFormat != "" {
		t := time.Now() // 必须是当前时间，而不是放在 AppendCreated 中获取的时间。
		e.AppendCreated = func(b *Buffer) { b.AppendTime(t, e.Logs().createdFormat) }
	}

	return e
}

func (e *Record) With(name string, val any) Logger {
	if ls, ok := val.(localeutil.Stringer); ok {
		val = ls.LocaleString(e.Logs().printer)
	}
	e.Params = append(e.Params, Pair{K: name, V: val})
	return e
}

func (e *Record) StdLogger() *log.Logger {
	return log.New(e.asWriter(), "", 0)
}

func (e *Record) Error(err error) { e.DepthError(2, err) }

// DepthError 输出 error 类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthError(depth int, err error) {
	if err == nil {
		panic("参数 err 不能为空")
	}

	switch ee := err.(type) {
	case xerrors.Formatter:
		e.AppendMessage = func(b *Buffer) { appendError(e.Logs().printer, b, ee) }
	case localeutil.Stringer:
		if pp := e.Logs().printer; pp != nil {
			e.AppendMessage = func(b *Buffer) { b.AppendString(ee.LocaleString(pp)) }
		} else { // e2 必然是实现了 error 接口的
			e.AppendMessage = func(b *Buffer) { b.AppendString(ee.(error).Error()) }
		}
	default:
		e.AppendMessage = func(b *Buffer) { b.AppendString(err.Error()) }
	}

	e.initLocationCreated(depth + 1).output()
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

func (e *Record) String(s string) { e.DepthString(2, s) }

// DepthString 输出字符串类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthString(depth int, s string) {
	e.AppendMessage = func(b *Buffer) { b.AppendString(s) }
	e.initLocationCreated(depth + 1).output()
}

func (e *Record) Print(v ...any) { e.DepthPrint(2, v...) }

// DepthPrint 输出任意类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthPrint(depth int, v ...any) {
	replaceLocaleString(e.Logs().printer, v)
	e.AppendMessage = func(b *Buffer) { b.Append(v...) }
	e.initLocationCreated(depth + 1).output()
}

func (e *Record) Printf(format string, v ...any) { e.DepthPrintf(2, format, v...) }

// DepthPrintf 输出任意类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthPrintf(depth int, format string, v ...any) {
	replaceLocaleString(e.Logs().printer, v)
	e.AppendMessage = func(b *Buffer) { b.Appendf(format, v...) }
	e.initLocationCreated(depth + 1).output()
}

func (e *Record) Println(v ...any) { e.DepthPrintln(2, v...) }

// DepthPrintln 输出任意类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasLocation] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthPrintln(depth int, v ...any) {
	replaceLocaleString(e.Logs().printer, v)
	e.AppendMessage = func(b *Buffer) { b.Appendln(v...) }
	e.initLocationCreated(depth + 1).output()
}

func (e *Record) output() {
	e.Logs().handler.Handle(e)
	if len(e.Params) < poolMaxParams {
		recordPool.Put(e)
	}
}

func replaceLocaleString(p *localeutil.Printer, v []any) {
	for i, val := range v {
		if ls, ok := val.(localeutil.Stringer); ok {
			v[i] = ls.LocaleString(p)
		}
	}
}
