// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/issue9/logs/v6/writers"
)

const poolMaxParams = 100

var recordPool = &sync.Pool{New: func() any { return &Record{} }}

type (
	// Record 每一条日志的表示
	Record struct {
		logs *Logs

		Level   Level
		Created string // 日志的生成时间

		// 日志的实际内容
		//
		// 如果要改变此值，请使用 Depth* 系列的方法。
		Message string

		// 以下表示日志的定位信息
		Path string

		// 额外的数据保存在此，比如由 [Logger.With] 添加的数据。
		Params []Pair
	}

	Pair struct {
		K string
		V any
	}
)

func (logs *Logs) NewRecord(lv Level) *Record {
	e := recordPool.Get().(*Record)

	e.logs = logs
	if e.Params != nil {
		e.Params = e.Params[:0]
	}
	e.Path = ""
	e.Message = ""
	if logs.createdFormat != "" {
		e.Created = time.Now().Format(logs.createdFormat)
	} else {
		e.Created = "" // 从 pool 中获取的值，必须要初始化。
	}
	e.Level = lv

	return e
}

// 转换成 io.Writer
//
// 仅供内部使用，因为 depth 值的关系，只有固定的调用层级关系才能正常显示行号。
func (e *Record) asWriter() io.Writer {
	return writers.WriteFunc(func(data []byte) (int, error) {
		e.DepthString(5, string(data))
		return len(data), nil
	})
}

func (e *Record) Logs() *Logs { return e.logs }

// depth 表示调用，1 表示调用 Location 的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么此调用将不产生任何内容。
func (e *Record) setLocation(depth int) *Record {
	if e.Logs().HasCaller() {
		_, p, l, _ := runtime.Caller(depth)
		e.Path = p + ":" + strconv.Itoa(l)
	}
	return e
}

func (e *Record) With(name string, val any) Logger {
	e.Params = append(e.Params, Pair{K: name, V: val})
	return e
}

func (e *Record) StdLogger() *log.Logger {
	return log.New(e.asWriter(), "", 0)
}

func (e *Record) Error(err error) { e.DepthError(2, err) }

// DepthError 输出 error 类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthError(depth int, err error) {
	if err != nil {
		e.Message = err.Error()
	}
	e.setLocation(depth + 1)
	e.output()
}

func (e *Record) String(s string) { e.DepthString(2, s) }

// DepthString 输出字符串类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthString(depth int, s string) {
	e.Message = s
	e.setLocation(depth + 1)
	e.output()
}

func (e *Record) Print(v ...any) { e.DepthPrint(2, v...) }

// DepthPrint 输出任意类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthPrint(depth int, v ...any) {
	if len(v) > 0 {
		e.Message = fmt.Sprint(v...)
	}
	e.setLocation(depth + 1)
	e.output()
}

func (e *Record) Printf(format string, v ...any) { e.DepthPrintf(2, format, v...) }

// DepthPrintf 输出任意类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthPrintf(depth int, format string, v ...any) {
	e.Message = fmt.Sprintf(format, v...)
	e.setLocation(depth + 1)
	e.output()
}

func (e *Record) Println(v ...any) { e.DepthPrintln(2, v...) }

// DepthPrintln 输出任意类型的内容到日志
//
// depth 表示调用，1 表示调用此方法的位置；
//
// 如果 [Logs.HasCaller] 为 false，那么 depth 将不起实际作用。
func (e *Record) DepthPrintln(depth int, v ...any) {
	if len(v) > 0 {
		e.Message = fmt.Sprintln(v...)
	}
	e.setLocation(depth + 1)
	e.output()
}

func (e *Record) output() {
	e.logs.handler.Handle(e)
	if len(e.Params) < poolMaxParams {
		recordPool.Put(e)
	}
}
