// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"io"
	"log"

	"github.com/issue9/logs/writers"
)

// 这是对log.New()参数的一个包装。
//
// 封装log.New()的参数，使其在转换成log.Logger实例之前，
// 类似于writers.Container，这样方便通过Register()注册。
type logWriter struct {
	prefix string
	flag   int
	c      *writers.Container
}

func newLogWriter(prefix string, flag int) *logWriter {
	return &logWriter{
		prefix: prefix,
		flag:   flag,
		c:      writers.NewContainer(),
	}
}

func (l *logWriter) Write(bs []byte) (int, error) {
	panic("该函数并未真正实现，仅为支持接口而设")
	return 0, nil
}

func (l *logWriter) Add(w io.Writer) error {
	return l.c.Add(w)
}

func (l *logWriter) toLogger() *log.Logger {
	return log.New(l.c, l.prefix, l.flag)
}
