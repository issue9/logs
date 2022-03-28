// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
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
	a.Equal(l.WARN(), emptyLoggerInst)

	l = New(NewTextWriter("2006", new(bytes.Buffer)))
	a.NotNil(l)
	l.Enable(LevelWarn, LevelError)

	a.Equal(l.FATAL(), emptyLoggerInst)

	_, ok := l.WARN().(*logger)
	a.True(ok)
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
	l := New(w)
	a.NotNil(l)
	l.Enable(LevelInfo, LevelError)

	info := l.StdLogger(LevelInfo)
	info.Print("abc")
	a.Contains(buf.String(), "logs_test.go:76")

	// Enable 未设置 LevelWarn
	buf.Reset()
	warn := l.StdLogger(LevelWarn)
	warn.Print("abc")
	a.Equal(buf.Len(), 0)
}
