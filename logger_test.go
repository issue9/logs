// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/issue9/assert/v3"
)

var (
	_ Logger    = &logger{}
	_ io.Writer = &logger{}

	_ Logger = &emptyLogger{}

	_ Logger = &Entry{}
)

func TestEntry_Location(t *testing.T) {
	a := assert.New(t, false)
	l := New(nil, Caller, Created)

	e := l.NewEntry(LevelWarn)
	a.NotNil(e)
	a.Empty(e.Path).Zero(e.Line)

	e.Location(1)
	a.True(strings.HasSuffix(e.Path, "logger_test.go")).Equal(e.Line, 33)
}

func TestLogger_location(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextWriter(MicroLayout, buf), Caller, Created)
	a.NotNil(l)
	l.Enable(LevelError)

	// Entry.Location
	l.ERROR().With("k1", "v1").
		Printf("Entry.Printf") // 位置记录此行
	val := buf.String()
	a.Contains(val, "logger_test.go:46").
		Contains(val, "k1=v1").
		Contains(val, "Entry.Printf")

	// Logs.Location
	buf.Reset()
	l.Errorf("Logs.%s", "Errorf")
	val = buf.String()
	a.Contains(val, "logger_test.go:54").
		Contains(val, "Logs.Errorf")

	// logger.Location
	buf.Reset()
	l.ERROR().Print("logger.Print")
	val = buf.String()
	a.Contains(val, "logger_test.go:61").
		Contains(val, "logger.Print")
}

func TestLogger_Error(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextWriter(MicroLayout, buf))
	a.NotNil(l)

	// Entry.Error
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

func TestLogger_String(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextWriter(MicroLayout, buf))
	a.NotNil(l)

	// Entry.String
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
