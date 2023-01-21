// SPDX-License-Identifier: MIT

package logs

import "log"

type (
	// Logger 日志接口
	Logger interface {
		Input

		// StdLogger 将当前对象转换成标准库的日志对象
		//
		// NOTE: 不要设置返回对象的 Prefix 和 Flag，这些配置项与当前模块的功能有重叠。
		// [log.Logger] 应该仅作为向 [Logger] 输入 [Entry.Message] 内容使用。
		StdLogger() *log.Logger
	}

	logger struct {
		lv     Level
		logs   *Logs
		enable bool
		std    *log.Logger
	}

	withLogger struct {
		l     *logger
		std   *log.Logger
		pairs []Pair
	}
)

// Write 实现 io.Writer 供 logs.StdLogger 使用
func (l *logger) Write(data []byte) (int, error) {
	if l.enable {
		ee := l.logs.NewEntry(l.lv)
		ee.Message = string(data)
		ee.Location(4)
		l.logs.Output(ee)
	}
	return len(data), nil
}

func (l *logger) StdLogger() *log.Logger {
	if l.std == nil {
		l.std = log.New(l, "", 0)
	}
	return l.std
}

func (l *logger) With(name string, val interface{}) Input {
	if l.enable {
		return l.logs.NewEntry(l.lv).With(name, val)
	}
	return emptyInputInst
}

func (l *logger) Error(err error) {
	if l.enable {
		l.logs.NewEntry(l.lv).err(3, err)
	}
}

func (l *logger) String(s string) {
	if l.enable {
		l.logs.NewEntry(l.lv).str(4, s)
	}
}

func (l *logger) Print(v ...interface{}) { l.print(4, v...) }

func (l *logger) Printf(format string, v ...interface{}) { l.printf(4, format, v...) }

func (l *logger) print(depth int, v ...interface{}) {
	if l.enable {
		l.logs.NewEntry(l.lv).print(depth, v...)
	}
}

func (l *logger) printf(depth int, format string, v ...interface{}) {
	if l.enable {
		l.logs.NewEntry(l.lv).printf(depth, format, v...)
	}
}

// With 创建带有指定参数的日志对象
//
// params 自动添加的参数，每条日志都将自动带上这些参数；
func (logs *Logs) With(lv Level, params map[string]interface{}) Logger {
	pairs := make([]Pair, 0, len(params))
	for k, v := range params {
		pairs = append(pairs, Pair{K: k, V: v})
	}

	return &withLogger{
		l:     logs.level(lv),
		pairs: pairs,
	}
}

func (l *withLogger) with() *Entry {
	if !l.l.enable {
		return nil
	}

	e := l.l.logs.NewEntry(l.l.lv)
	for _, pair := range l.pairs {
		e.With(pair.K, pair.V)
	}
	return e
}

func (l *withLogger) StdLogger() *log.Logger {
	if l.std == nil {
		l.std = log.New(l, "", 0)
	}
	return l.std
}

func (l *withLogger) With(k string, v interface{}) Input {
	return l.with().With(k, v)
}

func (l *withLogger) Error(err error) { l.with().err(3, err) }

func (l *withLogger) String(s string) { l.with().str(3, s) }

func (l *withLogger) Print(v ...interface{}) { l.with().print(3, v...) }

func (l *withLogger) Printf(format string, v ...interface{}) {
	l.with().printf(3, format, v...)
}

// Write 实现 io.Writer 供 logs.StdLogger 使用
func (l *withLogger) Write(data []byte) (int, error) {
	e := l.with()
	e.Location(4)
	e.Message = string(data)
	l.l.logs.Output(e)
	return len(data), nil
}
