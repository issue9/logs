// SPDX-License-Identifier: MIT

// Package logs 日志系统
//
// # 格式
//
// 提供了 [Writer] 接口用于处理输出的日志格式，用户可以自己实现，
// 系统也提供了几种常用的供用户选择。
//
// 同时还提供了 [Printer] 接口用于处理 [Logger.Print] 等方法输入的数据。
// [Printer] 一般用于对用户输入的数据进行二次处理，比如进行本地化翻译等。
//
// # Logger
//
// [Logger] 为实际的日志输出接口，提供多种 [Logger] 接口的实现。
//
// - [Logs.ERROR] 等为普通的日志对象；
// - [Logs.With] 返回的是带固定参数的日志对象；
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

	printer Printer
}

// New 声明 Logs 对象
//
// w 如果为 nil，则表示采用 [NewNopWriter]。
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
	if l.printer == nil {
		DefaultPrint(l)
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

func (logs *Logs) Info(v ...interface{}) { logs.level(LevelInfo).print(4, v...) }

func (logs *Logs) Infof(format string, v ...interface{}) {
	logs.level(LevelInfo).printf(4, format, v...)
}

func (logs *Logs) DEBUG() Logger { return logs.Logger(LevelDebug) }

func (logs *Logs) Debug(v ...interface{}) { logs.level(LevelDebug).print(4, v...) }

func (logs *Logs) Debugf(format string, v ...interface{}) {
	logs.level(LevelDebug).printf(4, format, v...)
}

func (logs *Logs) TRACE() Logger { return logs.Logger(LevelTrace) }

func (logs *Logs) Trace(v ...interface{}) { logs.level(LevelTrace).print(4, v...) }

func (logs *Logs) Tracef(format string, v ...interface{}) {
	logs.level(LevelTrace).printf(4, format, v...)
}

func (logs *Logs) WARN() Logger { return logs.Logger(LevelWarn) }

func (logs *Logs) Warn(v ...interface{}) { logs.level(LevelWarn).print(4, v...) }

func (logs *Logs) Warnf(format string, v ...interface{}) {
	logs.level(LevelWarn).printf(4, format, v...)
}

func (logs *Logs) ERROR() Logger { return logs.Logger(LevelError) }

func (logs *Logs) Error(v ...interface{}) { logs.level(LevelError).print(4, v...) }

func (logs *Logs) Errorf(format string, v ...interface{}) {
	logs.level(LevelError).printf(4, format, v...)
}

func (logs *Logs) FATAL() Logger { return logs.Logger(LevelFatal) }

func (logs *Logs) Fatal(v ...interface{}) { logs.level(LevelFatal).print(4, v...) }

func (logs *Logs) Fatalf(format string, v ...interface{}) {
	logs.level(LevelFatal).printf(4, format, v...)
}

// Logger 返回指定级别的日志接口
func (logs *Logs) Logger(lv Level) Logger { return logs.level(lv) }

func (logs *Logs) level(lv Level) *logger {
	if logs.w == nop {
		return logs.loggers[levelDisable]
	}
	return logs.loggers[lv]
}

func (logs *Logs) SetOutput(w Writer) { logs.w = w }

// Output 输出 [Entry] 对象
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
// NOTE: 不要设置 [log.Logger] 的 Prefix 和 flag，这些配置项 logs 本身有提供。
// [log.Logger] 应该仅作为输出 [Entry.Message] 内容使用。
func (logs *Logs) StdLogger(l Level) *log.Logger { return logs.level(l).stdLogger() }
