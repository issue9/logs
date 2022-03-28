// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var emptyLoggerInst = &emptyLogger{}

var entryPool = &sync.Pool{New: func() interface{} { return &Entry{} }}

type (
	Logger interface {
		// Value 为日志提供额外的参数
		Value(name string, val any) Logger

		Print(v ...interface{})
		Printf(format string, v ...interface{})
	}

	// Entry 每一条日志产生的数据
	Entry struct {
		Pairs   []Pair // 额外的数据保存在此，比如由 Logger.Value 添加的数据。
		Level   Level
		Created time.Time // 日志的生成时间，如果 IsZero 为 true，表示禁用该功能；
		Message string

		// 以下表示日志的定位信息，如果为空表示未启用定位信息。
		Path string
		Line int
	}

	Pair struct {
		K string
		V any
	}

	logger struct {
		lv   Level
		logs *Logs
		e    *Entry
	}

	emptyLogger struct{}
)

func NewEntry() *Entry {
	e := entryPool.Get().(*Entry)
	e.Reset()
	return e
}

func (e *Entry) Reset() {
	if e.Pairs != nil {
		e.Pairs = e.Pairs[:0]
	}
	e.Path = ""
	e.Line = 0
	e.Message = ""
	e.Created = time.Now()
	e.Level = 0
}

// Destroy 回收 Entry
//
// 非必须的操作，如果是经由 NewEntry 手动申请的 Entry，可以由此方法释放，在一定程序可能会增加性能。
func (e *Entry) Destroy() { entryPool.Put(e) }

// Location 记录位置信息到 Entry
//
// 会同时写入 e.Path 和 e.Line 两个值。
//
// depth 表示调用，1 表示调用 Location 的位置；
func (e *Entry) Location(depth int) { _, e.Path, e.Line, _ = runtime.Caller(depth) }

func newLogger(l *Logs, lv Level) *logger {
	e := NewEntry()
	e.Level = lv
	return &logger{logs: l, e: e, lv: lv}
}

// Write 实现 io.Writer 供 logs.StdLogger 使用
func (e *logger) Write(data []byte) (int, error) {
	e.e.Message = string(data)
	e.e.Location(4)
	e.logs.Output(e.e)
	return len(data), nil
}

func (e *logger) Value(name string, val interface{}) Logger {
	e.e.Pairs = append(e.e.Pairs, Pair{K: name, V: val})
	return e
}

func (e *logger) Print(v ...any) {
	if len(v) > 0 {
		e.e.Message = fmt.Sprint(v...)
	}
	e.e.Location(2)
	e.logs.Output(e.e)

	e.e.Reset() // 重置 e，可以复用该对象
	e.e.Level = e.lv
}

func (e *logger) Printf(format string, v ...interface{}) {
	e.e.Message = fmt.Sprintf(format, v...)
	e.e.Location(2)
	e.logs.Output(e.e)

	e.e.Reset()
	e.e.Level = e.lv
}

func (l *emptyLogger) Value(_ string, _ interface{}) Logger { return l }

func (l *emptyLogger) Print(_ ...interface{}) {}

func (l *emptyLogger) Printf(_ string, _ ...interface{}) {}

func (l *emptyLogger) Write(bs []byte) (int, error) { return len(bs), nil }
