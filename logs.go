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

type Logs struct {
	handler Handler
	loggers map[Level]*logger

	caller        bool // 是否需要生成调用位置信息
	createdFormat string
}

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

func (logs *Logs) DEBUG() Logger { return logs.Logger(LevelDebug) }

func (logs *Logs) TRACE() Logger { return logs.Logger(LevelTrace) }

func (logs *Logs) WARN() Logger { return logs.Logger(LevelWarn) }

func (logs *Logs) ERROR() Logger { return logs.Logger(LevelError) }

func (logs *Logs) FATAL() Logger { return logs.Logger(LevelFatal) }

// Logger 返回指定级别的日志接口
func (logs *Logs) Logger(lv Level) Logger { return logs.level(lv) }

func (logs *Logs) level(lv Level) *logger {
	if logs.handler == nop {
		return logs.loggers[levelDisable]
	}
	return logs.loggers[lv]
}

func (logs *Logs) SetHandler(h Handler) { logs.handler = h }
