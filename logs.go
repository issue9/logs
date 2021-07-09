// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/issue9/logs/v3/config"
)

// Logs 日志输出
type Logs struct {
	loggers []*logger
}

// New 声明 Logs 变量
//
// cfg 为配置项，可以为空，表示不输出任何信息，但是 Logs 实例是可用的状态。
func New(cfg *config.Config) (*Logs, error) {
	logs := &Logs{
		loggers: make([]*logger, 0, 6),
	}

	for _, level := range levels {
		logs.loggers = append(logs.loggers, newLogger(level, "", 0))
	}

	if cfg == nil {
		return logs, nil
	}

	for name, c := range cfg.Items {
		index, found := levels[name]
		if !found {
			panic("未知的二级元素名称:" + name)
		}

		ll, err := toWriter(name, c)
		if err != nil {
			return logs, err
		}
		logs.loggers[index] = ll.(*logger)
	}
	return logs, nil
}

// Logger 返回指定级别的日志操作实例
func (l *Logs) Logger(level int) *log.Logger {
	for _, item := range l.loggers {
		if item.level == level {
			return item.Logger
		}
	}
	return nil
}

// SetOutput 设置某一个类型的输出通道
//
// level 表示需要设置的通道，可以是多个值组合，比如 LevelInfo | LevelDebug 。
// 若将 w 设置为 nil 表示关闭此类型的输出。
//
// NOTE: 如果直接调用诸如 ERROR().SetOutput() 设置输出通道，
// 那么 Logs 将失去对该对象的管控，之后再调用 Logs.SetOutput 不会再启作用。
func (l *Logs) SetOutput(level int, w io.Writer) error {
	logs := l.logs(level)
	for _, item := range logs {
		if err := item.SetOutput(w); err != nil {
			return err
		}
	}
	return nil
}

// SetFlags 为所有的日志对象调用 SetFlags
func (l *Logs) SetFlags(level int, flag int) {
	logs := l.logs(level)
	for _, l := range logs {
		l.SetFlags(flag)
	}
}

// SetPrefix 为所有的日志对象调用 SetPrefix
func (l *Logs) SetPrefix(level int, p string) {
	logs := l.logs(level)
	for _, l := range logs {
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
func (l *Logs) INFO() *log.Logger { return l.Logger(LevelInfo) }

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
func (l *Logs) DEBUG() *log.Logger { return l.Logger(LevelDebug) }

// Debug 相当于 DEBUG().Println(v...) 的简写方式
func (l *Logs) Debug(v ...interface{}) { l.DEBUG().Output(2, fmt.Sprintln(v...)) }

// Debugf 相当于 DEBUG().Printf(format, v...) 的简写方式
func (l *Logs) Debugf(format string, v ...interface{}) {
	l.DEBUG().Output(2, fmt.Sprintf(format, v...))
}

// TRACE 获取 TRACE 级别的 log.Logger 实例
func (l *Logs) TRACE() *log.Logger { return l.Logger(LevelTrace) }

// Trace 相当于 TRACE().Println(v...) 的简写方式
func (l *Logs) Trace(v ...interface{}) { l.TRACE().Output(2, fmt.Sprintln(v...)) }

// Tracef 相当于 TRACE().Printf(format, v...) 的简写方式
func (l *Logs) Tracef(format string, v ...interface{}) {
	l.TRACE().Output(2, fmt.Sprintf(format, v...))
}

// WARN 获取 WARN 级别的 log.Logger 实例
func (l *Logs) WARN() *log.Logger { return l.Logger(LevelWarn) }

// Warn 相当于 WARN().Println(v...) 的简写方式
func (l *Logs) Warn(v ...interface{}) { l.WARN().Output(2, fmt.Sprintln(v...)) }

// Warnf 相当于 WARN().Printf(format, v...) 的简写方式
func (l *Logs) Warnf(format string, v ...interface{}) {
	l.WARN().Output(2, fmt.Sprintf(format, v...))
}

// ERROR 获取 ERROR 级别的 log.Logger 实例
func (l *Logs) ERROR() *log.Logger { return l.Logger(LevelError) }

// Error 相当于 ERROR().Println(v...) 的简写方式
func (l *Logs) Error(v ...interface{}) { l.ERROR().Output(2, fmt.Sprintln(v...)) }

// Errorf 相当于 ERROR().Printf(format, v...) 的简写方式
func (l *Logs) Errorf(format string, v ...interface{}) {
	l.ERROR().Output(2, fmt.Sprintf(format, v...))
}

// CRITICAL 获取 CRITICAL 级别的 log.Logger 实例
func (l *Logs) CRITICAL() *log.Logger { return l.Logger(LevelCritical) }

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

func (l *Logs) Print(level int, v ...interface{}) { l.print(level, 3, v...) }

func (l *Logs) Printf(level int, format string, v ...interface{}) {
	l.printf(level, 3, fmt.Sprintf(format, v...))
}

// Fatal 输出错误信息然后退出程序
func (l *Logs) Fatal(level int, code int, v ...interface{}) { l.fatal(level, code, v...) }

// Fatalf 输出错误信息然后退出程序
func (l *Logs) Fatalf(level int, code int, format string, v ...interface{}) {
	l.fatalf(level, code, format, v...)
}

func (l *Logs) panic(level int, v ...interface{}) {
	s := fmt.Sprint(v...)
	l.print(level, 4, s)
	l.Close()
	panic(s)
}

func (l *Logs) panicf(level int, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.print(level, 4, msg)
	l.Close()
	panic(msg)
}

// Panic 输出错误信息然后触发 panic
func (l *Logs) Panic(level int, v ...interface{}) { l.panic(level, v...) }

// Panicf 输出错误信息然后触发 panic
func (l *Logs) Panicf(level int, format string, v ...interface{}) {
	l.panicf(level, format, v...)
}

func (l *Logs) all(msg string) {
	for _, l := range l.loggers {
		l.Output(3, msg)
	}
}

func (l *Logs) print(level, deep int, v ...interface{}) {
	logs := l.logs(level)
	for _, l := range logs {
		l.Output(deep, fmt.Sprintln(v...))
	}
}

func (l *Logs) printf(level, deep int, format string, v ...interface{}) {
	logs := l.logs(level)
	for _, l := range logs {
		l.Output(deep, fmt.Sprintf(format, v...))
	}
}

func (l *Logs) fatal(level int, code int, v ...interface{}) {
	l.print(level, 4, v...)
	l.Close()
	os.Exit(code)
}

func (l *Logs) fatalf(level int, code int, format string, v ...interface{}) {
	l.printf(level, 4, format, v...)
	l.Close()
	os.Exit(code)
}
