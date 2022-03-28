// SPDX-License-Identifier: MIT

// Package logs 日志系统
package logs

import (
	"log"
	"sync"

	"github.com/issue9/sliceutil"
)

type Logs struct {
	mu     sync.Mutex
	w      Writer
	levels map[Level]bool
}

func New(w Writer) *Logs {
	if w == nil {
		w = NewNopWriter()
	}

	levels := make(map[Level]bool, 6)
	for lv := range levelStrings {
		levels[lv] = true
	}

	return &Logs{
		w:      w,
		levels: levels,
	}
}

// Enable 允许的日志通道
//
// 调用此函数之后，所有不在 level 参数的通道都将被关闭。
func (logs *Logs) Enable(level ...Level) {
	for lv := range logs.levels {
		logs.levels[lv] = sliceutil.Exists(level, func(l Level) bool { return l == lv })
	}
}

func (logs *Logs) IsEnable(l Level) bool { return logs.levels[l] }

func (logs *Logs) SetOutput(w Writer) { logs.w = w }

func (logs *Logs) INFO() Logger { return logs.level(LevelInfo) }

func (logs *Logs) Info(v ...any) { logs.INFO().Print(v...) }

func (logs *Logs) Infof(format string, v ...any) { logs.INFO().Printf(format, v...) }

func (logs *Logs) DEBUG() Logger { return logs.level(LevelDebug) }

func (logs *Logs) Debug(v ...any) { logs.DEBUG().Print(v...) }

func (logs *Logs) Debugf(format string, v ...any) { logs.DEBUG().Printf(format, v...) }

func (logs *Logs) TRACE() Logger { return logs.level(LevelTrace) }

func (logs *Logs) Trace(v ...any) { logs.TRACE().Print(v...) }

func (logs *Logs) Tracef(format string, v ...any) { logs.TRACE().Printf(format, v...) }

func (logs *Logs) WARN() Logger { return logs.level(LevelWarn) }

func (logs *Logs) Warn(v ...any) { logs.WARN().Print(v...) }

func (logs *Logs) Warnf(format string, v ...any) { logs.WARN().Printf(format, v...) }

func (logs *Logs) ERROR() Logger { return logs.level(LevelError) }

func (logs *Logs) Error(v ...any) { logs.ERROR().Print(v...) }

func (logs *Logs) Errorf(format string, v ...any) { logs.ERROR().Printf(format, v...) }

func (logs *Logs) FATAL() Logger { return logs.level(LevelFatal) }

func (logs *Logs) Fatal(v ...any) { logs.FATAL().Print(v...) }

func (logs *Logs) Fatalf(format string, v ...any) { logs.FATAL().Printf(format, v...) }

func (logs *Logs) level(lv Level) Logger {
	if logs.IsEnable(lv) {
		return newEntry(logs, lv)
	}
	return emptyLoggerInst
}

// Output 输出 Entry 对象
//
// 相对于其它方法，该方法比较自由，可以由 e 决定最终输出到哪儿，内容也由用户定义。
func (logs *Logs) Output(e *Entry) {
	logs.mu.Lock()
	defer logs.mu.Unlock()
	logs.w.WriteEntry(e)
}

// StdLogger 转换成标准库的 Logger
func (logs *Logs) StdLogger(l Level) *log.Logger {
	if logs.IsEnable(l) {
		return log.New(newEntry(logs, l), "", 0)
	}
	return log.New(emptyLoggerInst, "", 0)
}
