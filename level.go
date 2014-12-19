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

	"github.com/issue9/logs/writer"
)

const (
	LevelInfo = iota
	LevelDebug
	LevelTrace
	LevelWarn
	LevelError
	LevelCritical
)

// 一个分级的日志系统
type LevelLogger struct {
	logs map[int]*logWriter
}

var _ writer.FlushAdder = &LevelLogger{}
var _ io.Writer = &LevelLogger{}

// 从一个xml配置文件初始一个LevelLogger实例
func NewFromFile(file string) (*LevelLogger, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewFromXml(f)
}

// 从一个xml初始化LevelLogger
func NewFromXml(r io.Reader) (*LevelLogger, error) {
	cfg, err := loadFromXml(r)
	if err != nil {
		return nil, err
	}

	w, err := cfg.toWriter()
	if err != nil {
		return nil, err
	}

	log, ok := w.(*LevelLogger)
	if !ok {
		return nil, errors.New("无法转换成*LevelLogger")
	}

	return log, nil
}

// 仅为实现接口，不作任何输出
func (l *LevelLogger) Write(bs []byte) (int, error) {
	panic("该接口没有具体实现，请使用Println()实现相同的功能")
	return 0, nil
}

// w只能是logWriter实例的，否则会返回错误信息。
func (l *LevelLogger) Add(w io.Writer) error {
	lw, ok := w.(*logWriter)
	if !ok {
		return fmt.Errorf("必须为logWriter接口")
	}

	lw.initLogger()
	l.logs[lw.level] = lw
	return nil
}

// writer.FlushAdder.Flush()
func (l *LevelLogger) Flush() (size int, err error) {
	for _, w := range l.logs {
		size, err = w.Flush()
	}
	return
}

// 将指定的level的日志转换成log.Logger实例
func (l *LevelLogger) ToStdLogger(level int) (*log.Logger, bool) {
	w, ok := l.logs[level]
	if !ok {
		return nil, false
	}

	return w.log, true
}

// 向指定level的日志输出一行信息
func (l *LevelLogger) Println(level int, v ...interface{}) {
	if w, found := l.logs[level]; found {
		w.log.Println(v...)
	}
}

// 向指定level的日志输出一条信息
func (l *LevelLogger) Printf(level int, format string, v ...interface{}) {
	if w, found := l.logs[level]; found {
		w.log.Printf(format, v...)
	}
}

// Info相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Info(v ...interface{}) {
	l.Println(LevelInfo, v...)
}

// Infof根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Infof(format string, v ...interface{}) {
	l.Printf(LevelInfo, format, v...)
}

// Debug相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Debug(v ...interface{}) {
	l.Println(LevelDebug, v...)
}

// Debugf根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Debugf(format string, v ...interface{}) {
	l.Printf(LevelDebug, format, v...)
}

// Trace相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Trace(v ...interface{}) {
	l.Println(LevelTrace, v...)
}

// Tracef根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Tracef(format string, v ...interface{}) {
	l.Printf(LevelTrace, format, v...)
}

// Warn相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Warn(v ...interface{}) {
	l.Println(LevelWarn, v...)
}

// Warnf根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Warnf(format string, v ...interface{}) {
	l.Printf(LevelWarn, format, v...)
}

// Error相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Error(v ...interface{}) {
	l.Println(LevelError, v...)
}

// Errorf根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Errorf(format string, v ...interface{}) {
	l.Printf(LevelError, format, v...)
}

// Critical相当于LevelLogger.Println(v...)的简写方式
func (l *LevelLogger) Critical(v ...interface{}) {
	l.Println(LevelCritical, v...)
}

// Criticalf根目录于LevelLogger.Printf(format, v...)的简写方式
func (l *LevelLogger) Criticalf(format string, v ...interface{}) {
	l.Printf(LevelCritical, format, v...)
}

func init() {
	fn := func(args map[string]string) (io.Writer, error) {
		return &LevelLogger{logs: make(map[int]*logWriter)}, nil
	}

	if !Register("logs", fn) {
		panic("无法注册logs初始化函数")
	}
}
