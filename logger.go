// SPDX-License-Identifier: MIT

package logs

import (
	"io"
	"log"
)

var disabledLogger = &disableLogger{}

type (
	// Logger 日志接口
	Logger interface {
		// With 为日志提供额外的参数
		//
		// 返回值是当前对象。
		With(name string, val any) Logger

		// Error 将一条错误信息作为一条日志输出
		//
		// 这是 Print 的特化版本，在已知类型为 error 时，
		// 采用此方法会比 Print(err) 有更好的性能。
		//
		// 如果 err 实现了 [xerrors.FormatError] 接口，同时也会打印调用信息。
		Error(err error)

		// String 将字符串作为一条日志输出
		//
		// 这是 Print 的特化版本，在已知类型为字符串时，
		// 采用此方法会比 Print(s) 有更好的性能。
		String(s string)

		// 输出一条日志信息
		Print(v ...any)
		Println(v ...any)
		Printf(format string, v ...any)

		// StdLogger 将当前对象转换成标准库的日志对象
		//
		// NOTE: 不要设置返回对象的 Prefix 和 Flag，这些配置项与当前模块的功能有重叠。
		// [log.Logger] 应该仅作为向 [Logger] 输入 [Record.Message] 内容使用。
		StdLogger() *log.Logger
	}

	logger struct {
		lv     Level
		logs   *Logs
		enable bool
	}

	withLogger struct {
		l     *logger
		pairs []Pair
	}

	disableLogger struct{}
)

func (l *logger) StdLogger() *log.Logger {
	w := io.Discard
	if l.enable {
		w = l.logs.NewRecord(l.lv).asWriter()
	}
	return log.New(w, "", 0)
}

func (l *logger) With(name string, val any) Logger {
	if l.enable {
		return l.logs.NewRecord(l.lv).With(name, val)
	}
	return disabledLogger
}

func (l *logger) Error(err error) {
	if l.enable {
		l.logs.NewRecord(l.lv).DepthError(3, err)
	}
}

func (l *logger) String(s string) {
	if l.enable {
		l.logs.NewRecord(l.lv).DepthString(2, s)
	}
}

func (l *logger) Print(v ...any) {
	if l.enable {
		l.logs.NewRecord(l.lv).DepthPrint(2, v...)
	}
}

func (l *logger) Println(v ...any) {
	if l.enable {
		l.logs.NewRecord(l.lv).DepthPrintln(2, v...)
	}
}

func (l *logger) Printf(format string, v ...any) {
	if l.enable {
		l.logs.NewRecord(l.lv).DepthPrintf(2, format, v...)
	}
}

// With 创建带有指定参数的日志对象
//
// params 自动添加的参数，每条日志都将自动带上这些参数；
func (logs *Logs) With(lv Level, params map[string]any) Logger {
	l := logs.level(lv)
	if !l.enable {
		return disabledLogger
	}

	pairs := make([]Pair, 0, len(params))
	for k, v := range params {
		pairs = append(pairs, Pair{K: k, V: v})
	}

	return &withLogger{
		l:     logs.level(lv),
		pairs: pairs,
	}
}

func (l *withLogger) with() *Record {
	e := l.l.logs.NewRecord(l.l.lv)
	for _, pair := range l.pairs {
		e.With(pair.K, pair.V)
	}
	return e
}

func (l *withLogger) StdLogger() *log.Logger {
	return log.New(l.with().asWriter(), "", 0)
}

func (l *withLogger) With(k string, v any) Logger {
	return l.with().With(k, v)
}

func (l *withLogger) Error(err error) { l.with().DepthError(2, err) }

func (l *withLogger) String(s string) { l.with().DepthString(2, s) }

func (l *withLogger) Print(v ...any) { l.with().DepthPrint(2, v...) }

func (l *withLogger) Println(v ...any) { l.with().DepthPrintln(2, v...) }

func (l *withLogger) Printf(format string, v ...any) {
	l.with().DepthPrintf(2, format, v...)
}

func (l *disableLogger) With(_ string, _ any) Logger { return l }

func (l *disableLogger) Error(_ error) {}

func (l *disableLogger) String(_ string) {}

func (l *disableLogger) Print(_ ...any) {}

func (l *disableLogger) Printf(_ string, _ ...any) {}

func (l *disableLogger) Println(_ ...any) {}

// 空对象构建一个不输出任何内容的实例
func (l *disableLogger) StdLogger() *log.Logger { return log.New(io.Discard, "", 0) }
