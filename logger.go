// SPDX-License-Identifier: MIT

package logs

import (
	"io"
	"log"
	"sync"

	"github.com/issue9/localeutil"

	"github.com/issue9/logs/v7/writers"
)

var loggerPool = &sync.Pool{New: func() any { return &Logger{} }}

// Logger 日志对象
type Logger struct {
	lv   Level
	logs *Logs
	h    Handler
}

// IsEnable 当前日志是否会真实输出内容
//
// 此返回值与 [Logs.IsEnable] 返回的值是相同的。
func (l *Logger) IsEnable() bool { return l.logs.IsEnable(l.Level()) }

// Level 当前日志的类别
func (l *Logger) Level() Level { return l.lv }

// AppendAttrs 添加新属性
//
// 不会影响之前调用 [Logger.New] 生成的对象。
func (l *Logger) AppendAttrs(attrs map[string]any) {
	l.h = l.h.New(l.logs.detail, l.Level(), map2Slice(l.logs.printer, attrs))
}

// With 创建 [Recorder] 对象
func (l *Logger) With(name string, val any) Recorder {
	if !l.IsEnable() {
		return disabledRecorder
	}

	r := withRecordPool.Get().(*withRecorder)
	r.l = l
	r.r = l.logs.NewRecord().with(name, val)
	return r
}

func (l *Logger) Error(err error) {
	if l.IsEnable() {
		l.logs.NewRecord().DepthError(3, err).Output(l)
	}
}

func (l *Logger) String(s string) {
	if l.IsEnable() {
		l.logs.NewRecord().DepthString(3, s).Output(l)
	}
}

func (l *Logger) LocaleString(s localeutil.Stringer) {
	if l.IsEnable() {
		l.logs.NewRecord().DepthLocaleString(3, s).Output(l)
	}
}

func (l *Logger) Print(v ...any) {
	if l.IsEnable() {
		l.logs.NewRecord().DepthPrint(3, v...).Output(l)
	}
}

func (l *Logger) Println(v ...any) {
	if l.IsEnable() {
		l.logs.NewRecord().DepthPrintln(3, v...).Output(l)
	}
}

func (l *Logger) Printf(format string, v ...any) {
	if l.IsEnable() {
		l.logs.NewRecord().DepthPrintf(3, format, v...).Output(l)
	}
}

// New 根据当前对象派生新的 [Logger]
//
// 新对象会继承当前对象的 [Logger.attrs] 同时还拥有参数 attrs。
func (l *Logger) New(attrs map[string]any) *Logger {
	ll := loggerPool.Get().(*Logger)
	ll.lv = l.lv
	ll.logs = l.logs
	ll.h = l.Handler().New(l.logs.detail, l.Level(), map2Slice(l.logs.printer, attrs))
	return ll
}

// LogLogger 将当前对象转换成标准库的日志对象
//
// NOTE: 不要设置返回对象的 Prefix 和 Flag，这些配置项与当前模块的功能有重叠。
// [log.Logger] 应该仅作为向 [Logger] 输入 [Record.Message] 内容使用。
func (l *Logger) LogLogger() *log.Logger {
	w := io.Discard
	if l.IsEnable() {
		w = l.asWriter()
	}
	return log.New(w, "", 0)
}

// 仅供 [Logger.LogLogger] 使用，因为 depth 值的关系，只有固定的调用层级关系才能正常显示行号。
func (l *Logger) asWriter() io.Writer {
	return writers.WriteFunc(func(data []byte) (int, error) {
		l.logs.NewRecord().DepthString(6, string(data)).Output(l)
		return len(data), nil
	})
}

// Handler 返回关联的 [Handler] 对象
func (l *Logger) Handler() Handler { return l.h }

// 仅用于 AttrLogs 对象中的生成的 Logger
func (l *Logger) free() { loggerPool.Put(l) }
