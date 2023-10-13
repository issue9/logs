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

// WithLocale 是否带本地化信息
//
// 设置了此值为影响以下几个方法中实现了 [localeutil.Stringer] 的参数：
//   - Logger.Error
func WithLocale(p *localeutil.Printer) Option {
	return func(l *Logs) { l.printer = p }
}

// Created 是否显示记录的创建时间
func WithCreated(layout string) Option {
	return func(l *Logs) { l.createdFormat = layout }
}

// Caller 是否显示记录的定位信息
func WithCaller() Option { return func(l *Logs) { l.caller = true } }

// HasCaller 是否包含定位信息
func (logs *Logs) HasCaller() bool { return logs.caller }

// HasCreated 是否包含时间信息
func (logs *Logs) HasCreated() bool { return logs.createdFormat != "" }

// SetCaller 是否显示位置信息
func (logs *Logs) SetCaller(v bool) { logs.caller = v }

// SetCreated 是否显示时间信息
func (logs *Logs) SetCreated(v string) { logs.createdFormat = v }
