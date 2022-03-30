// SPDX-License-Identifier: MIT

// Package logs 日志系统
package logs

import (
	"log"
	"sync"
)

type Logs struct {
	mu      sync.Mutex
	w       Writer
	loggers map[Level]*logger

	// 是否需要生成调用位置信息和日志生成时间
	caller, created bool
}

type Option func(*Logs)

// Caller 是否显示记录的定位信息
func Caller(l *Logs) { l.caller = true }

// Created 是否显示记录的创建时间
func Created(l *Logs) { l.created = true }

// New 声明 Logs 对象
//
// w 如果为 nil，则表示采用 NewNopWriter。
func New(w Writer, o ...Option) *Logs {
	if w == nil {
		w = NewNopWriter()
	}
	l := &Logs{w: w}

	l.loggers = make(map[Level]*logger, len(levelStrings))
	for lv := range levelStrings {
		l.loggers[lv] = &logger{
			logs:   l,
			lv:     lv,
			enable: lv != levelDisable,
		}
	}

	for _, opt := range o {
		opt(l)
	}
	return l
}

func (logs *Logs) SetOutput(w Writer) { logs.w = w }

// Enable 允许的日志通道
//
// 调用此函数之后，所有不在 level 参数的通道都将被关闭。
func (logs *Logs) Enable(level ...Level) {
	exists := func(lv Level) bool {
		if lv == levelDisable {
			return false
		}

		for _, l := range level {
			if l == lv {
				return true
			}
		}
		return false
	}

	for _, l := range logs.loggers {
		l.enable = exists(l.lv)
	}
}

func (logs *Logs) IsEnable(l Level) bool { return logs.loggers[l].enable }

func (logs *Logs) INFO() Logger { return logs.level(LevelInfo) }

func (logs *Logs) Info(v ...interface{}) { logs.INFO().Print(v...) }

func (logs *Logs) Infof(format string, v ...interface{}) { logs.INFO().Printf(format, v...) }

func (logs *Logs) DEBUG() Logger { return logs.level(LevelDebug) }

func (logs *Logs) Debug(v ...interface{}) { logs.DEBUG().Print(v...) }

func (logs *Logs) Debugf(format string, v ...interface{}) { logs.DEBUG().Printf(format, v...) }

func (logs *Logs) TRACE() Logger { return logs.level(LevelTrace) }

func (logs *Logs) Trace(v ...interface{}) { logs.TRACE().Print(v...) }

func (logs *Logs) Tracef(format string, v ...interface{}) { logs.TRACE().Printf(format, v...) }

func (logs *Logs) WARN() Logger { return logs.level(LevelWarn) }

func (logs *Logs) Warn(v ...interface{}) { logs.WARN().Print(v...) }

func (logs *Logs) Warnf(format string, v ...interface{}) { logs.WARN().Printf(format, v...) }

func (logs *Logs) ERROR() Logger { return logs.level(LevelError) }

func (logs *Logs) Error(v ...interface{}) { logs.ERROR().Print(v...) }

func (logs *Logs) Errorf(format string, v ...interface{}) { logs.ERROR().Printf(format, v...) }

func (logs *Logs) FATAL() Logger { return logs.level(LevelFatal) }

func (logs *Logs) Fatal(v ...interface{}) { logs.FATAL().Print(v...) }

func (logs *Logs) Fatalf(format string, v ...interface{}) { logs.FATAL().Printf(format, v...) }

func (logs *Logs) level(lv Level) *logger {
	if logs.w == nop {
		return logs.loggers[levelDisable]
	}
	return logs.loggers[lv]
}

// Output 输出 Entry 对象
//
// 相对于其它方法，该方法比较自由，可以由 e 决定最终输出到哪儿，内容也由用户定义。
func (logs *Logs) Output(e *Entry) {
	logs.mu.Lock()
	defer logs.mu.Unlock()

	logs.w.WriteEntry(e)

	if len(e.Params) < poolMaxParams {
		entryPool.Put(e)
	}
}

// StdLogger 转换成标准库的 Logger
//
// NOTE: 不要设置 log.Logger 的 Prefix 和 flag，这些配置项 logs 本身有提供。
// log.Logger 应该仅作为输出 Entry.Message 内容使用。
func (logs *Logs) StdLogger(l Level) *log.Logger { return log.New(logs.level(l), "", 0) }

// HasCaller 是否包含定位信息
func (logs *Logs) HasCaller() bool { return logs.caller }

// HasCreated 是否包含时间信息
func (logs *Logs) HasCreated() bool { return logs.created }
