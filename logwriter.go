// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/issue9/logs/writer"
)

// 这是对log.Logger的一个包装，将其包装成io.Writer和writer.FlushAdder
// 接口。
//
// log.Logger并未提供更换output的功能，为了达到writer.FlushAdder.Add()
// 函数的功能，只能将所有的log.New()函数缓存起来，直到调用initLogger()
// 时才正直初始log.Logger()实例。但是之后就不能再调用Add()方法添加新的
// io.Writer实例了。
type logWriter struct {
	level  int
	prefix string
	flag   int
	c      *writer.Container
	log    *log.Logger
}

var _ writer.FlushAdder = &logWriter{}
var _ io.Writer = &logWriter{}

// io.Writer.Write()
func (l *logWriter) Write(bs []byte) (int, error) {
	panic("该函数并未真正实现，仅为支持接口而设")
	return 0, nil
}

// writer.FlushAdder.Add()
func (l *logWriter) Add(w io.Writer) error {
	if l.log != nil {
		return errors.New("已经初始化成logger，不能再添加新的io.Writer实例")
	}

	l.c.Add(w)
	return nil
}

// writer.FlushAdder.Flush()
func (l *logWriter) Flush() (int, error) {
	return l.c.Flush()
}

// initLogger 根据当前情况生成log.Logger实例。
// 当多次调用，将触发panic。
func (l *logWriter) initLogger() {
	if l.log != nil {
		panic("log.Logger已经生成")
	}

	l.log = log.New(l.c, l.prefix, l.flag)
}

var flagMap = map[string]int{
	"log.ldate":         log.Ldate,
	"log.ltime":         log.Ltime,
	"log.lmicroseconds": log.Lmicroseconds,
	"log.llongfile":     log.Llongfile,
	"log.lshortfile":    log.Lshortfile,
	"log.lstdflags":     log.LstdFlags,
}

func logWriterInitializer(level int, args map[string]string) (io.Writer, error) {
	flagStr, found := args["flag"]
	if !found || (flagStr == "") {
		flagStr = "log.lstdflags"
	}

	flag, found := flagMap[strings.ToLower(flagStr)]
	if !found {
		return nil, fmt.Errorf("未知的Flag参数:[%v]", flagStr)
	}

	prefix, _ := args["prefix"]

	return &logWriter{
		level:  level,
		flag:   flag,
		prefix: prefix,
		c:      writer.NewContainer(),
	}, nil
}

func init() {
	reg := func(levelName string, level int) {
		fn := func(args map[string]string) (io.Writer, error) {
			return logWriterInitializer(level, args)
		}
		if !Register(levelName, fn) {
			panic(fmt.Sprintf("注册[%v]未成功", levelName))
		}
	}

	reg("info", LevelInfo)
	reg("debug", LevelDebug)
	reg("trace", LevelTrace)
	reg("warn", LevelWarn)
	reg("error", LevelError)
	reg("critical", LevelCritical)
}
