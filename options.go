// SPDX-License-Identifier: MIT

package logs

import "github.com/issue9/localeutil"

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

// WithLocale 指定本地化信息
//
// 如果为 nil，那么将禁用本地化输出。
//
// 设置了此值为影响以下几个方法中实现了 [localeutil.Stringer] 的参数：
//   - Logger.Error
func WithLocale(p *localeutil.Printer) Option {
	return func(l *Logs) { l.printer = p }
}

// WithDetail 错误信息的调用堆栈
//
// 如果向日志输出的是类型为 err 的信息，是否显示其调用堆栈。
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