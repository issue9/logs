// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"testing"

	"github.com/issue9/assert/v2"
)

func TestLogs_IsEnable(t *testing.T) {
	a := assert.New(t, false)

	l := New()
	a.NotNil(l)
	a.True(l.IsEnable(LevelInfo)).
		True(l.IsEnable(LevelFatal)).
		True(l.IsEnable(LevelWarn))

	l.Enable(LevelWarn, LevelError)
	a.False(l.IsEnable(LevelInfo)).
		False(l.IsEnable(LevelFatal)).
		True(l.IsEnable(LevelWarn)).
		True(l.IsEnable(LevelError))

	ll := l.WARN()
	_, ok := ll.(*logger)
	a.True(ok)

	ll = l.FATAL()
	a.Equal(ll, emptyLoggerInst)
}

func TestLogsLoggers(t *testing.T) {
	a := assert.New(t, false)
	l := New()
	a.NotNil(l)
	buf := new(bytes.Buffer)
	w := NewWriter(TextFormat("2006-01-02"), buf)
	l.SetOutput(w, LevelInfo, LevelWarn, LevelDebug, LevelTrace, LevelError, LevelFatal)

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
	l := New()
	a.NotNil(l)
	buf := new(bytes.Buffer)
	w := NewWriter(TextFormat("2006-01-02"), buf)
	l.SetOutput(w, LevelInfo, LevelError)

	info := l.StdLogger(LevelInfo)
	info.Print("abc")
	val := buf.String()
	a.Contains(val, "logs_test.go:70")

	// SetOutput 未设置 LevelWarn
	buf.Reset()
	warn := l.StdLogger(LevelWarn)
	warn.Print("abc")
	val = buf.String()
	a.Empty(val)
}
