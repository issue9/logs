// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/issue9/logs/v2/config"
)

// 定义了一些日志的类型
const (
	LevelInfo = iota
	LevelTrace
	LevelDebug
	LevelWarn
	LevelError
	LevelCritical
	levelSize
)

var levels = map[string]int{
	"info":     LevelInfo,
	"trace":    LevelTrace,
	"debug":    LevelDebug,
	"warn":     LevelWarn,
	"error":    LevelError,
	"critical": LevelCritical,
}

var defaultLogs = New()

// Logs 日志输出
type Logs struct {
	loggers []*logger
}

// New 声明 Logs 变量
//
// 需要调用 InitFromXMLFile 或是 InitFromXMLString 进行具体的初始化。
func New() *Logs {
	logs := &Logs{
		loggers: make([]*logger, levelSize, levelSize),
	}

	for index := range logs.loggers {
		logs.loggers[index] = newLogger("", 0)
	}

	return logs
}

// Init 从 config.Config 中初始化整个 logs 系统
func (logs *Logs) Init(cfg *config.Config) error {
	for name, c := range cfg.Items {
		index, found := levels[name]
		if !found {
			panic("未知的二级元素名称:" + name)
		}

		l, err := toWriter(name, c)
		if err != nil {
			return err
		}

		logs.loggers[index] = l.(*logger)
	}

	return nil
}

// InitFromXMLFile 从一个 XML 文件中初始化日志系统。
//
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
//
// Deprecated: 只能由 Init 进行初始化
func (logs *Logs) InitFromXMLFile(path string) error {
	cfg, err := config.ParseXMLFile(path)
	if err != nil {
		return err
	}
	return logs.Init(cfg)
}

// InitFromXMLString 从一个 XML 字符串初始化日志系统。
//
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
//
// Deprecated: 只能由 Init 进行初始化
func (logs *Logs) InitFromXMLString(str string) error {
	cfg, err := config.ParseXMLString(str)
	if err != nil {
		return err
	}
	return logs.Init(cfg)
}

// SetOutput 设置某一个类型的输出通道
//
// 若将 w 设置为 nil 等同于 iotuil.Discard，即关闭此类型的输出。
func (logs *Logs) SetOutput(level int, w io.Writer, prefix string, flag int) error {
	if level < 0 || level > levelSize {
		return errors.New("无效的 level 值")
	}

	logs.loggers[level].setOutput(w, prefix, flag)
	return nil
}

// Flush 输出所有的缓存内容。
// 若是通过 os.Exit() 退出程序的，在执行之前，
// 一定记得调用 Flush() 输出可能缓存的日志内容。
func (logs *Logs) Flush() {
	for _, l := range logs.loggers {
		l.container.Flush()
	}
}

// INFO 获取 INFO 级别的 log.Logger 实例，在未指定 info 级别的日志时，该实例返回一个 nil。
func (logs *Logs) INFO() *log.Logger {
	return logs.loggers[LevelInfo].log
}

// Info 相当于 INFO().Println(v...) 的简写方式
// Info 函数默认是带换行符的，若需要不带换行符的，请使用 DEBUG().Print() 函数代替。
// 其它相似函数也有类型功能。
func (logs *Logs) Info(v ...interface{}) {
	logs.INFO().Output(2, fmt.Sprintln(v...))
}

// Infof 相当于 INFO().Printf(format, v...) 的简写方式
func (logs *Logs) Infof(format string, v ...interface{}) {
	logs.INFO().Output(2, fmt.Sprintf(format, v...))
}

// DEBUG 获取 DEBUG 级别的 log.Logger 实例，在未指定 debug 级别的日志时，该实例返回一个 nil。
func (logs *Logs) DEBUG() *log.Logger {
	return logs.loggers[LevelDebug].log
}

// Debug 相当于 DEBUG().Println(v...) 的简写方式
func (logs *Logs) Debug(v ...interface{}) {
	logs.DEBUG().Output(2, fmt.Sprintln(v...))
}

// Debugf 相当于 DEBUG().Printf(format, v...) 的简写方式
func (logs *Logs) Debugf(format string, v ...interface{}) {
	logs.DEBUG().Output(2, fmt.Sprintf(format, v...))
}

// TRACE 获取 TRACE 级别的 log.Logger 实例，在未指定 trace 级别的日志时，该实例返回一个 nil。
func (logs *Logs) TRACE() *log.Logger {
	return logs.loggers[LevelTrace].log
}

// Trace 相当于 TRACE().Println(v...) 的简写方式
func (logs *Logs) Trace(v ...interface{}) {
	logs.TRACE().Output(2, fmt.Sprintln(v...))
}

// Tracef 相当于 TRACE().Printf(format, v...) 的简写方式
func (logs *Logs) Tracef(format string, v ...interface{}) {
	logs.TRACE().Output(2, fmt.Sprintf(format, v...))
}

// WARN 获取 WARN 级别的 log.Logger 实例，在未指定 warn 级别的日志时，该实例返回一个 nil。
func (logs *Logs) WARN() *log.Logger {
	return logs.loggers[LevelWarn].log
}

// Warn 相当于 WARN().Println(v...) 的简写方式
func (logs *Logs) Warn(v ...interface{}) {
	logs.WARN().Output(2, fmt.Sprintln(v...))
}

// Warnf 相当于 WARN().Printf(format, v...) 的简写方式
func (logs *Logs) Warnf(format string, v ...interface{}) {
	logs.WARN().Output(2, fmt.Sprintf(format, v...))
}

// ERROR 获取 ERROR 级别的 log.Logger 实例，在未指定 error 级别的日志时，该实例返回一个 nil。
func (logs *Logs) ERROR() *log.Logger {
	return logs.loggers[LevelError].log
}

// Error 相当于 ERROR().Println(v...) 的简写方式
func (logs *Logs) Error(v ...interface{}) {
	logs.ERROR().Output(2, fmt.Sprintln(v...))
}

// Errorf 相当于 ERROR().Printf(format, v...) 的简写方式
func (logs *Logs) Errorf(format string, v ...interface{}) {
	logs.ERROR().Output(2, fmt.Sprintf(format, v...))
}

// CRITICAL 获取 CRITICAL 级别的 log.Logger 实例，在未指定 critical 级别的日志时，该实例返回一个 nil。
func (logs *Logs) CRITICAL() *log.Logger {
	return logs.loggers[LevelCritical].log
}

// Critical 相当于 CRITICAL().Println(v...)的简写方式
func (logs *Logs) Critical(v ...interface{}) {
	logs.CRITICAL().Output(2, fmt.Sprintln(v...))
}

// Criticalf 相当于 CRITICAL().Printf(format, v...) 的简写方式
func (logs *Logs) Criticalf(format string, v ...interface{}) {
	logs.CRITICAL().Output(2, fmt.Sprintf(format, v...))
}

// All 向所有的日志输出内容。
func (logs *Logs) All(v ...interface{}) {
	logs.all(v...)
}

// Allf 向所有的日志输出内容。
func (logs *Logs) Allf(format string, v ...interface{}) {
	logs.allf(format, v...)
}

// Fatal 输出错误信息，然后退出程序。
func (logs *Logs) Fatal(code int, v ...interface{}) {
	logs.all(v...)
	logs.Flush()
	os.Exit(code)
}

// Fatalf 输出错误信息，然后退出程序。
func (logs *Logs) Fatalf(code int, format string, v ...interface{}) {
	logs.allf(format, v...)
	logs.Flush()
	os.Exit(code)
}

// Panic 输出错误信息，然后触发 panic。
func (logs *Logs) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	logs.all(s)
	logs.Flush()
	panic(s)
}

// Panicf 输出错误信息，然后触发 panic。
func (logs *Logs) Panicf(format string, v ...interface{}) {
	logs.allf(format, v...)
	logs.Flush()
	panic(fmt.Sprintf(format, v...))
}

func (logs *Logs) all(v ...interface{}) {
	for _, l := range logs.loggers {
		l.log.Output(3, fmt.Sprintln(v...))
	}
}

func (logs *Logs) allf(format string, v ...interface{}) {
	for _, l := range logs.loggers {
		l.log.Output(3, fmt.Sprintf(format, v...))
	}
}

// Init 从 config.Config 中初始化整个 logs 系统
func Init(cfg *config.Config) error {
	return defaultLogs.Init(cfg)
}

// InitFromXMLFile 从一个 XML 文件中初始化日志系统。
//
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
func InitFromXMLFile(path string) error {
	return defaultLogs.InitFromXMLFile(path)
}

// InitFromXMLString 从一个 XML 字符串初始化日志系统。
//
// 再次调用该函数，将会根据新的配置文件重新初始化日志系统。
func InitFromXMLString(str string) error {
	return defaultLogs.InitFromXMLString(str)
}

// SetOutput 设置某一个类型的输出通道
//
// 若将 w 设置为 nil 等同于 iotuil.Discard，即关闭此类型的输出。
func SetOutput(level int, w io.Writer, prefix string, flag int) error {
	return defaultLogs.SetOutput(level, w, prefix, flag)
}

// Flush 输出所有的缓存内容。
// 若是通过 os.Exit() 退出程序的，在执行之前，
// 一定记得调用 Flush() 输出可能缓存的日志内容。
func Flush() {
	defaultLogs.Flush()
}

// INFO 获取 INFO 级别的 log.Logger 实例，在未指定 info 级别的日志时，该实例返回一个 nil。
func INFO() *log.Logger {
	return defaultLogs.INFO()
}

// Info 相当于 INFO().Println(v...) 的简写方式
// Info 函数默认是带换行符的，若需要不带换行符的，请使用 DEBUG().Print() 函数代替。
// 其它相似函数也有类型功能。
func Info(v ...interface{}) {
	defaultLogs.INFO().Output(2, fmt.Sprintln(v...))
}

// Infof 相当于 INFO().Printf(format, v...) 的简写方式
func Infof(format string, v ...interface{}) {
	defaultLogs.INFO().Output(2, fmt.Sprintf(format, v...))
}

// DEBUG 获取 DEBUG 级别的 log.Logger 实例，在未指定 debug 级别的日志时，该实例返回一个 nil。
func DEBUG() *log.Logger {
	return defaultLogs.loggers[LevelDebug].log
}

// Debug 相当于 DEBUG().Println(v...) 的简写方式
func Debug(v ...interface{}) {
	defaultLogs.DEBUG().Output(2, fmt.Sprintln(v...))
}

// Debugf 相当于 DEBUG().Printf(format, v...) 的简写方式
func Debugf(format string, v ...interface{}) {
	defaultLogs.DEBUG().Output(2, fmt.Sprintf(format, v...))
}

// TRACE 获取 TRACE 级别的 log.Logger 实例，在未指定 trace 级别的日志时，该实例返回一个 nil。
func TRACE() *log.Logger {
	return defaultLogs.loggers[LevelTrace].log
}

// Trace 相当于 TRACE().Println(v...) 的简写方式
func Trace(v ...interface{}) {
	defaultLogs.TRACE().Output(2, fmt.Sprintln(v...))
}

// Tracef 相当于 TRACE().Printf(format, v...) 的简写方式
func Tracef(format string, v ...interface{}) {
	defaultLogs.TRACE().Output(2, fmt.Sprintf(format, v...))
}

// WARN 获取 WARN 级别的 log.Logger 实例，在未指定 warn 级别的日志时，该实例返回一个 nil。
func WARN() *log.Logger {
	return defaultLogs.loggers[LevelWarn].log
}

// Warn 相当于 WARN().Println(v...) 的简写方式
func Warn(v ...interface{}) {
	defaultLogs.WARN().Output(2, fmt.Sprintln(v...))
}

// Warnf 相当于 WARN().Printf(format, v...) 的简写方式
func Warnf(format string, v ...interface{}) {
	defaultLogs.WARN().Output(2, fmt.Sprintf(format, v...))
}

// ERROR 获取 ERROR 级别的 log.Logger 实例，在未指定 error 级别的日志时，该实例返回一个 nil。
func ERROR() *log.Logger {
	return defaultLogs.loggers[LevelError].log
}

// Error 相当于 ERROR().Println(v...) 的简写方式
func Error(v ...interface{}) {
	defaultLogs.ERROR().Output(2, fmt.Sprintln(v...))
}

// Errorf 相当于 ERROR().Printf(format, v...) 的简写方式
func Errorf(format string, v ...interface{}) {
	defaultLogs.ERROR().Output(2, fmt.Sprintf(format, v...))
}

// CRITICAL 获取 CRITICAL 级别的 log.Logger 实例，在未指定 critical 级别的日志时，该实例返回一个 nil。
func CRITICAL() *log.Logger {
	return defaultLogs.loggers[LevelCritical].log
}

// Critical 相当于 CRITICAL().Println(v...)的简写方式
func Critical(v ...interface{}) {
	defaultLogs.CRITICAL().Output(2, fmt.Sprintln(v...))
}

// Criticalf 相当于 CRITICAL().Printf(format, v...) 的简写方式
func Criticalf(format string, v ...interface{}) {
	defaultLogs.CRITICAL().Output(2, fmt.Sprintf(format, v...))
}

// All 向所有的日志输出内容。
func All(v ...interface{}) {
	defaultLogs.All(v...)
}

// Allf 向所有的日志输出内容。
func Allf(format string, v ...interface{}) {
	defaultLogs.Allf(format, v...)
}

// Fatal 输出错误信息，然后退出程序。
func Fatal(code int, v ...interface{}) {
	defaultLogs.Fatal(code, v...)
}

// Fatalf 输出错误信息，然后退出程序。
func Fatalf(code int, format string, v ...interface{}) {
	defaultLogs.Fatalf(code, format, v...)
}

// Panic 输出错误信息，然后触发 panic。
func Panic(v ...interface{}) {
	defaultLogs.Panic(v...)
}

// Panicf 输出错误信息，然后触发 panic。
func Panicf(format string, v ...interface{}) {
	defaultLogs.Panicf(format, v...)
}
