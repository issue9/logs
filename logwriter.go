// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"io"
	"log"

	"github.com/issue9/logs/writer"
)

// 这是对log.New()参数的一个包装。
//
// log.Logger并未提供更换output的功能，
// 为了达到writer.FlushAdder.Add()函数的功能，
// 只能将所有的log.New()函数参数缓存起来，
// 直到调用toLogger()时才正直初始化成log.Logger()实例。
// 但是之后就不能再调用Add()方法添加新的io.Writer实例了。
type logWriter struct {
	prefix string
	flag   int
	c      *writer.Container
}

func newLogWriter(prefix string, flag int) *logWriter {
	return &logWriter{
		prefix: prefix,
		flag:   flag,
		c:      writer.NewContainer(),
	}
}

// io.Writer.Write()
func (l *logWriter) Write(bs []byte) (int, error) {
	panic("该函数并未真正实现，仅为支持接口而设")
	return 0, nil
}

// writer.FlushAdder.Add()
func (l *logWriter) Add(w io.Writer) error {
	return l.c.Add(w)
}

func (l *logWriter) toLogger() *log.Logger {
	return log.New(l.c, l.prefix, l.flag)
}
