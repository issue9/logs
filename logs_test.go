// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"errors"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestLogsLoggers(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	w := NewTextHandler(MicroLayout, buf)
	l := New(w, Caller, Created)
	a.NotNil(l)
	l.Enable(LevelInfo, LevelWarn, LevelDebug, LevelTrace, LevelError, LevelFatal)

	testLogger := func(a *assert.Assertion, p func(...any), pf func(string, ...any), w *bytes.Buffer) {
		p("p1")
		val := w.String()
		a.Contains(val, "p1").
			Contains(val, "logs_test.go:22") // 行数是否正确

		pf("p2")
		val = w.String()
		a.Contains(val, "p2").
			Contains(val, "logs_test.go:27") // 行数是否正确
	}

	testLogger(a, l.INFO().Print, l.INFO().Printf, buf)
	testLogger(a, l.DEBUG().Print, l.DEBUG().Printf, buf)
	testLogger(a, l.TRACE().Print, l.TRACE().Printf, buf)
	testLogger(a, l.WARN().Print, l.WARN().Printf, buf)
	testLogger(a, l.ERROR().Print, l.ERROR().Printf, buf)
	testLogger(a, l.FATAL().Print, l.FATAL().Printf, buf)
}

func TestLogs_IsEnable(t *testing.T) {
	a := assert.New(t, false)

	l := New(nil)
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

	// WARN 属于 enable，但是 logs.w 为 nop
	ll, ok := l.WARN().(*logger)
	a.True(ok).False(ll.enable)

	buf := new(bytes.Buffer)
	l = New(NewTextHandler(MicroLayout, buf))
	a.NotNil(l)
	l.Enable(LevelWarn, LevelError)

	ll, ok = l.WARN().(*logger)
	a.True(ok).True(ll.enable)

	ll, ok = l.FATAL().(*logger)
	a.True(ok).False(ll.enable)

	// enable=false，disabledLogger.With
	buf.Reset()
	inst := l.FATAL()
	a.Equal(inst.With("k1", "v1").With("k2", "v2"), disabledLogger)
	inst.With("k2", "v2").Error(errors.New("err"))
	a.Zero(buf.Len())

	// 运行过程中调整了 Level 的值
	l.Enable(LevelFatal)
	inst = l.FATAL()
	a.NotEqual(inst.With("k1", "v1"), disabledLogger) // k1=v1 并未保存
	inst.With("k2", "v2").Error(errors.New("err"))
	a.NotContains(buf.String(), "k1=v1").
		Contains(buf.String(), "k2=v2").
		Contains(buf.String(), "err")
}
