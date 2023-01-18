// SPDX-License-Identifier: MIT

package logs

import (
	"log"
	"runtime"
	"sync"
	"time"
)

const poolMaxParams = 100

var emptyLoggerInst = &emptyLogger{}

var entryPool = &sync.Pool{New: func() interface{} { return &Entry{} }}

type (
	// Logger 日志输出接口
	Logger interface {
		// With 为日志提供额外的参数
		With(name string, val interface{}) Logger

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

	// Entry 每一条日志产生的数据
	Entry struct {
		logs *Logs

		Level   Level
		Created time.Time // 日志的生成时间
		Message string

		// 以下表示日志的定位信息
		Path string
		Line int

		// 额外的数据保存在此，比如由 Logger.With 添加的数据。
		Params []Pair
	}

	Pair struct {
		K string
		V interface{}
	}

	logger struct {
		lv     Level
		logs   *Logs
		enable bool
		std    *log.Logger
	}

	emptyLogger struct{}
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

// Location 记录位置信息到 Entry
//
// 会同时写入 e.Path 和 e.Line 两个值。
//
// depth 表示调用，1 表示调用 Location 的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么此调用将不产生任何内容。
func (e *Entry) Location(depth int) *Entry {
	if e.Logs().HasCaller() {
		_, e.Path, e.Line, _ = runtime.Caller(depth)
	}
	return e
}

func (e *Entry) With(name string, val interface{}) Logger {
	e.Params = append(e.Params, Pair{K: name, V: val})
	return e
}

func (e *Entry) Error(err error) { e.err(3, err) }

func (e *Entry) err(depth int, err error) {
	if err != nil {
		e.Message = e.logs.printer.Error(err)
	}
	e.Location(depth)
	e.logs.Output(e)
}

func (e *Entry) String(s string) { e.str(3, s) }

func (e *Entry) str(depth int, s string) {
	e.Message = e.logs.printer.String(s)
	e.Location(depth)
	e.logs.Output(e)
}

func (e *Entry) Print(v ...interface{}) { e.print(3, v...) }

func (e *Entry) print(depth int, v ...interface{}) {
	if len(v) > 0 {
		e.Message = e.logs.printer.Print(v...)
	}
	e.Location(depth)
	e.logs.Output(e)
}

func (e *Entry) Printf(format string, v ...interface{}) { e.printf(3, format, v...) }

func (e *Entry) printf(depth int, format string, v ...interface{}) {
	e.Message = e.logs.printer.Printf(format, v...)
	e.Location(depth)
	e.logs.Output(e)
}

// Write 实现 io.Writer 供 logs.StdLogger 使用
func (l *logger) Write(data []byte) (int, error) {
	if l.enable {
		ee := l.logs.NewEntry(l.lv)
		ee.Message = string(data)
		ee.Location(4)
		l.logs.Output(ee)
	}
	return len(data), nil
}

func (l *logger) stdLogger() *log.Logger {
	if l.std == nil {
		l.std = log.New(l, "", 0)
	}
	return l.std
}

func (l *logger) With(name string, val interface{}) Logger {
	if l.enable {
		return l.logs.NewEntry(l.lv).With(name, val)
	}
	return emptyLoggerInst
}

func (l *logger) Error(err error) {
	if l.enable {
		l.logs.NewEntry(l.lv).err(3, err)
	}
}

func (l *logger) String(s string) {
	if l.enable {
		l.logs.NewEntry(l.lv).str(4, s)
	}
}

func (l *logger) Print(v ...interface{}) { l.print(4, v...) }

func (l *logger) Printf(format string, v ...interface{}) { l.printf(4, format, v...) }

func (l *logger) print(depth int, v ...interface{}) {
	if l.enable {
		l.logs.NewEntry(l.lv).print(depth, v...)
	}
}

func (l *logger) printf(depth int, format string, v ...interface{}) {
	if l.enable {
		l.logs.NewEntry(l.lv).printf(depth, format, v...)
	}
}

func (l *emptyLogger) With(_ string, _ interface{}) Logger { return l }

func (l *emptyLogger) Error(_ error) {}

func (l *emptyLogger) String(_ string) {}

func (l *emptyLogger) Print(_ ...interface{}) {}

func (l *emptyLogger) Printf(_ string, _ ...interface{}) {}
