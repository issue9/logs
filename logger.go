// SPDX-License-Identifier: MIT

package logs

import (
	"io"
	"log"

	"github.com/issue9/logs/v3/writers"
)

// 扩展 log.Logger，使可以同时输出到多个日志通道
type logger struct {
	*log.Logger
	container *writers.Container
}

func newLogger(prefix string, flag int) *logger {
	c := writers.NewContainer()
	return &logger{
		Logger:    log.New(c, prefix, flag),
		container: c,
	}
}

// SetOutput 重新设置输出通道
//
// 如果还有内容未输出，则会先输出内容。 如果 w 为 nil，取消该通道的输出。
func (l *logger) SetOutput(w io.Writer) error {
	if err := l.container.Flush(); err != nil {
		return err
	}

	l.container.Clear()

	if w == nil {
		return nil
	}

	return l.container.Add(w)
}
