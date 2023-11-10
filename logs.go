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
	loggers map[Level]*Logger

	levels        []Level
	attrs         map[string]any
	location      bool
	detail        bool
	createdFormat string
	printer       *localeutil.Printer
}

// Marshaler 定义了序列化日志属性的方法
//
// [Recorder.With] 的 val 如果实现了该接口，
// 那么在传递进去之后会调用该接口转换成字符串之后保存。
type Marshaler interface {
	MarshalLog() string
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
		panic("参数 h 不能为空")
	}

	l := &Logs{
		levels:  AllLevels(),
		attrs:   make(map[string]any, 10),
		loggers: make(map[Level]*Logger, len(levelStrings)),
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
	}

	return l
}

// Enable 允许的日志通道
//
// 调用此函数之后，所有不在 level 参数的通道都将被关闭。
func (logs *Logs) Enable(level ...Level) {
	// TODO(go1.21): 采用 slices.Clone 代替
	ls := make([]Level, 0, len(level))
	logs.levels = append(ls, level...)
}

// AppendAttrs 为所有的 [Logger] 对象添加属性
func (logs *Logs) AppendAttrs(attrs map[string]any) {
	for _, l := range logs.loggers {
		l.AppendAttrs(attrs)
	}
}

// IsEnable 指定级别日志是否会真实被启用
//
// 如果设置了 [Handler] 为空值或是未在 [Logs.Enable] 中指定都将返回 false。
func (logs *Logs) IsEnable(l Level) bool {
	// TODO(go1.21): 采用 slices.Index 代替
	return sliceutil.Exists(logs.levels, func(v Level, _ int) bool { return v == l })
}

func (logs *Logs) INFO() *Logger { return logs.Logger(LevelInfo) }

func (logs *Logs) DEBUG() *Logger { return logs.Logger(LevelDebug) }

func (logs *Logs) TRACE() *Logger { return logs.Logger(LevelTrace) }

func (logs *Logs) WARN() *Logger { return logs.Logger(LevelWarn) }

func (logs *Logs) ERROR() *Logger { return logs.Logger(LevelError) }

func (logs *Logs) FATAL() *Logger { return logs.Logger(LevelFatal) }

// Logger 返回指定级别的日志对象
func (logs *Logs) Logger(lv Level) *Logger { return logs.loggers[lv] }
