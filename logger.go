// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const poolMaxParams = 100

var emptyLoggerInst = &emptyLogger{}

var entryPool = &sync.Pool{New: func() interface{} { return &Entry{} }}

type (
	Logger interface {
		// Value 为日志提供额外的参数
		Value(name string, val interface{}) Logger

		// Error 输出 error 接口数据
		//
		// 这是 Print 的特化版本，在输出 error 时性能上会比直接用 Print(err) 好。
		Error(err error)

		Print(v ...interface{})
		Printf(format string, v ...interface{})
	}

	// Entry 每一条日志产生的数据
	Entry struct {
		logs *Logs

		Level   Level     `json:"level"`
		Created time.Time `json:"created,omitempty"` // 日志的生成时间，如果 IsZero 为 true，表示禁用该功能；
		Message string    `json:"message"`

		// 以下表示日志的定位信息，如果为空表示未启用定位信息。
		Path string `json:"path,omitempty"`
		Line int    `json:"line,omitempty"`

		// 额外的数据保存在此，比如由 Logger.Value 添加的数据。
		Params []Pair `json:"params,omitempty"`
	}

	Pair struct {
		K string
		V interface{}
	}

	logger struct {
		lv     Level
		logs   *Logs
		enable bool
	}

	emptyLogger struct{}
)

func (logs *Logs) NewEntry(lv Level) *Entry {
	e := entryPool.Get().(*Entry)
	e.Reset(logs, lv)
	return e
}

func (e *Entry) Reset(l *Logs, lv Level) {
	e.logs = l

	if e.Params != nil {
		e.Params = e.Params[:0]
	}
	e.Path = ""
	e.Line = 0
	e.Message = ""
	if l.HasCreated() {
		e.Created = time.Now()
	} else {
		e.Created = time.Time{}
	}
	e.Level = lv
}

func (e *Entry) Logs() *Logs { return e.logs }

// Location 记录位置信息到 Entry
//
// 会同时写入 e.Path 和 e.Line 两个值。
//
// depth 表示调用，1 表示调用 Location 的位置；
//
// 如果 Logs.HasCaller 为 false，那么此调用将不产生任何内容。
func (e *Entry) Location(depth int) {
	if e.Logs().HasCaller() {
		_, e.Path, e.Line, _ = runtime.Caller(depth)
	}
}

func (e *Entry) Value(name string, val interface{}) Logger {
	e.Params = append(e.Params, Pair{K: name, V: val})
	return e
}

func (e *Entry) Error(err error) { e.err(3, err) }

func (e *Entry) err(depth int, err error) {
	if err != nil {
		e.Message = err.Error()
	}
	e.Location(depth)
	e.logs.Output(e)
}

func (e *Entry) Print(v ...interface{}) { e.print(3, v...) }

func (e *Entry) print(depth int, v ...interface{}) {
	if len(v) > 0 {
		e.Message = fmt.Sprint(v...)
	}
	e.Location(depth)
	e.logs.Output(e)
}

func (e *Entry) Printf(format string, v ...interface{}) { e.printf(3, format, v...) }

func (e *Entry) printf(depth int, format string, v ...interface{}) {
	e.Message = fmt.Sprintf(format, v...)
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

func (l *logger) Value(name string, val interface{}) Logger {
	if l.enable {
		return l.logs.NewEntry(l.lv).Value(name, val)
	}
	return emptyLoggerInst
}

func (l *logger) Error(err error) {
	if l.enable {
		l.logs.NewEntry(l.lv).err(3, err)
	}
}

func (l *logger) Print(v ...interface{}) {
	if l.enable {
		l.logs.NewEntry(l.lv).print(3, v...)
	}
}

func (l *logger) Printf(format string, v ...interface{}) {
	if l.enable {
		l.logs.NewEntry(l.lv).printf(3, format, v...)
	}
}

func (l *emptyLogger) Value(_ string, _ interface{}) Logger { return l }

func (l *emptyLogger) Error(_ error) {}

func (l *emptyLogger) Print(_ ...interface{}) {}

func (l *emptyLogger) Printf(_ string, _ ...interface{}) {}
