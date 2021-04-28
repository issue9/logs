// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/issue9/logs/v3/config"
)

// 目前支持的日志类型
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
func New() *Logs {
	logs := &Logs{
		loggers: make([]*logger, levelSize),
	}

	for index := range logs.loggers {
		logs.loggers[index] = newLogger("", 0)
	}

	return logs
}

// Init 从 config.Config 中初始化整个 logs 系统
func (l *Logs) Init(cfg *config.Config) error {
	for name, c := range cfg.Items {
		index, found := levels[name]
		if !found {
			panic("未知的二级元素名称:" + name)
		}

		ll, err := toWriter(name, c)
		if err != nil {
			return err
		}
		l.loggers[index] = ll.(*logger)
	}

	return nil
}

// SetOutput 设置某一个类型的输出通道
//
// 若将 w 设置为 nil 表示关闭此类型的输出。
//
// NOTE: 如果直接调用诸如 ERROR().SetOutput() 设置输出通道，
// 那么 Logs 将失去对该对象的管控，之后再调用 Logs.SetOutput 不会再启作用。
func (l *Logs) SetOutput(level int, w io.Writer) error {
	if level >= LevelInfo && level < levelSize {
		return l.loggers[level].SetOutput(w)
	}
	panic(fmt.Sprintf("无效的 level 值：%d", level))
}

// SetFlags 为所有的日志对象调用 SetFlags
func (l *Logs) SetFlags(flag int) {
	for _, l := range l.loggers {
		l.SetFlags(flag)
	}
}

// SetPrefix 为所有的日志对象调用 SetPrefix
func (l *Logs) SetPrefix(p string) {
	for _, l := range l.loggers {
		l.SetPrefix(p)
	}
}

// Flush 输出所有的缓存内容
func (l *Logs) Flush() error {
	for _, l := range l.loggers {
		if err := l.container.Flush(); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭所有的输出通道
//
// 若是通过 os.Exit() 退出程序的，在执行之前，
// 一定记得调用 Close() 输出可能缓存的日志内容。
func (l *Logs) Close() error {
	for _, l := range l.loggers {
		if err := l.container.Close(); err != nil {
			return err
		}
	}
	return nil
}

// INFO 获取 INFO 级别的 log.Logger 实例
func (l *Logs) INFO() *log.Logger { return l.loggers[LevelInfo].Logger }

// Info 相当于 INFO().Println(v...) 的简写方式
//
// Info 函数默认是带换行符的，若需要不带换行符的，请使用 DEBUG().Print() 函数代替。
// 其它相似函数也有类型功能。
func (l *Logs) Info(v ...interface{}) { l.INFO().Output(2, fmt.Sprintln(v...)) }

// Infof 相当于 INFO().Printf(format, v...) 的简写方式
func (l *Logs) Infof(format string, v ...interface{}) {
	l.INFO().Output(2, fmt.Sprintf(format, v...))
}

// DEBUG 获取 DEBUG 级别的 log.Logger 实例
func (l *Logs) DEBUG() *log.Logger { return l.loggers[LevelDebug].Logger }

// Debug 相当于 DEBUG().Println(v...) 的简写方式
func (l *Logs) Debug(v ...interface{}) { l.DEBUG().Output(2, fmt.Sprintln(v...)) }

// Debugf 相当于 DEBUG().Printf(format, v...) 的简写方式
func (l *Logs) Debugf(format string, v ...interface{}) {
	l.DEBUG().Output(2, fmt.Sprintf(format, v...))
}

// TRACE 获取 TRACE 级别的 log.Logger 实例
func (l *Logs) TRACE() *log.Logger { return l.loggers[LevelTrace].Logger }

// Trace 相当于 TRACE().Println(v...) 的简写方式
func (l *Logs) Trace(v ...interface{}) { l.TRACE().Output(2, fmt.Sprintln(v...)) }

// Tracef 相当于 TRACE().Printf(format, v...) 的简写方式
func (l *Logs) Tracef(format string, v ...interface{}) {
	l.TRACE().Output(2, fmt.Sprintf(format, v...))
}

// WARN 获取 WARN 级别的 log.Logger 实例
func (l *Logs) WARN() *log.Logger { return l.loggers[LevelWarn].Logger }

// Warn 相当于 WARN().Println(v...) 的简写方式
func (l *Logs) Warn(v ...interface{}) { l.WARN().Output(2, fmt.Sprintln(v...)) }

// Warnf 相当于 WARN().Printf(format, v...) 的简写方式
func (l *Logs) Warnf(format string, v ...interface{}) {
	l.WARN().Output(2, fmt.Sprintf(format, v...))
}

// ERROR 获取 ERROR 级别的 log.Logger 实例
//
// 在未指定 error 级别的日志时，该实例返回一个 nil。
func (l *Logs) ERROR() *log.Logger { return l.loggers[LevelError].Logger }

// Error 相当于 ERROR().Println(v...) 的简写方式
func (l *Logs) Error(v ...interface{}) { l.ERROR().Output(2, fmt.Sprintln(v...)) }

// Errorf 相当于 ERROR().Printf(format, v...) 的简写方式
func (l *Logs) Errorf(format string, v ...interface{}) {
	l.ERROR().Output(2, fmt.Sprintf(format, v...))
}

// CRITICAL 获取 CRITICAL 级别的 log.Logger 实例
func (l *Logs) CRITICAL() *log.Logger { return l.loggers[LevelCritical].Logger }

// Critical 相当于 CRITICAL().Println(v...)的简写方式
func (l *Logs) Critical(v ...interface{}) { l.CRITICAL().Output(2, fmt.Sprintln(v...)) }

// Criticalf 相当于 CRITICAL().Printf(format, v...) 的简写方式
func (l *Logs) Criticalf(format string, v ...interface{}) {
	l.CRITICAL().Output(2, fmt.Sprintf(format, v...))
}

// All 向所有的日志输出内容
func (l *Logs) All(v ...interface{}) { l.all(fmt.Sprint(v...)) }

// Allf 向所有的日志输出内容
func (l *Logs) Allf(format string, v ...interface{}) {
	l.all(fmt.Sprintf(format, v...))
}

// Fatal 输出错误信息然后退出程序
func (l *Logs) Fatal(code int, v ...interface{}) {
	l.all(fmt.Sprintln(v...))
	l.Close()
	os.Exit(code)
}

// Fatalf 输出错误信息然后退出程序
func (l *Logs) Fatalf(code int, format string, v ...interface{}) {
	l.all(fmt.Sprintf(format, v...))
	l.Close()
	os.Exit(code)
}

// Panic 输出错误信息然后触发 panic
func (l *Logs) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.all(s)
	l.Close()
	panic(s)
}

// Panicf 输出错误信息然后触发 panic
func (l *Logs) Panicf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.all(msg)
	l.Close()
	panic(msg)
}

func (l *Logs) all(msg string) {
	for _, l := range l.loggers {
		l.Output(3, msg)
	}
}

// Default 返回当前模块中全局函数使用的 *Logs 对象
func Default() *Logs { return defaultLogs }

// Init 从 config.Config 中初始化整个 logs 系统
func Init(cfg *config.Config) error { return Default().Init(cfg) }

// Flush 输出所有的缓存内容
func Flush() error { return Default().Flush() }

// Close 关闭所有的输出通道
//
// 若是通过 os.Exit() 退出程序的，在执行之前，
// 一定记得调用 Close() 输出可能缓存的日志内容。
func Close() error { return Default().Close() }

// INFO 获取 INFO 级别的 log.Logger 实例
func INFO() *log.Logger { return Default().INFO() }

// Info 相当于 INFO().Println(v...) 的简写方式
//
// Info 函数默认是带换行符的，若需要不带换行符的，请使用 DEBUG().Print() 函数代替。
// 其它相似函数也有类型功能。
func Info(v ...interface{}) { INFO().Output(2, fmt.Sprintln(v...)) }

// Infof 相当于 INFO().Printf(format, v...) 的简写方式
func Infof(format string, v ...interface{}) {
	INFO().Output(2, fmt.Sprintf(format, v...))
}

// DEBUG 获取 DEBUG 级别的 log.Logger 实例
func DEBUG() *log.Logger { return Default().DEBUG() }

// Debug 相当于 DEBUG().Println(v...) 的简写方式
func Debug(v ...interface{}) { DEBUG().Output(2, fmt.Sprintln(v...)) }

// Debugf 相当于 DEBUG().Printf(format, v...) 的简写方式
func Debugf(format string, v ...interface{}) {
	DEBUG().Output(2, fmt.Sprintf(format, v...))
}

// TRACE 获取 TRACE 级别的 log.Logger 实例
func TRACE() *log.Logger { return Default().TRACE() }

// Trace 相当于 TRACE().Println(v...) 的简写方式
func Trace(v ...interface{}) { TRACE().Output(2, fmt.Sprintln(v...)) }

// Tracef 相当于 TRACE().Printf(format, v...) 的简写方式
func Tracef(format string, v ...interface{}) {
	TRACE().Output(2, fmt.Sprintf(format, v...))
}

// WARN 获取 WARN 级别的 log.Logger 实例
func WARN() *log.Logger { return Default().WARN() }

// Warn 相当于 WARN().Println(v...) 的简写方式
func Warn(v ...interface{}) { WARN().Output(2, fmt.Sprintln(v...)) }

// Warnf 相当于 WARN().Printf(format, v...) 的简写方式
func Warnf(format string, v ...interface{}) {
	WARN().Output(2, fmt.Sprintf(format, v...))
}

// ERROR 获取 ERROR 级别的 log.Logger 实例
func ERROR() *log.Logger { return Default().ERROR() }

// Error 相当于 ERROR().Println(v...) 的简写方式
func Error(v ...interface{}) { ERROR().Output(2, fmt.Sprintln(v...)) }

// Errorf 相当于 ERROR().Printf(format, v...) 的简写方式
func Errorf(format string, v ...interface{}) {
	ERROR().Output(2, fmt.Sprintf(format, v...))
}

// CRITICAL 获取 CRITICAL 级别的 log.Logger 实例
func CRITICAL() *log.Logger { return Default().CRITICAL() }

// Critical 相当于 CRITICAL().Println(v...)的简写方式
func Critical(v ...interface{}) { CRITICAL().Output(2, fmt.Sprintln(v...)) }

// Criticalf 相当于 CRITICAL().Printf(format, v...) 的简写方式
func Criticalf(format string, v ...interface{}) {
	CRITICAL().Output(2, fmt.Sprintf(format, v...))
}

// All 向所有的日志输出内容
func All(v ...interface{}) { Default().all(fmt.Sprintln(v...)) }

// Allf 向所有的日志输出内容
func Allf(format string, v ...interface{}) { Default().all(fmt.Sprintf(format, v...)) }

// Fatal 输出错误信息然后退出程序
func Fatal(code int, v ...interface{}) {
	Default().all(fmt.Sprint(v...))
	Close()
	os.Exit(code)
}

// Fatalf 输出错误信息然后退出程序
func Fatalf(code int, format string, v ...interface{}) {
	Default().all(fmt.Sprintf(format, v...))
	Close()
	os.Exit(code)
}

// Panic 输出错误信息然后触发 panic
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	Default().all(s)
	Close()
	panic(s)
}

// Panicf 输出错误信息然后触发 panic
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	Default().all(s)
	Close()
	panic(s)
}
