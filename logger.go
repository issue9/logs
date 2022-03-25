// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"runtime"
	"time"
)

var emptyLoggerInst = &emptyLogger{}

type (
	Logger interface {
		// Value 为日志提供额外的参数
		Value(name string, val any) Logger

		// Location 将当前行作为日志输出的定位信息
		//
		// 如果未调用该函数，那么在 Print 和 Printf 中会隐藏调用。
		Location() Logger

		Print(v ...interface{})
		Printf(format string, v ...interface{})
	}

	Entry struct {
		logs *Logs

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

	logger struct {
		logs   *Logs
		level  Level
		enable bool
		w      Writer
	}

	emptyLogger struct{}
)

func newEntry(l *Logs) *Entry { return &Entry{logs: l} }

func (e *Entry) setLevel(l Level) *Entry {
	e.Level = l
	return e
}

func (e *Entry) Value(name string, val any) Logger {
	e.Pairs = append(e.Pairs, Pair{K: name, V: val})
	return e
}

func (e *Entry) Location() Logger { return e.location(2) }

func (e *Entry) location(depth int) Logger {
	_, e.Path, e.Line, _ = runtime.Caller(depth)
	return e
}

func (e *Entry) Print(v ...any) {
	if len(v) > 0 {
		e.Message = fmt.Sprint(v...)
	}
	e.location(2)
	e.logs.output(e)
}

func (e *Entry) Printf(format string, v ...any) {
	e.Message = fmt.Sprintf(format, v...)
	e.logs.output(e)
}

func (l *logger) Value(name string, val any) Logger {
	return newEntry(l.logs).setLevel(l.level).Value(name, val)
}

func (l *logger) Location() Logger {
	return newEntry(l.logs).setLevel(l.level).location(2)
}

func (l *logger) Print(v ...any) {
	newEntry(l.logs).setLevel(l.level).Print(v...)
}

func (l *logger) Printf(format string, v ...any) {
	newEntry(l.logs).setLevel(l.level).Printf(format, v...)
}

func (l *emptyLogger) Value(_ string, _ any) Logger { return l }

func (l *emptyLogger) Location() Logger { return l }

func (l *emptyLogger) Print(_ ...any) {}

func (l *emptyLogger) Printf(_ string, _ ...any) {}
