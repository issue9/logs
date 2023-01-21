// SPDX-License-Identifier: MIT

package logs

import "fmt"

type Option func(*Logs)

// Printer 对 [Input] 输入的内容进行二次处理
//
// 每个函数分别对 [Input] 相应的输入方法，对其提供的内容进行格式化。
type Printer interface {
	// Error 对 [Input.Error] 提供的内容进行二次处理
	Error(error) string

	// String 对 [Input.String] 提供的内容进行二次处理
	String(string) string

	// Printf 对 [Input.Printf] 提供的内容进行二次处理
	Printf(string, ...interface{}) string

	// Print 对 [Input.Print] 提供的内容进行二次处理
	Print(...interface{}) string
}

type defaultPrinter struct{}

func (p *defaultPrinter) Error(err error) string { return err.Error() }

func (p *defaultPrinter) String(s string) string { return s }

func (p *defaultPrinter) Print(v ...interface{}) string {
	return fmt.Sprint(v...)
}

func (p *defaultPrinter) Printf(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}

func DefaultPrint(l *Logs) { Print(&defaultPrinter{})(l) }

// Print 自定义 [Printer] 接口
func Print(f Printer) Option { return func(l *Logs) { l.printer = f } }

// Caller 是否显示记录的定位信息
func Caller(l *Logs) { l.caller = true }

// Created 是否显示记录的创建时间
func Created(l *Logs) { l.created = true }

// HasCaller 是否包含定位信息
func (logs *Logs) HasCaller() bool { return logs.caller }

// HasCreated 是否包含时间信息
func (logs *Logs) HasCreated() bool { return logs.created }

// SetCaller 是否显示位置信息
func (logs *Logs) SetCaller(v bool) { logs.caller = v }

// SetCreated 是否显示时间信息
func (logs *Logs) SetCreated(v bool) { logs.created = v }

// SetPrinter 设置 [Printer] 对象
//
// 如果 p 为 nil，表示采用默认值。
func (logs *Logs) SetPrinter(p Printer) {
	if p == nil {
		p = &defaultPrinter{}
	}
	logs.printer = p
}
