// SPDX-License-Identifier: MIT

// Package logs 日志系统
package logs

import (
	"sync"

	"github.com/issue9/sliceutil"
)

type Level int8

// 目前支持的日志类型
const (
	LevelInfo Level = iota
	LevelTrace
	LevelDebug
	LevelWarn
	LevelError
	LevelFatal
	levelSize
)

var levels = map[Level]string{
	LevelInfo:  "INFO",
	LevelTrace: "TRAC",
	LevelDebug: "DBUG",
	LevelWarn:  "WARN",
	LevelError: "ERRO",
	LevelFatal: "FATL",
}

func (l Level) String() string { return levels[l] }

type Logs struct {
	mux    sync.Mutex // 防止多个 logger 对象引用同一个 writer 造成混合输入的情况
	levels map[Level]*logger
}

func New() *Logs {
	l := make(map[Level]*logger, levelSize)
	for lv := range levels {
		l[lv] = &logger{}
	}

	return &Logs{levels: l}
}

// EnableLevels 允许的 Level
func (logs *Logs) EnableLevels(level ...Level) {
	for lv, logger := range logs.levels {
		logger.enable = sliceutil.Exists(level, func(e Level) bool { return e == lv })
	}
}

func (logs *Logs) IsEnable(l Level) bool { return logs.levels[l].enable }

func (logs *Logs) SetOutput(w Writer, level ...Level) {
	for _, lv := range level {
		logs.levels[lv].w = w
	}
}

func (logs *Logs) INFO() Logger {
	l := logs.levels[LevelInfo]
	if l.enable {
		return l
	}
	return emptyLoggerInst
}

func (logs *Logs) Info(v ...any) { logs.INFO().Print(v...) }

func (logs *Logs) Infof(format string, v ...any) { logs.INFO().Printf(format, v...) }

func (logs *Logs) DEBUG() Logger {
	l := logs.levels[LevelDebug]
	if l.enable {
		return l
	}
	return emptyLoggerInst
}

func (logs *Logs) Debug(v ...any) { logs.DEBUG().Print(v...) }

func (logs *Logs) Debugf(format string, v ...any) { logs.DEBUG().Printf(format, v...) }

func (logs *Logs) TRACE() Logger {
	l := logs.levels[LevelTrace]
	if l.enable {
		return l
	}
	return emptyLoggerInst
}

func (logs *Logs) Trace(v ...any) { logs.TRACE().Print(v...) }

func (logs *Logs) Tracef(format string, v ...any) { logs.TRACE().Printf(format, v...) }

func (logs *Logs) WARN() Logger {
	l := logs.levels[LevelWarn]
	if l.enable {
		return l
	}
	return emptyLoggerInst
}

func (logs *Logs) Warn(v ...any) { logs.WARN().Print(v...) }

func (logs *Logs) Warnf(format string, v ...any) { logs.WARN().Printf(format, v...) }

func (logs *Logs) ERROR() Logger {
	l := logs.levels[LevelError]
	if l.enable {
		return l
	}
	return emptyLoggerInst
}

func (logs *Logs) Error(v ...any) { logs.ERROR().Print(v...) }

func (logs *Logs) Errorf(format string, v ...any) { logs.ERROR().Printf(format, v...) }

func (logs *Logs) FATAL() Logger {
	l := logs.levels[LevelFatal]
	if l.enable {
		return l
	}
	return emptyLoggerInst
}

func (logs *Logs) Fatal(v ...any) { logs.FATAL().Print(v...) }

func (logs *Logs) Fatalf(format string, v ...any) { logs.FATAL().Printf(format, v...) }

func (logs *Logs) output(e *Entry) {
	logs.mux.Lock()
	defer logs.mux.Unlock()

	logs.levels[e.Level].w.WriteEntry(e)
}
