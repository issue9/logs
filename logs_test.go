// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"errors"
	"testing"

	"github.com/issue9/assert/v2"
)

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

	// WARN 属于 enable，但是没有 logs.w 为 Nop
	ll, ok := l.WARN().(*logger)
	a.True(ok).False(ll.enable)

	buf := new(bytes.Buffer)
	l = New(NewTextWriter("2006", buf))
	a.NotNil(l)
	l.Enable(LevelWarn, LevelError)

	ll, ok = l.WARN().(*logger)
	a.True(ok).True(ll.enable)

	ll, ok = l.FATAL().(*logger)
	a.True(ok).False(ll.enable)

	// enable=false，Value emptyLoggerInst
	buf.Reset()
	inst := l.FATAL().Value("k1", "v1")
	a.Equal(inst, emptyLoggerInst)
	inst.Value("k2", "v2").Error(errors.New("err"))
	a.Zero(buf.Len())
}

func TestLogsLoggers(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	w := NewTextWriter("2006-01-02", buf)
	l := New(w)
	a.NotNil(l)
	l.Enable(LevelInfo, LevelWarn, LevelDebug, LevelTrace, LevelError, LevelFatal)

	testLogger := func(a *assert.Assertion, p func(...interface{}), pf func(string, ...interface{}), w *bytes.Buffer) {
		p("p1")
		val := w.String()
		a.Contains(val, "p1")

		pf("p2")
		val = w.String()
		a.Contains(val, "p2")
	}

	testLogger(a, l.Info, l.Infof, buf)
	testLogger(a, l.Debug, l.Debugf, buf)
	testLogger(a, l.Trace, l.Tracef, buf)
	testLogger(a, l.Warn, l.Warnf, buf)
	testLogger(a, l.Error, l.Errorf, buf)
	testLogger(a, l.Fatal, l.Fatalf, buf)
}

func TestLogs_StdLogger(t *testing.T) {
	a := assert.New(t, false)
	buf := new(bytes.Buffer)
	w := NewTextWriter("2006-01-02", buf)
	l := New(w, Created, Caller)
	a.NotNil(l)
	l.Enable(LevelInfo, LevelError)

	info := l.StdLogger(LevelInfo)
	info.Print("abc")
	a.Contains(buf.String(), "logs_test.go:87")

	// Enable 未设置 LevelWarn
	buf.Reset()
	warn := l.StdLogger(LevelWarn)
	warn.Print("abc")
	a.Equal(buf.Len(), 0)
}
