// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"io"
	"log"

	"github.com/issue9/localeutil"
	"github.com/issue9/sliceutil"

	"github.com/issue9/logs/v6/writers"
)

// Logger 日志对象
type Logger struct {
	lv    Level
	logs  *Logs
	attrs []Attr
}

func appendAttrs(attrs []Attr, p *localeutil.Printer, m map[string]any) []Attr {
	if len(m) == 0 {
		panic("参数 attrs 不能为空")
	}

	for k := range m {
		if sliceutil.Exists(attrs, func(v Attr, _ int) bool { return v.K == k }) {
			panic(fmt.Sprintf("存在同名的元素 %s", k))
		}
	}

	return append(attrs, map2Slice(p, m)...)
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
	l.attrs = appendAttrs(l.attrs, l.logs.printer, attrs)
}

func (l *Logger) With(name string, val any) Recorder {
	if l.IsEnable() {
		return l.newRecord().With(name, val)
	}
	return disabledRecorder
}

func (l *Logger) newRecord() *Record {
	r := l.logs.NewRecord(l.lv)
	for _, p := range l.attrs {
		r.With(p.K, p.V)
	}
	return r
}

func (l *Logger) Error(err error) {
	if l.IsEnable() {
		l.newRecord().DepthError(3, err)
	}
}

func (l *Logger) String(s string) {
	if l.IsEnable() {
		l.newRecord().DepthString(2, s)
	}
}

func (l *Logger) Print(v ...any) {
	if l.IsEnable() {
		l.newRecord().DepthPrint(2, v...)
	}
}

func (l *Logger) Println(v ...any) {
	if l.IsEnable() {
		l.newRecord().DepthPrintln(2, v...)
	}
}

func (l *Logger) Printf(format string, v ...any) {
	if l.IsEnable() {
		l.newRecord().DepthPrintf(2, format, v...)
	}
}

// New 根据当前对象派生新的 [Logger]
//
// 新对象会继承当前对象的 [Logger.attrs] 同时还拥有参数 attrs。
func (l *Logger) New(attrs map[string]any) *Logger {
	pairs := make([]Attr, 0, len(l.attrs)+len(attrs))
	pairs = append(pairs, l.attrs...)
	pairs = appendAttrs(pairs, l.logs.printer, attrs)

	return &Logger{
		lv:    l.lv,
		logs:  l.logs,
		attrs: pairs,
	}
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
		l.newRecord().DepthString(5, string(data))
		return len(data), nil
	})
}
