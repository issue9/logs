// SPDX-License-Identifier: MIT

package logs

import (
	"runtime"
	"sync"
	"time"
)

const poolMaxParams = 100

var emptyInputInst = &emptyInput{}

var entryPool = &sync.Pool{New: func() interface{} { return &Entry{} }}

type (
	// Input 日志输入提供的接口
	Input interface {
		// With 为日志提供额外的参数
		//
		// 返回值是当前对象。
		With(name string, val interface{}) Input

		// Error 将一条错误信息作为一条日志输出
		//
		// 这是 Print 的特化版本，在已知类型为 error 时，
		// 采用此方法会比 Print(err) 有更好的性能。
		Error(err error)

		// String 将字符串作为一条日志输出
		//
		// 这是 Print 的特化版本，在已知类型为字符串时，
		// 采用此方法会比 Print(s) 有更好的性能。
		String(s string)

		// 输出一条日志信息
		Print(v ...interface{})
		Printf(format string, v ...interface{})
	}

	// Entry 每一条日志的表示对象
	Entry struct {
		logs *Logs

		Level   Level
		Created time.Time // 日志的生成时间

		// 日志的实际内容
		//
		// 如果要改变此值，请使用 Depth* 系列的方法。
		Message string

		// 以下表示日志的定位信息
		Path string
		Line int

		// 额外的数据保存在此，比如由 [Logger.With] 添加的数据。
		Params []Pair
	}

	Pair struct {
		K string
		V interface{}
	}

	emptyInput struct{}
)

func (logs *Logs) NewEntry(lv Level) *Entry {
	e := entryPool.Get().(*Entry)

	e.logs = logs
	if e.Params != nil {
		e.Params = e.Params[:0]
	}
	e.Path = ""
	e.Line = 0
	e.Message = ""
	if logs.HasCreated() {
		e.Created = time.Now()
	} else {
		e.Created = time.Time{}
	}
	e.Level = lv

	return e
}

func (e *Entry) Logs() *Logs { return e.logs }

// Location 记录位置信息
//
// 会同时写入 e.Path 和 e.Line 两个值。
//
// depth 表示调用，1 表示调用 Location 的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么此调用将不产生任何内容。
//
// Deprecated: 可以直接使用 [Entry.DepthError] 等方法代替
func (e *Entry) Location(depth int) *Entry {
	return e.setLocation(depth + 1)
}

// depth 表示调用，1 表示调用 Location 的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么此调用将不产生任何内容。
func (e *Entry) setLocation(depth int) *Entry {
	if e.Logs().HasCaller() {
		_, e.Path, e.Line, _ = runtime.Caller(depth)
	}
	return e
}

func (e *Entry) With(name string, val interface{}) Input {
	e.Params = append(e.Params, Pair{K: name, V: val})
	return e
}

func (e *Entry) Error(err error) { e.DepthError(2, err) }

// DepthError 输出 error 类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么 depth 将不起实际作用。
func (e *Entry) DepthError(depth int, err error) {
	if err != nil {
		e.Message = e.logs.printer.Error(err)
	}
	e.setLocation(depth + 1)
	e.logs.output(e)
}

func (e *Entry) String(s string) { e.DepthString(2, s) }

// DepthString 输出字符串类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么 depth 将不起实际作用。
func (e *Entry) DepthString(depth int, s string) {
	e.Message = e.logs.printer.String(s)
	e.setLocation(depth + 1)
	e.logs.output(e)
}

func (e *Entry) Print(v ...interface{}) { e.DepthPrint(2, v...) }

// DepthPrint 输出任意类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么 depth 将不起实际作用。
func (e *Entry) DepthPrint(depth int, v ...interface{}) {
	if len(v) > 0 {
		e.Message = e.logs.printer.Print(v...)
	}
	e.setLocation(depth + 1)
	e.logs.output(e)
}

func (e *Entry) Printf(format string, v ...interface{}) { e.DepthPrintf(2, format, v...) }

// DepthPrintf 输出任意类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么 depth 将不起实际作用。
func (e *Entry) DepthPrintf(depth int, format string, v ...interface{}) {
	e.Message = e.logs.printer.Printf(format, v...)
	e.setLocation(depth + 1)
	e.logs.output(e)
}

func (l *emptyInput) With(_ string, _ interface{}) Input { return l }

func (l *emptyInput) Error(_ error) {}

func (l *emptyInput) String(_ string) {}

func (l *emptyInput) Print(_ ...interface{}) {}

func (l *emptyInput) Printf(_ string, _ ...interface{}) {}