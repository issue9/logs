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

func newEntry(l *Logs) *entry { return &entry{logs: l, e: &Entry{}} }

func (e *Entry) location(depth int) {
	if e.Path == "" {
		_, e.Path, e.Line, _ = runtime.Caller(depth)
	}
}

func (e *Entry) print(depth int, v ...interface{}) {
	if len(v) > 0 {
		e.Message = fmt.Sprint(v...)
	}
	e.location(depth)
	e.Created = time.Now()
}

func (e *Entry) printf(depth int, format string, v ...interface{}) {
	e.Message = fmt.Sprintf(format, v...)
	e.Created = time.Now()
	e.location(depth)
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
