// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestLogs(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout))
	a.NotNil(l)
	l.Enable(LevelInfo, LevelWarn, LevelDebug, LevelTrace, LevelError, LevelFatal)

	testLogger := func(a *assert.Assertion, p func(...any), pf func(string, ...any), w *bytes.Buffer) {
		a.TB().Helper()

		p("p1")
		val := w.String()
		a.Contains(val, "p1").
			Contains(val, "logs_test.go:24") // 行数是否正确

		pf("p2")
		val = w.String()
		a.Contains(val, "p2").
			Contains(val, "logs_test.go:29") // 行数是否正确
	}

	testLogger(a, l.INFO().Print, l.INFO().Printf, buf)
	testLogger(a, l.DEBUG().Print, l.DEBUG().Printf, buf)
	testLogger(a, l.TRACE().Print, l.TRACE().Printf, buf)
	testLogger(a, l.WARN().Print, l.WARN().Printf, buf)
	testLogger(a, l.ERROR().Print, l.ERROR().Printf, buf)
	testLogger(a, l.FATAL().Print, l.FATAL().Printf, buf)
}

func TestLogs_AppendAttrs(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf))
	a.NotNil(l)

	l.AppendAttrs(map[string]any{"a1": 1})
	l.INFO().Print("123")
	a.Equal(buf.String(), "[INFO] 123 a1=1\n")
}

func TestLogs_New(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf))
	a.NotNil(l)

	al := l.New(map[string]any{"a1": 1})
	a.NotNil(al)

	a.Equal(al.INFO(), al.INFO()) // 确保不会每次都构建新对象

	al.ERROR().String("abc")
	a.Equal(buf.String(), "[ERRO] abc a1=1\n")

	FreeAttrLogs(al)
}

func TestAttrLogs_AppendAttrs(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf))
	a.NotNil(l)

	al := l.New(map[string]any{"a1": 1})
	a.NotNil(al)

	al.TRACE().String("abc")
	a.Equal(buf.String(), "[TRAC] abc a1=1\n")

	buf.Reset()
	al.AppendAttrs(map[string]any{"a2": 2})
	al.TRACE().String("abc")
	a.Equal(buf.String(), "[TRAC] abc a1=1 a2=2\n")
}

func TestLogs_IsEnable(t *testing.T) {
	a := assert.New(t, false)

	l := New(NewJSONHandler(os.Stdout))
	a.NotNil(l)
	l.Enable(LevelInfo, LevelWarn, LevelDebug, LevelTrace, LevelError, LevelFatal)
	a.True(l.IsEnable(LevelInfo)).
		True(l.IsEnable(LevelFatal)).
		True(l.IsEnable(LevelWarn))

	l.Enable(LevelWarn, LevelError)
	a.False(l.IsEnable(LevelInfo)).
		False(l.IsEnable(LevelFatal)).
		True(l.IsEnable(LevelWarn)).
		True(l.IsEnable(LevelError))

	buf := new(bytes.Buffer)
	l = New(NewTextHandler(buf))
	a.NotNil(l)
	l.Enable(LevelWarn, LevelError)

	a.NotEqual(l.WARN(), disabledRecorder)

	a.False(l.FATAL().IsEnable())

	// enable=false，disabledLogger.With
	buf.Reset()
	inst := l.FATAL()
	a.Equal(inst.With("k1", "v1").With("k2", "v2"), disabledRecorder)
	inst.With("k2", "v2").Error(errors.New("err"))
	a.Zero(buf.Len())

	// 运行过程中调整了 Level 的值
	l.Enable(LevelFatal)
	inst = l.FATAL()
	a.NotEqual(inst.With("k1", "v1"), disabledRecorder) // k1=v1 并未保存
	inst.With("k2", "v2").Error(errors.New("err"))
	a.NotContains(buf.String(), "k1=v1").
		Contains(buf.String(), "k2=v2").
		Contains(buf.String(), "err")
}
