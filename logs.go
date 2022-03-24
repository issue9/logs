// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Logs 日志输出
type Logs struct {
	loggers map[int]*logger
}

func New() (*Logs, error) {
	return &Logs{
		loggers: map[int]*logger{
			LevelInfo:     newLogger("", 0),
			LevelTrace:    newLogger("", 0),
			LevelDebug:    newLogger("", 0),
			LevelWarn:     newLogger("", 0),
			LevelError:    newLogger("", 0),
			LevelCritical: newLogger("", 0),
		},
	}, nil
}

// Logger 返回指定级别的日志操作实例
//
// level 不能以组合的形式出现；
func (l *Logs) Logger(level int) *log.Logger {
	for key, item := range l.loggers {
		if key == level {
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
	return l.walk(level, func(ll *logger) error {
		return ll.SetOutput(w)
	})
}

// SetFlags 为所选的日志对象调用 SetFlags
func (l *Logs) SetFlags(level int, flag int) {
	l.walk(level, func(ll *logger) error {
		ll.SetFlags(flag)
		return nil
	})
}

// SetPrefix 为所选的日志对象调用 SetPrefix
func (l *Logs) SetPrefix(level int, p string) {
	l.walk(level, func(ll *logger) error {
		ll.SetPrefix(p)
		return nil
	})
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

// Fatal 输出错误信息然后退出程序
func (l *Logs) Fatal(level int, code int, v ...interface{}) {
	l.Print(level, 1, v...)
	l.Flush()
	os.Exit(code)
}

// Fatalf 输出错误信息然后退出程序
func (l *Logs) Fatalf(level int, code int, format string, v ...interface{}) {
	l.Printf(level, 1, format, v...)
	l.Flush()
	os.Exit(code)
}

// Panic 输出错误信息然后触发 panic
func (l *Logs) Panic(level int, v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Print(level, 1, s)
	l.Flush()
	panic(s)
}

// Panicf 输出错误信息然后触发 panic
func (l *Logs) Panicf(level int, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Print(level, 1, msg)
	l.Flush()
	panic(msg)
}

func (l *Logs) all(msg string) {
	for _, l := range l.loggers {
		l.Output(3, msg)
	}
}

// Print 向指定的通道输出信息
//
// level 表示需要设置的通道，可以是多个值组合，比如 LevelInfo | LevelDebug；
// deep 为 0 时，表示调用者；
func (l *Logs) Print(level, deep int, v ...interface{}) {
	deep += 4 // 保证 walk 为 1
	l.walk(level, func(ll *logger) error {
		return ll.Output(deep, fmt.Sprintln(v...))
	})
}

// Printf 向指定的通道输出信息
//
// level 表示需要设置的通道，可以是多个值组合，比如 LevelInfo | LevelDebug；
// deep 为 0 时，表示调用者；
func (l *Logs) Printf(level, deep int, format string, v ...interface{}) {
	deep += 4 // 保证 walk 为 1
	l.walk(level, func(ll *logger) error {
		return ll.Output(deep, fmt.Sprintf(format, v...))
	})
}
