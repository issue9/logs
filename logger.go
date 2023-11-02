// SPDX-License-Identifier: MIT

package logs

import (
	"io"
	"log"
)

// TODO 直接使用 nil?
var disabledLogger = &disableLogger{}

type (
	// Logger 日志接口
	Logger interface {
		// With 为日志提供额外的参数
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
		lv    Level
		logs  *Logs
		pairs []Pair
	}

	disableLogger struct{}
)

func (l *logger) StdLogger() *log.Logger { return log.New(l.newRecord().asWriter(), "", 0) }

func (l *logger) With(name string, val any) Logger { return l.newRecord().With(name, val) }

func (l *logger) newRecord() *Record {
	r := l.logs.NewRecord(l.lv)
	for _, p := range l.pairs {
		r.With(p.K, p.V)
	}
	return r
}

func (l *logger) Error(err error) { l.newRecord().DepthError(3, err) }

func (l *logger) String(s string) { l.newRecord().DepthString(2, s) }

func (l *logger) Print(v ...any) { l.newRecord().DepthPrint(2, v...) }

func (l *logger) Println(v ...any) { l.newRecord().DepthPrintln(2, v...) }

func (l *logger) Printf(format string, v ...any) { l.newRecord().DepthPrintf(2, format, v...) }

// With 创建带有指定参数的日志对象
//
// attrs 自动添加的参数，每条日志都将自动带上这些参数；
func (logs *Logs) With(lv Level, attrs map[string]any) Logger {
	if l := logs.level(lv); l == disabledLogger {
		return l
	}

	pairs := make([]Pair, 0, len(attrs)+len(logs.attrs))
	pairs = append(pairs, attrs2Pairs(logs.printer, logs.attrs)...)
	pairs = append(pairs, attrs2Pairs(logs.printer, attrs)...)

	return &logger{
		lv:    lv,
		logs:  logs,
		pairs: pairs,
	}
}

func (l *disableLogger) With(_ string, _ any) Logger { return l }

func (l *disableLogger) Error(_ error) {}

func (l *disableLogger) String(_ string) {}

func (l *disableLogger) Print(_ ...any) {}

func (l *disableLogger) Printf(_ string, _ ...any) {}

func (l *disableLogger) Println(_ ...any) {}

// 空对象构建一个不输出任何内容的实例
func (l *disableLogger) StdLogger() *log.Logger { return log.New(io.Discard, "", 0) }
