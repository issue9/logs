// SPDX-License-Identifier: MIT

package logs

import (
	"io"
	"log"

	"github.com/issue9/logs/v6/writers"
)

// Logger 日志对象
type Logger struct {
	lv    Level
	logs  *Logs
	pairs []Pair
}

func (l *Logger) isEnable() bool { return l.logs.IsEnable(l.Level()) }

func (l *Logger) Level() Level { return l.lv }

func (l *Logger) With(name string, val any) Recorder {
	if l.isEnable() {
		return l.newRecord().With(name, val)
	}
	return disabledRecorder
}

func (l *Logger) newRecord() *Record {
	r := l.logs.NewRecord(l.lv)
	for _, p := range l.pairs {
		r.With(p.K, p.V)
	}
	return r
}

func (l *Logger) Error(err error) {
	if l.isEnable() {
		l.newRecord().DepthError(3, err)
	}
}

func (l *Logger) String(s string) {
	if l.isEnable() {
		l.newRecord().DepthString(2, s)
	}
}

func (l *Logger) Print(v ...any) {
	if l.isEnable() {
		l.newRecord().DepthPrint(2, v...)
	}
}

func (l *Logger) Println(v ...any) {
	if l.isEnable() {
		l.newRecord().DepthPrintln(2, v...)
	}
}

func (l *Logger) Printf(format string, v ...any) {
	if l.isEnable() {
		l.newRecord().DepthPrintf(2, format, v...)
	}
}

// New 根据当前对象派生新的 [Logger]
//
// 新对象会继承当前对象的 [Logger.attrs] 同时还拥有参数 attrs。
func (l *Logger) New(attrs map[string]any) *Logger {
	if len(attrs) == 0 {
		panic("参数 attrs 不能为空")
	}

	pairs := make([]Pair, 0, len(l.pairs)+len(attrs))
	pairs = append(pairs, l.pairs...)
	pairs = append(pairs, attrs2Pairs(l.logs.printer, attrs)...)

	return &Logger{
		lv:    l.lv,
		logs:  l.logs,
		pairs: pairs,
	}
}

// StdLogger 将当前对象转换成标准库的日志对象
//
// NOTE: 不要设置返回对象的 Prefix 和 Flag，这些配置项与当前模块的功能有重叠。
// [log.Logger] 应该仅作为向 [Logger] 输入 [Record.Message] 内容使用。
func (l *Logger) StdLogger() *log.Logger {
	w := io.Discard
	if l.isEnable() {
		w = l.asWriter()
	}
	return log.New(w, "", 0)
}

// 转换成 io.Writer
//
// 仅供内部使用，因为 depth 值的关系，只有固定的调用层级关系才能正常显示行号。
func (l *Logger) asWriter() io.Writer {
	return writers.WriteFunc(func(data []byte) (int, error) {
		l.newRecord().DepthString(5, string(data))
		return len(data), nil
	})
}
