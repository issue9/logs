// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"errors"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestLogger_location(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout))
	a.NotNil(l)
	l.Enable(LevelError)

	// Record.Location
	l.ERROR().With("k1", "v1").
		Printf("Record.Printf") // 位置记录此行
	val := buf.String()
	a.Contains(val, "logger_test.go:22").
		Contains(val, "k1=v1").
		Contains(val, "Record.Printf")

	// Logs.Location
	buf.Reset()
	l.ERROR().Printf("Logs.%s", "Errorf")
	val = buf.String()
	a.Contains(val, "logger_test.go:30").
		Contains(val, "Logs.Errorf")

	// logger.Location
	buf.Reset()
	l.ERROR().Print("logger.Print")
	val = buf.String()
	a.Contains(val, "logger_test.go:37").
		Contains(val, "logger.Print")

	buf.Reset()
	l.SetLocation(false)
	l.ERROR().Printf("caller=false")
	val = buf.String()
	a.NotContains(val, "logger_test.go").
		Contains(val, "caller=false")
}

func TestLogger_Error(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf))
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

func TestLogger_New(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf), WithLocation(true), WithAttrs(map[string]any{"a1": "v1"}))
	a.NotNil(l)

	err := l.ERROR().New(map[string]any{"k1": "v1"})
	err.Printf("err1")
	a.Contains(buf.String(), "err1").
		Contains(buf.String(), "k1=v1").
		Contains(buf.String(), "a1=v1").
		Contains(buf.String(), "logger_test.go:76")

	buf.Reset()
	err.With("k2", "v2").Printf("err2")
	a.Contains(buf.String(), "err2").
		Contains(buf.String(), "k1=v1").
		Contains(buf.String(), "a1=v1").
		Contains(buf.String(), "k2=v2").
		NotContains(buf.String(), "err1")

	buf.Reset()
	err.With("k3", "v3").Print("err3")
	a.Contains(buf.String(), "err3").
		Contains(buf.String(), "k1=v1").
		Contains(buf.String(), "a1=v1").
		Contains(buf.String(), "k3=v3").
		NotContains(buf.String(), "err1").
		NotContains(buf.String(), "k2=v2").
		NotContains(buf.String(), "err2")

	buf.Reset()
	err.Error(errors.New("err1"))
	a.Contains(buf.String(), "err1").
		Contains(buf.String(), "a1=v1").
		Contains(buf.String(), "k1=v1")

	buf.Reset()
	l.Enable(LevelDebug)
	err = l.ERROR().New(map[string]any{"k2": "v2"})
	r := err.With("k", "v")
	a.False(l.ERROR().IsEnable()).
		False(err.IsEnable()).
		Equal(r, disabledRecorder).
		Empty(buf.String())
}

func TestLogger_String(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf))
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

func TestLogger_LogLogger(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	w := NewTextHandler(buf)
	l := New(w, WithCreated(MicroLayout), WithLocation(true))
	a.NotNil(l)
	l.Enable(LevelInfo, LevelError)

	// logger.LogLogger

	info := l.INFO().LogLogger()
	info.Print("abc")
	a.Contains(buf.String(), "logger_test.go:146") // 行数是否正确

	// Enable 未设置 LevelWarn
	buf.Reset()
	warn := l.WARN().LogLogger()
	warn.Print("abc")
	a.Equal(buf.Len(), 0)
}
