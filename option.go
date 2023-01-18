// SPDX-License-Identifier: MIT

package logs

import "fmt"

type Option func(*Logs)

// Printer 对 [Logger] 输入的内容进行格式化
//
// 每个函数分别对 [Logger] 相应的输入方法，对其提供的内容进行格式化。
type Printer interface {
	// Error 格式由 [Logger.Error] 提供的内容
	Error(error) string

	// Printf 格式化由 [Logger.Printf] 提供的内容
	Printf(string, ...interface{}) string

	// Print 格式化由 [Logger.Print] 提供的内容
	Print(...interface{}) string
}

type defaultPrinter struct{}

// Caller 是否显示记录的定位信息
func Caller(l *Logs) { l.caller = true }

// Created 是否显示记录的创建时间
func Created(l *Logs) { l.created = true }

func DefaultPrint(l *Logs) { Print(&defaultPrinter{})(l) }

// Print 自定义 [Printer] 接口
func Print(f Printer) Option { return func(l *Logs) { l.printer = f } }

// HasCaller 是否包含定位信息
func (logs *Logs) HasCaller() bool { return logs.caller }

// HasCreated 是否包含时间信息
func (logs *Logs) HasCreated() bool { return logs.created }

func (f *defaultPrinter) Error(err error) string { return err.Error() }

func (f *defaultPrinter) Print(v ...interface{}) string {
	return fmt.Sprint(v...)
}

func (f *defaultPrinter) Printf(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}
