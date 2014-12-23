// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"log"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/logs/writer"
)

var _ writer.WriteAdder = &logWriter{}

// logWriterTestWriter1的输出内容保存在这里
var logWriterTestWriter1Content []byte

type logWriterTestWriter1 struct {
}

func (w *logWriterTestWriter1) Write(bs []byte) (int, error) {
	logWriterTestWriter1Content = append(logWriterTestWriter1Content, bs...)
	return len(bs), nil
}

// logWriterTestWriter2的输出内容保存在这里
var logWriterTestWriter2Content []byte

type logWriterTestWriter2 struct {
}

func (w *logWriterTestWriter2) Write(bs []byte) (int, error) {
	logWriterTestWriter2Content = append(logWriterTestWriter2Content, bs...)
	return len(bs), nil
}

func TestLogWriter(t *testing.T) {
	a := assert.New(t)

	lw := newLogWriter("", log.LstdFlags)
	a.NotNil(lw)

	err := lw.Add(&logWriterTestWriter1{})
	a.NotError(err)
	err = lw.Add(&logWriterTestWriter2{})
	a.NotError(err)

	l := lw.toLogger()
	l.Println("abcd")

	a.True(len(logWriterTestWriter1Content) > 0)
	a.Equal(logWriterTestWriter1Content, logWriterTestWriter2Content)
}
