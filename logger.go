// SPDX-License-Identifier: MIT

package logs

import (
	"errors"
	"io"
	"log"
	"strings"

	"github.com/issue9/logs/v3/config"
	"github.com/issue9/logs/v3/writers"
)

var flagMap = map[string]int{
	"none":              0,
	"log.ldate":         log.Ldate,
	"log.ltime":         log.Ltime,
	"log.lmicroseconds": log.Lmicroseconds,
	"log.llongfile":     log.Llongfile,
	"log.lshortfile":    log.Lshortfile,
	"log.lmsgprefix":    log.Lmsgprefix,
	"log.lstdflags":     log.LstdFlags,
}

// 扩展 log.Logger，使可以同时输出到多个日志通道
type logger struct {
	*log.Logger
	container *writers.Container
	level     int
}

func newLogger(level int, prefix string, flag int) *logger {
	c := writers.NewContainer()
	return &logger{
		Logger:    log.New(c, prefix, flag),
		container: c,
		level:     level,
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

// 该接口仅为兼容 toWriter 所使用。不应该直接调用。
//
// 当然如果直接调用该接口，也能将内容正确输出到日志。
func (l *logger) Write(data []byte) (int, error) {
	return l.container.Write(data)
}

// Add 可以让 toWriter 直接调用添加 io.Writer 实现
func (l *logger) Add(w io.Writer) error {
	return l.container.Add(w)
}

func loggerInitializer(level int) WriterInitializer {
	return func(cfg *config.Config) (io.Writer, error) {
		flag, err := parseFlag(cfg.Attrs["flag"])
		if err != nil {
			return nil, err
		}

		return newLogger(level, cfg.Attrs["prefix"], flag), nil
	}
}

// 将 log.Ldate|log.Ltime 的字符串转换成正确的值
func parseFlag(flagStr string) (int, error) {
	flagStr = strings.TrimSpace(flagStr)
	if len(flagStr) == 0 {
		return 0, nil
	}

	strs := strings.Split(flagStr, "|")
	ret := 0

	for _, str := range strs {
		str = strings.ToLower(strings.TrimSpace(str))
		flag, found := flagMap[str]
		if !found {
			return 0, errors.New("无效的 flag:" + str)
		}
		ret |= flag
	}

	return ret, nil
}
