// SPDX-License-Identifier: MIT

// Package logs 日志系统
//
// # 格式
//
// 提供了 [Handler] 接口用于处理输出的日志格式，用户可以自己实现，
// 系统也提供了几种常用的供用户选择。
//
// # Logger
//
// [Logger] 为实际的日志输出接口，提供多种 [Logger] 的实现。
//   - [Logs.ERROR] 等为普通的日志对象；
//   - [Logs.With] 返回的是带固定参数的日志对象；
package logs

import (
	"github.com/issue9/localeutil"
	"github.com/issue9/sliceutil"
)

type Logs struct {
	handler Handler
	loggers map[Level]Logger
	attrs   map[string]any

	location, detail bool
	createdFormat    string
	printer          *localeutil.Printer
}

func attrs2Pairs(p *localeutil.Printer, attrs map[string]any) []Pair {
	pairs := make([]Pair, 0, len(attrs))

	if p == nil {
		for k, v := range attrs {
			pairs = append(pairs, Pair{K: k, V: v})
		}
	} else {
		for k, v := range attrs {
			if ls, ok := v.(localeutil.Stringer); ok {
				v = ls.LocaleString(p)
			}
			pairs = append(pairs, Pair{K: k, V: v})
		}
	}

	return pairs
}

// New 声明 Logs 对象
//
// h 如果为 nil，则表示采用 [NewNopHandler]。
func New(h Handler, o ...Option) *Logs {
	if h == nil {
		h = NewNopHandler()
	}

	l := &Logs{
		handler: h,
		loggers: make(map[Level]Logger, len(levelStrings)),
	}
	for _, opt := range o {
		opt(l)
	}

	for lv := range levelStrings {
		l.loggers[lv] = &logger{
			logs:  l,
			lv:    lv,
			pairs: attrs2Pairs(l.printer, l.attrs),
		}
	}

	return l
}

// Enable 允许的日志通道
//
// 调用此函数之后，所有不在 level 参数的通道都将被关闭。
func (logs *Logs) Enable(level ...Level) {
	for lv, l := range logs.loggers {
		if !sliceutil.Exists(level, func(ll Level, _ int) bool { return ll == lv }) { // 不存在于启用列表
			logs.loggers[lv] = disabledLogger
			continue
		}

		if l != disabledLogger {
			continue
		}

		logs.loggers[lv] = &logger{
			logs:  logs,
			lv:    lv,
			pairs: attrs2Pairs(logs.printer, logs.attrs),
		}
	}
}

func (logs *Logs) IsEnable(l Level) bool { return logs.loggers[l] != disabledLogger }

func (logs *Logs) INFO() Logger { return logs.Logger(LevelInfo) }

func (logs *Logs) DEBUG() Logger { return logs.Logger(LevelDebug) }

func (logs *Logs) TRACE() Logger { return logs.Logger(LevelTrace) }

func (logs *Logs) WARN() Logger { return logs.Logger(LevelWarn) }

func (logs *Logs) ERROR() Logger { return logs.Logger(LevelError) }

func (logs *Logs) FATAL() Logger { return logs.Logger(LevelFatal) }

// Logger 返回指定级别的日志接口
func (logs *Logs) Logger(lv Level) Logger { return logs.level(lv) }

func (logs *Logs) level(lv Level) Logger {
	if logs.handler == nop {
		return disabledLogger
	}
	return logs.loggers[lv]
}

func (logs *Logs) SetHandler(h Handler) { logs.handler = h }
