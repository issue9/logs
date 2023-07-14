// SPDX-License-Identifier: MIT

// Package logs 日志系统
//
// # 格式
//
// 提供了 [Handler] 接口用于处理输出的日志格式，用户可以自己实现，
// 系统也提供了几种常用的供用户选择。
//
// 同时还提供了 [Printer] 接口用于处理由 [Logger.Print] 等方法输入的数据。
// [Printer] 一般用于对用户输入的数据进行二次处理，比如进行本地化翻译等。
//
// # Logger
//
// [Logger] 为实际的日志输出接口，提供多种 [Logger] 接口的实现。
//   - [Logs.ERROR] 等为普通的日志对象；
//   - [Logs.With] 返回的是带固定参数的日志对象；
package logs

import "sync"

type Logs struct {
	mu      sync.Mutex
	handler Handler
	loggers map[Level]*logger

	caller, created bool // 是否需要生成调用位置信息和日志生成时间
}

type Option func(*Logs)

// New 声明 Logs 对象
//
// h 如果为 nil，则表示采用 [NewNopHandler]。
func New(h Handler, o ...Option) *Logs {
	if h == nil {
		h = NewNopHandler()
	}
	l := &Logs{handler: h}

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

func (logs *Logs) INFO() Logger { return logs.Logger(LevelInfo) }

func (logs *Logs) Info(v ...any) { logs.level(LevelInfo).print(3, v...) }

func (logs *Logs) Infof(format string, v ...any) {
	logs.level(LevelInfo).printf(3, format, v...)
}

func (logs *Logs) DEBUG() Logger { return logs.Logger(LevelDebug) }

func (logs *Logs) Debug(v ...any) { logs.level(LevelDebug).print(3, v...) }

func (logs *Logs) Debugf(format string, v ...any) {
	logs.level(LevelDebug).printf(3, format, v...)
}

func (logs *Logs) TRACE() Logger { return logs.Logger(LevelTrace) }

func (logs *Logs) Trace(v ...any) { logs.level(LevelTrace).print(3, v...) }

func (logs *Logs) Tracef(format string, v ...any) {
	logs.level(LevelTrace).printf(3, format, v...)
}

func (logs *Logs) WARN() Logger { return logs.Logger(LevelWarn) }

func (logs *Logs) Warn(v ...any) { logs.level(LevelWarn).print(3, v...) }

func (logs *Logs) Warnf(format string, v ...any) {
	logs.level(LevelWarn).printf(3, format, v...)
}

func (logs *Logs) ERROR() Logger { return logs.Logger(LevelError) }

func (logs *Logs) Error(v ...any) { logs.level(LevelError).print(3, v...) }

func (logs *Logs) Errorf(format string, v ...any) {
	logs.level(LevelError).printf(3, format, v...)
}

func (logs *Logs) FATAL() Logger { return logs.Logger(LevelFatal) }

func (logs *Logs) Fatal(v ...any) { logs.level(LevelFatal).print(3, v...) }

func (logs *Logs) Fatalf(format string, v ...any) {
	logs.level(LevelFatal).printf(3, format, v...)
}

// Logger 返回指定级别的日志接口
func (logs *Logs) Logger(lv Level) Logger { return logs.level(lv) }

func (logs *Logs) level(lv Level) *logger {
	if logs.handler == nop {
		return logs.loggers[levelDisable]
	}
	return logs.loggers[lv]
}

func (logs *Logs) SetHandler(w Handler) { logs.handler = w }

func (logs *Logs) output(e *Record) {
	logs.mu.Lock()
	defer logs.mu.Unlock()

	logs.handler.Handle(e)

	if len(e.Params) < poolMaxParams {
		recordPool.Put(e)
	}
}

// Caller 是否显示记录的定位信息
func Caller(l *Logs) { l.caller = true }

// Created 是否显示记录的创建时间
func Created(l *Logs) { l.created = true }

// HasCaller 是否包含定位信息
func (logs *Logs) HasCaller() bool { return logs.caller }

// HasCreated 是否包含时间信息
func (logs *Logs) HasCreated() bool { return logs.created }

// SetCaller 是否显示位置信息
func (logs *Logs) SetCaller(v bool) { logs.caller = v }

// SetCreated 是否显示时间信息
func (logs *Logs) SetCreated(v bool) { logs.created = v }
