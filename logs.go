// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

// Package logs 高性能的日志系统
package logs

import (
	"slices"
	"sync"

	"github.com/issue9/localeutil"
)

var attrLogsPool = &sync.Pool{
	New: func() any { return &AttrLogs{loggers: make(map[Level]*Logger, LevelFatal+1)} },
}

type Logs struct {
	loggers map[Level]*Logger

	levels        []Level
	attrs         map[string]any
	location      bool
	detail        bool
	createdFormat string
	printer       *localeutil.Printer
}

// AttrLogs 带有固定属性的日志
type AttrLogs struct {
	attrs   map[string]any
	loggers map[Level]*Logger
	logs    *Logs
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
		h = NewNopHandler()
	}

	l := &Logs{
		loggers: make(map[Level]*Logger, len(levelStrings)),

		levels: AllLevels(),
		attrs:  make(map[string]any, 10),
	}
	for _, opt := range o {
		opt(l)
	}

	for lv := range levelStrings {
		l.loggers[lv] = &Logger{
			logs: l,
			lv:   lv,
			h:    h.New(l.detail, lv, map2Slice(l.printer, l.attrs)),
		}
	}

	return l
}

// Enable 允许的日志通道
//
// 调用此函数之后，所有不在 level 参数的通道都将被关闭。
func (logs *Logs) Enable(level ...Level) { logs.levels = slices.Clone(level) }

// AppendAttrs 为所有的 [Logger] 对象添加属性
func (logs *Logs) AppendAttrs(attrs map[string]any) {
	for _, l := range logs.loggers {
		l.AppendAttrs(attrs)
	}
}

// IsEnable 指定级别日志是否会真实被启用
func (logs *Logs) IsEnable(l Level) bool { return slices.Index(logs.levels, l) >= 0 }

func (logs *Logs) INFO() *Logger { return logs.Logger(LevelInfo) }

func (logs *Logs) DEBUG() *Logger { return logs.Logger(LevelDebug) }

func (logs *Logs) TRACE() *Logger { return logs.Logger(LevelTrace) }

func (logs *Logs) WARN() *Logger { return logs.Logger(LevelWarn) }

func (logs *Logs) ERROR() *Logger { return logs.Logger(LevelError) }

func (logs *Logs) FATAL() *Logger { return logs.Logger(LevelFatal) }

// Logger 返回指定级别的日志对象
func (logs *Logs) Logger(lv Level) *Logger { return logs.loggers[lv] }

// New 声明一组带有 attrs 属性的日志
func (logs *Logs) New(attrs map[string]any) *AttrLogs {
	l := attrLogsPool.Get().(*AttrLogs)
	l.attrs = attrs
	clear(l.loggers)
	l.logs = logs
	return l
}

// IsEnable 指定级别日志是否会真实被启用
func (logs *AttrLogs) IsEnable(l Level) bool { return logs.logs.IsEnable(l) }

func (logs *AttrLogs) INFO() *Logger { return logs.Logger(LevelInfo) }

func (logs *AttrLogs) DEBUG() *Logger { return logs.Logger(LevelDebug) }

func (logs *AttrLogs) TRACE() *Logger { return logs.Logger(LevelTrace) }

func (logs *AttrLogs) WARN() *Logger { return logs.Logger(LevelWarn) }

func (logs *AttrLogs) ERROR() *Logger { return logs.Logger(LevelError) }

func (logs *AttrLogs) FATAL() *Logger { return logs.Logger(LevelFatal) }

// Logger 返回指定级别的日志对象
func (logs *AttrLogs) Logger(lv Level) *Logger {
	if _, found := logs.loggers[lv]; !found {
		logs.loggers[lv] = logs.logs.Logger(lv).New(logs.attrs)
	}
	return logs.loggers[lv]
}

// AppendAttrs 为所有的 [Logger] 对象添加属性
func (l *AttrLogs) AppendAttrs(attrs map[string]any) {
	for _, ll := range l.loggers {
		if ll != nil {
			ll.AppendAttrs(attrs)
		}
	}
	for k, v := range attrs {
		l.attrs[k] = v
	}
}

func (l *AttrLogs) NewRecord() *Record { return l.logs.NewRecord() }

// FreeAttrLogs 回收 [AttrLogs]
//
// 如果需要频繁地生成 [AttrLogs] 且其生命周期都有固定的销毁时间点，
// 可以用此方法达到复用 [AttrLogs] 以达到些许性能提升。
//
// NOTE: 此操作会让 logs 不再可用。
func FreeAttrLogs(logs *AttrLogs) {
	for _, l := range logs.loggers {
		if l != nil {
			l.free()
		}
	}

	attrLogsPool.Put(logs)
}
