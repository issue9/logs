// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"errors"
	"sync"
	"testing"

	"github.com/issue9/assert/v3"
)

var (
	_ Logger = &logger{}
	_ Logger = &withLogger{}
	_ Logger = &emptyLogger{}
)

func TestLogger_location(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextWriter(MicroLayout, buf), Caller, Created)
	a.NotNil(l)
	l.Enable(LevelError)

	// Record.Location
	l.ERROR().With("k1", "v1").
		Printf("Record.Printf") // 位置记录此行
	val := buf.String()
	a.Contains(val, "logger_test.go:29").
		Contains(val, "k1=v1").
		Contains(val, "Record.Printf")

	// Logs.Location
	buf.Reset()
	l.Errorf("Logs.%s", "Errorf")
	val = buf.String()
	a.Contains(val, "logger_test.go:37").
		Contains(val, "Logs.Errorf")

	// logger.Location
	buf.Reset()
	l.ERROR().Print("logger.Print")
	val = buf.String()
	a.Contains(val, "logger_test.go:44").
		Contains(val, "logger.Print")

	buf.Reset()
	l.SetCaller(false)
	l.ERROR().Printf("caller=false")
	val = buf.String()
	a.NotContains(val, "logger_test.go").
		Contains(val, "caller=false")
}

func TestLogger_Error(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextWriter(MicroLayout, buf))
	a.NotNil(l)

	// Record.Error
	l.ERROR().With("k1", "v1").
		Error(errors.New("err"))
	val := buf.String()
	a.Contains(val, "err").
		Contains(val, "k1=v1")

	// logger.Error
	buf.Reset()
	l.DEBUG().Error(errors.New("info"))
	a.Contains(buf.String(), "info")
}

func TestLogger_With(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextWriter(MicroLayout, buf), Caller)
	a.NotNil(l)

	err := l.With(LevelError, map[string]any{"k1": "v1"})
	err.Printf("err1")
	a.Contains(buf.String(), "err1").
		Contains(buf.String(), "k1=v1").
		Contains(buf.String(), "logger_test.go:83")

	buf.Reset()
	err.With("k2", "v2").Printf("err2")
	a.Contains(buf.String(), "err2").
		Contains(buf.String(), "k1=v1").
		Contains(buf.String(), "k2=v2").
		NotContains(buf.String(), "err1")

	buf.Reset()
	err.With("k3", "v3").Print("err3")
	a.Contains(buf.String(), "err3").
		Contains(buf.String(), "k1=v1").
		Contains(buf.String(), "k3=v3").
		NotContains(buf.String(), "err1").
		NotContains(buf.String(), "k2=v2").
		NotContains(buf.String(), "err2")

	buf.Reset()
	err.Error(errors.New("err1"))
	a.Contains(buf.String(), "err1").
		Contains(buf.String(), "k1=v1")

	buf.Reset()
	l.Enable(LevelDebug)
	err = l.With(LevelError, map[string]any{"k2": "v2"})
	err.Println("err")
	a.Equal(err, emptyLLoggerInst).
		NotNil(err.StdLogger()).
		Empty(buf.String())
}

func TestLogger_String(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextWriter(MicroLayout, buf))
	a.NotNil(l)

	// Record.String
	l.ERROR().With("k1", "v1").
		String("string")
	val := buf.String()
	a.Contains(val, "string").
		Contains(val, "k1=v1")

	// logger.String
	buf.Reset()
	l.DEBUG().String("info")
	a.Contains(buf.String(), "info")
}

func TestLogger_StdLogger(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	w := NewTextWriter(MicroLayout, buf)
	l := New(w, Created, Caller)
	a.NotNil(l)
	l.Enable(LevelInfo, LevelError)

	// logger.StdLogger

	info := l.INFO().StdLogger()
	info.Print("abc")
	a.Contains(buf.String(), "logger_test.go:148") // 行数是否正确

	// Enable 未设置 LevelWarn
	buf.Reset()
	warn := l.WARN().StdLogger()
	warn.Print("abc")
	a.Equal(buf.Len(), 0)

	// withLogger.StdLogger

	buf.Reset()
	err := l.With(LevelError, map[string]any{"k1": "v1"}).StdLogger()
	err.Print("abc")
	a.Contains(buf.String(), "logger_test.go:161"). // 行数是否正确
							Contains(buf.String(), "k1=v1")
}

func TestLogger_Printf(t *testing.T) {
	a := assert.New(t, false)

	buf := new(bytes.Buffer)
	l := New(NewTextWriter(MicroLayout, buf))
	a.NotNil(l)
	l.Enable(LevelError)

	g := sync.WaitGroup{}
	err := l.ERROR()

	g.Add(1)
	go func() {
		err.Printf("这是一段不可分割的文字内容 1")
		g.Done()
	}()

	g.Add(1)
	go func() {
		err.Printf("这是一段不可分割的文字内容 2")
		g.Done()
	}()

	g.Add(1)
	go func() {
		err.Printf("这是一段不可分割的文字内容 3")
		g.Done()
	}()

	g.Wait()

	a.Contains(buf.String(), "这是一段不可分割的文字内容 1")
	a.Contains(buf.String(), "这是一段不可分割的文字内容 2")
	a.Contains(buf.String(), "这是一段不可分割的文字内容 3")
}
