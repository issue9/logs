// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/issue9/logs/internal/config"
	"github.com/issue9/logs/writers"
)

// 保存info,warn等6个预定义log.Logger的io.Writer接口实例，
// 方便在关闭日志时，输出其中缓存的内容。
var conts = writers.NewContainer()

// 预定义的6个log.Logger实例。
var (
	info     *log.Logger = nil
	warn     *log.Logger = nil
	_error   *log.Logger = nil
	debug    *log.Logger = nil
	trace    *log.Logger = nil
	critical *log.Logger = nil
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
		info = nil
		critical = nil
		debug = nil
		trace = nil
		warn = nil
		_error = nil
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
// 若是通过os.Exit()退出程序的，在执行之前，
// 一定记得调用Flush()输出可能缓存的日志内容。
func Flush() {
	conts.Flush()
}

// 获取INFO级别的log.Logger实例，在未指定info级别的日志时，该实例返回一个nil。
func INFO() *log.Logger {
	return info
}

// Info相当于INFO().Println(v...)的简写方式
func Info(v ...interface{}) {
	if info == nil {
		return
	}

	info.Println(v...)
}

// Infof相当于INFO().Printf(format, v...)的简写方式
func Infof(format string, v ...interface{}) {
	if info == nil {
		return
	}

	info.Printf(format, v...)
}

// 获取DEBUG级别的log.Logger实例，在未指定debug级别的日志时，该实例返回一个nil。
func DEBUG() *log.Logger {
	return debug
}

// Debug相当于DEBUG().Println(v...)的简写方式
func Debug(v ...interface{}) {
	if debug == nil {
		return
	}

	debug.Println(v...)
}

// Debugf相当于DEBUG().Printf(format, v...)的简写方式
func Debugf(format string, v ...interface{}) {
	if debug == nil {
		return
	}

	debug.Printf(format, v...)
}

// 获取TRACE级别的log.Logger实例，在未指定trace级别的日志时，该实例返回一个nil。
func TRACE() *log.Logger {
	return trace
}

// Trace相当于TRACE().Println(v...)的简写方式
func Trace(v ...interface{}) {
	if trace == nil {
		return
	}

	trace.Println(v...)
}

// Tracef相当于TRACE().Printf(format, v...)的简写方式
func Tracef(format string, v ...interface{}) {
	if trace == nil {
		return
	}

	trace.Printf(format, v...)
}

// 获取WARN级别的log.Logger实例，在未指定warn级别的日志时，该实例返回一个nil。
func WARN() *log.Logger {
	return warn
}

// Warn相当于WARN().Println(v...)的简写方式
func Warn(v ...interface{}) {
	if warn == nil {
		return
	}

	warn.Println(v...)
}

// Warnf相当于WARN().Printf(format, v...)的简写方式
func Warnf(format string, v ...interface{}) {
	if warn == nil {
		return
	}

	warn.Printf(format, v...)
}

// 获取ERROR级别的log.Logger实例，在未指定error级别的日志时，该实例返回一个nil。
func ERROR() *log.Logger {
	return _error
}

// Error相当于ERROR().Println(v...)的简写方式
func Error(v ...interface{}) {
	if _error == nil {
		return
	}

	_error.Println(v...)
}

// Errorf相当于ERROR().Printf(format, v...)的简写方式
func Errorf(format string, v ...interface{}) {
	if _error == nil {
		return
	}

	_error.Printf(format, v...)
}

// 获取CRITICAL级别的log.Logger实例，在未指定critical级别的日志时，该实例返回一个nil。
func CRITICAL() *log.Logger {
	return critical
}

// Critical相当于CRITICAL().Println(v...)的简写方式
func Critical(v ...interface{}) {
	if critical == nil {
		return
	}

	critical.Println(v...)
}

// Criticalf相当于CRITICAL().Printf(format, v...)的简写方式
func Criticalf(format string, v ...interface{}) {
	if critical == nil {
		return
	}

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

// 输出错误信息，然后退出程序。
func Fatal(v ...interface{}) {
	All(v...)
	Flush()
	os.Exit(2)
}

// 输出错误信息，然后退出程序。
func Fatalf(format string, v ...interface{}) {
	Allf(format, v...)
	Flush()
	os.Exit(2)
}

// 输出错误信息，然后触发panic。
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	All(s)
	Flush()
	panic(s)
}

// 输出错误信息，然后触发panic。
func Panicf(format string, v ...interface{}) {
	Allf(format, v...)
	Flush()
	panic(fmt.Sprintf(format, v...))
}
