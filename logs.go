// SPDX-License-Identifier: MIT

// Package logs 日志系统
//
// # 格式
//
// 提供了 [Handler] 接口用于处理输出的日志格式，用户可以自己实现，
// 系统也提供了几种常用的供用户选择。
package logs

import (
	"github.com/issue9/localeutil"
	"github.com/issue9/sliceutil"
)

type Logs struct {
	handler Handler
	loggers map[Level]*Logger
	enables map[Level]bool

	attrs            map[string]any // 仅用于被 Option 函数存取，没有其它用处。
	location, detail bool
	createdFormat    string
	printer          *localeutil.Printer
}

func map2Slice(p *localeutil.Printer, attrs map[string]any) []Attr {
	pairs := make([]Attr, 0, len(attrs))

	if p == nil {
		for k, v := range attrs {
			pairs = append(pairs, Attr{K: k, V: v})
		}
	} else {
		for k, v := range attrs {
			if ls, ok := v.(localeutil.Stringer); ok {
				v = ls.LocaleString(p)
			}
			pairs = append(pairs, Attr{K: k, V: v})
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
		loggers: make(map[Level]*Logger, len(levelStrings)),
		enables: make(map[Level]bool, len(levelStrings)),

		attrs: make(map[string]any, 10),
	}
	for _, opt := range o {
		opt(l)
	}

	attrs := map2Slice(l.printer, l.attrs)
	for lv := range levelStrings {
		l.loggers[lv] = &Logger{
			logs: l,
			lv:   lv,
			h:    h.New(l.detail, lv, attrs),
		}
		l.enables[lv] = true
	}

	return l
}

// Enable 允许的日志通道
//
// 调用此函数之后，所有不在 level 参数的通道都将被关闭。
func (logs *Logs) Enable(level ...Level) {
	for lv := range logs.enables {
		logs.enables[lv] = sliceutil.Exists(level, func(ll Level, _ int) bool { return ll == lv })
	}
}

// IsEnable 指定级别日志是否会真实被启用
//
// 如果设置了 [Handler] 为空值或是未在 [Logs.Enable] 中指定都将返回 false。
func (logs *Logs) IsEnable(l Level) bool { return logs.enables[l] && logs.handler != nop }

func (logs *Logs) INFO() *Logger { return logs.Logger(LevelInfo) }

func (logs *Logs) DEBUG() *Logger { return logs.Logger(LevelDebug) }

func (logs *Logs) TRACE() *Logger { return logs.Logger(LevelTrace) }

func (logs *Logs) WARN() *Logger { return logs.Logger(LevelWarn) }

func (logs *Logs) ERROR() *Logger { return logs.Logger(LevelError) }

func (logs *Logs) FATAL() *Logger { return logs.Logger(LevelFatal) }

// Logger 返回指定级别的日志对象
func (logs *Logs) Logger(lv Level) *Logger { return logs.loggers[lv] }
