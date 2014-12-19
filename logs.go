// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"io"
	"log"
)

// 默认实现
var std *LevelLogger

// 初始化一个默认的LevelLogger实例，之后可以
// 直接调用logs.Info(...)等函数。
func InitFromFile(file string) (err error) {
	if std, err = NewFromFile(file); err != nil {
		return err
	}

	return nil
}

// 从一个xml字符串的reader中初始化levellogger
func InitFromXml(r io.Reader) (err error) {
	if std, err = NewFromXml(r); err != nil {
		return err
	}

	return nil
}

func Flush() (int, error) {
	return std.Flush()
}

// 将指定的level的日志转换成log.Logger实例
// 需要先调用Init(...)函数进行初始化。
func ToStdLogger(level int) (log *log.Logger, ok bool) {
	return std.ToStdLogger(level)
}

// 向指定level的日志输出一行信息
// 需要先调用Init(...)函数进行初始化。
func Println(level int, v ...interface{}) {
	std.Println(level, v...)
}

// 向指定level的日志输出一条信息
// 需要先调用Init(...)函数进行初始化。
func Printf(level int, format string, v ...interface{}) {
	std.Printf(level, format, v...)
}

// Info相当于Println(v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Info(v ...interface{}) {
	std.Info(v...)
}

// Infof根目录于Printf(format, v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

// Debug相当于Println(v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Debug(v ...interface{}) {
	std.Debug(v...)
}

// Debugf根目录于Printf(format, v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

// Trace相当于Println(v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Trace(v ...interface{}) {
	std.Trace(v...)
}

// Tracef根目录于Printf(format, v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Tracef(format string, v ...interface{}) {
	std.Tracef(format, v...)
}

// Warn相当于Println(v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Warn(v ...interface{}) {
	std.Warn(v...)
}

// Warnf根目录于Printf(format, v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Warnf(format string, v ...interface{}) {
	std.Warnf(format, v...)
}

// Error相当于Println(v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Error(v ...interface{}) {
	std.Error(v...)
}

// Errorf根目录于Printf(format, v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

// Critical相当于Println(v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Critical(v ...interface{}) {
	std.Critical(v...)
}

// Criticalf根目录于Printf(format, v...)的简写方式
// 需要先调用Init(...)函数进行初始化。
func Criticalf(format string, v ...interface{}) {
	std.Criticalf(format, v...)
}
