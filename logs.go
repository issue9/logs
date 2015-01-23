// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/issue9/logs/internal/config"
	"github.com/issue9/logs/writers"
)

// 默认所有日志的写入文件。
var discardLog = log.New(ioutil.Discard, "", log.LstdFlags)

// 保存info,warn等6个预定义log.Logger的io.Writer接口实例，
// 方便在关闭日志时，输出其中缓存的内容。
var conts = writers.NewContainer()

// 预定义的6个log.Logger实例。
var (
	info     = discardLog
	warn     = discardLog
	_error   = discardLog
	debug    = discardLog
	trace    = discardLog
	critical = discardLog
)

// 从一个xml文件中初始化日志系统。
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
func InitFromXMLFile(path string) error {
	cfg, err := config.ParseXMLFile(path)
	if err != nil {
		return err
	}
	return initFromConfig(cfg)
}

// 从一个xml字符串初始化日志系统。
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
func InitFromXMLString(xml string) error {
	cfg, err := config.ParseXMLString(xml)
	if err != nil {
		return err
	}
	return initFromConfig(cfg)
}

// 从config.Config中初始化整个log系统
func initFromConfig(cfg *config.Config) error {
	if conts.Len() > 0 { // 加载新配置文件。先输出旧的内容。
		Flush()
		conts.Clear()
		// 重置为空值
		info = discardLog
		critical = discardLog
		debug = discardLog
		trace = discardLog
		warn = discardLog
		_error = discardLog
	}

	for name, c := range cfg.Items {
		writer, err := toWriter(c)
		if err != nil {
			return err
		}

		w, ok := writer.(*logWriter)
		if !ok {
			return errors.New("initFromConfig:二级元素必须为logWriter类型")
		}
		switch name {
		case "info":
			info = w.toLogger()
		case "warn":
			warn = w.toLogger()
		case "debug":
			debug = w.toLogger()
		case "error":
			_error = w.toLogger()
		case "trace":
			trace = w.toLogger()
		case "critical":
			critical = w.toLogger()
		}
		conts.Add(w.c)
	}

	return nil
}

// 输出所有的缓存内容。
func Flush() {
	conts.Flush()
}

// 返回info日志实例。
func INFO() *log.Logger {
	return info
}

// Info相当于info.Println(v...)的简写方式
func Info(v ...interface{}) {
	info.Println(v...)
}

// Infof相当于info.Printf(format, v...)的简写方式
func Infof(format string, v ...interface{}) {
	info.Printf(format, v...)
}

// 返回debug级别的日志。
func DEBUG() *log.Logger {
	return debug
}

// Debug相当于debug.Println(v...)的简写方式
func Debug(v ...interface{}) {
	debug.Println(v...)
}

// Debugf相当于debug.Printf(format, v...)的简写方式
func Debugf(format string, v ...interface{}) {
	debug.Printf(format, v...)
}

// 返回trace级别的日志。
func TRACE() *log.Logger {
	return trace
}

// Trace相当于trace.Println(v...)的简写方式
func Trace(v ...interface{}) {
	trace.Println(v...)
}

// Tracef相当于trace.Printf(format, v...)的简写方式
func Tracef(format string, v ...interface{}) {
	trace.Printf(format, v...)
}

// 返回warn级别的日志。
func WARN() *log.Logger {
	return warn
}

// Warn相当于warn.Println(v...)的简写方式
func Warn(v ...interface{}) {
	warn.Println(v...)
}

// Warnf相当于warn.Printf(format, v...)的简写方式
func Warnf(format string, v ...interface{}) {
	warn.Printf(format, v...)
}

// 返回error级别的日志。
func ERROR() *log.Logger {
	return _error
}

// Error相当于_error.Println(v...)的简写方式
func Error(v ...interface{}) {
	_error.Println(v...)
}

// Errorf相当于_error.Printf(format, v...)的简写方式
func Errorf(format string, v ...interface{}) {
	_error.Printf(format, v...)
}

// 返回critical级别的日志。
func CRITICAL() *log.Logger {
	return critical
}

// Critical相当于critical.Println(v...)的简写方式
func Critical(v ...interface{}) {
	critical.Println(v...)
}

// Criticalf相当于critical.Printf(format, v...)的简写方式
func Criticalf(format string, v ...interface{}) {
	critical.Printf(format, v...)
}

// 向所有的日志输出内容。
func All(v ...interface{}) {
	Info(v...)
	Debug(v...)
	Trace(v...)
	Warn(v...)
	Error(v...)
	Critical(v...)
}

// 向所有的日志输出内容。
func Allf(format string, v ...interface{}) {
	Infof(format, v...)
	Debugf(format, v...)
	Tracef(format, v...)
	Warnf(format, v...)
	Errorf(format, v...)
	Criticalf(format, v...)
}
