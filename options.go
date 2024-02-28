// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"

	"github.com/issue9/localeutil"
)

// 常用的日志时间格式
const (
	DateMilliLayout = "2006-01-02T15:04:05.000"
	DateMicroLayout = "2006-01-02T15:04:05.000000"
	DateNanoLayout  = "2006-01-02T15:04:05.000000000"

	MilliLayout = "15:04:05.000"
	MicroLayout = "15:04:05.000000"
	NanoLayout  = "15:04:05.000000000"
)

type Option func(*Logs)

// WithStd 是否接管默认的日志处理程序
//
// 如果是 go1.21 之前的版本，会调用 [log.SetOutput] 管理默认日志的输出，输出到 [logs.INFO]；
// 如果是 go1.21 之后的版本，会调用 [slog.SetDefault] 管理默认日志的输出；
func WithStd() Option { return func(l *Logs) { withStd(l) } }

// WithLevels 指定启用的日志级别
//
// 之后也可以通过 [Logs.Enable] 进行修改。
func WithLevels(lv ...Level) Option { return func(l *Logs) { l.levels = lv } }

// WithLocale 指定本地化信息
//
// 如果为 nil，那么将禁用本地化输出，如果多次调用，则以最后一次为准。
//
// 设置了此值为影响以下几个方法中实现了 [localeutil.Stringer] 的参数：
//   - Recorder.Error 中的 error 类型参数；
//   - Recorder.Print/Printf/Println 中的 any 类型参数；
//   - Recorder.With 中的 any 类型参数
func WithLocale(p *localeutil.Printer) Option { return func(l *Logs) { l.printer = p } }

// WithAttrs 为日志添加附加的固定字段
func WithAttrs(attrs map[string]any) Option {
	// NOTE: 无法确保 WithLocale 在 WithAttrs 之前调用，
	// 所以此处直接将 map 类型保存在 Logs，而不是处理后的 slice。

	return func(l *Logs) {
		for k, v := range attrs {
			if _, found := l.attrs[k]; found { // 可能多次调用 WithAttrs 造成重复元素
				panic(fmt.Sprintf("已经存在名称为 %s 的元素", k))
			}
			l.attrs[k] = v
		}
	}
}

// WithDetail 错误信息的调用堆栈
//
// 如果向日志输出的是类型为 err 的信息，是否显示其调用堆栈。
//
// NOTE: 该设置仅对 [Recorder.Error] 方法有效果，
// 如果将 err 传递给 [Recorder.Printf] 等方法，则遵照 [fmt.Appendf] 进行处理。
func WithDetail(v bool) Option { return func(l *Logs) { l.detail = v } }

// WithCreated 指定日期的格式
//
// 如果 layout 为空将会禁用日期显示。
func WithCreated(layout string) Option {
	return func(l *Logs) { l.createdFormat = layout }
}

// WithLocation 是否显示定位信息
func WithLocation(v bool) Option { return func(l *Logs) { l.location = v } }

// HasLocation 是否包含定位信息
func (logs *Logs) HasLocation() bool { return logs.location }

// CreatedFormat created 的时间格式
//
// 如果返回空值，表示禁用在日志中显示时间信息。
func (logs *Logs) CreatedFormat() string { return logs.createdFormat }

// SetLocation 设置是否输出位置信息
func (logs *Logs) SetLocation(v bool) { logs.location = v }

// SetCreated 指定日期的格式
//
// 如果 v 为空将会禁用日期显示。
func (logs *Logs) SetCreated(v string) { logs.createdFormat = v }
