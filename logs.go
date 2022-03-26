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
	LevelInfo Level = iota + 1
	LevelTrace
	LevelDebug
	LevelWarn
	LevelError
	LevelFatal
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

func (l Level) MarshalText() ([]byte, error) { return []byte(l.String()), nil }

type Logs struct {
	mux    sync.Mutex // 防止多个 logger 对象引用同一个 writer 造成混合输入的情况
	levels map[Level]*logger
}

func New() *Logs {
	logs := &Logs{}

	l := make(map[Level]*logger, len(levels))
	for lv := range levels {
		l[lv] = &logger{
			enable: true,
			level:  lv,
			logs:   logs,
		}
	}
	logs.levels = l

	return logs
}

// Enable 允许的日志通道
//
// 默认情况下所有的通道都是允许的。
// 调用此函数之后，所有不在 level 参数的通道都将被关闭。
func (logs *Logs) Enable(level ...Level) {
	for lv, logger := range logs.levels {
		logger.enable = sliceutil.Exists(level, func(l Level) bool { return l == lv })
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

// Output 输出 Entry 对象
//
// 相对于其它方法，该方法比较自由，可以由 e 决定最终输出到哪儿，内容也由用户定义。
func (logs *Logs) Output(e *Entry) {
	logs.mux.Lock()
	defer logs.mux.Unlock()
	logs.levels[e.Level].w.WriteEntry(e)
}
