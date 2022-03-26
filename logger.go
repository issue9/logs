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

	Entry struct {
		Pairs   []Pair
		Level   Level
		Created time.Time
		Message string
		Path    string
		Line    int
	}

	Pair struct {
		K string
		V any
	}

	entry struct {
		logs *Logs
		e    *Entry
	}

	logger struct {
		logs   *Logs
		level  Level
		enable bool
		w      Writer
	}

	emptyLogger struct{}
)

func NewEntry() *Entry {
	ee := entryPool.Get().(*Entry)
	if ee.Pairs != nil {
		ee.Pairs = ee.Pairs[:0]
	}
	ee.Path = ""
	ee.Line = 0
	ee.Message = ""
	ee.Created = time.Now()
	ee.Level = 0

	return ee
}

func newEntry(l *Logs) *entry { return &entry{logs: l, e: NewEntry()} }

// Location 记录位置信息到 Entry
//
// 会同时写入 e.Path 和 e.Line 两个值。
//
// depth 表示调用，1 表示调用 Location 的位置；
func (e *Entry) Location(depth int) {
	_, e.Path, e.Line, _ = runtime.Caller(depth)
}

func (e *Entry) print(depth int, v ...interface{}) {
	if len(v) > 0 {
		e.Message = fmt.Sprint(v...)
	}
	e.Location(depth)
}

func (e *Entry) printf(depth int, format string, v ...interface{}) {
	e.Message = fmt.Sprintf(format, v...)
	e.Location(depth)
}

func (e *entry) setLevel(l Level) *entry {
	e.e.Level = l
	return e
}

func (e *entry) Value(name string, val interface{}) Logger {
	e.e.Pairs = append(e.e.Pairs, Pair{K: name, V: val})
	return e
}

func (e *entry) Print(v ...any) {
	e.e.print(3, v...)
	e.logs.Output(e.e)
}

func (e *entry) Printf(format string, v ...interface{}) {
	e.e.printf(3, format, v...)
	e.logs.Output(e.e)
}

func (l *logger) Value(name string, val interface{}) Logger {
	return newEntry(l.logs).setLevel(l.level).Value(name, val)
}

func (l *logger) Print(v ...interface{}) {
	newEntry(l.logs).setLevel(l.level).Print(v...)
}

func (l *logger) Printf(format string, v ...interface{}) {
	newEntry(l.logs).setLevel(l.level).Printf(format, v...)
}

func (l *emptyLogger) Value(_ string, _ interface{}) Logger { return l }

func (l *emptyLogger) Print(_ ...interface{}) {}

func (l *emptyLogger) Printf(_ string, _ ...interface{}) {}
