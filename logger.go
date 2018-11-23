// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/issue9/logs/writers"
)

var flagMap = map[string]int{
	"none":              0,
	"log.ldate":         log.Ldate,
	"log.ltime":         log.Ltime,
	"log.lmicroseconds": log.Lmicroseconds,
	"log.llongfile":     log.Llongfile,
	"log.lshortfile":    log.Lshortfile,
	"log.lstdflags":     log.LstdFlags,
}

// 扩展 log.Logger，使可以同时输出到多个日志通道
type logger struct {
	flush writers.Flusher // 如果当前的 log 的 io.Writer 实例是个容器，则此处保存此容器的指针。
	log   *log.Logger     // 要确保这些值不能为空，因为要保证对应的 ERROR() 等函数的返回值是始终可用的。
}

func (l *logger) set(w io.Writer, prefix string, flag int) {
	if w == nil {
		l.flush = nil
		l.log.SetOutput(ioutil.Discard)
		return
	}

	l.log.SetFlags(flag)
	l.log.SetPrefix(prefix)
	l.log.SetOutput(w)
	if f, ok := w.(writers.Flusher); ok {
		l.flush = f
	}
}

func loggerInitializer(args map[string]string) (io.Writer, error) {
	return writers.NewContainer(), nil
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
