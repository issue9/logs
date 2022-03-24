// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"log"
	"testing"

	"github.com/issue9/assert/v2"
)

type testLogs struct {
	logs    *Logs
	buffers map[int]*bytes.Buffer
	a       *assert.Assertion
}

func newLogs(a *assert.Assertion) *testLogs {
	l, err := New()
	a.NotError(err).NotNil(l)

	bufs := map[int]*bytes.Buffer{
		LevelInfo:     new(bytes.Buffer),
		LevelTrace:    new(bytes.Buffer),
		LevelDebug:    new(bytes.Buffer),
		LevelWarn:     new(bytes.Buffer),
		LevelError:    new(bytes.Buffer),
		LevelCritical: new(bytes.Buffer),
	}

	a.NotError(l.SetOutput(LevelInfo, bufs[LevelInfo]))
	a.NotError(l.SetOutput(LevelDebug, bufs[LevelDebug]))
	a.NotError(l.SetOutput(LevelError, bufs[LevelError]))
	a.NotError(l.SetOutput(LevelTrace, bufs[LevelTrace]))
	a.NotError(l.SetOutput(LevelWarn, bufs[LevelWarn]))
	a.NotError(l.SetOutput(LevelCritical, bufs[LevelCritical]))

	return &testLogs{
		logs:    l,
		buffers: bufs,
		a:       a,
	}
}

func (l *testLogs) chk(level int, content string) {
	l.a.TB().Helper()
	for lv, buf := range l.buffers {
		if lv&level == lv {
			l.a.Contains(buf.String(), content)
		}
	}
}

func (l *testLogs) reset() {
	for _, b := range l.buffers {
		b.Reset()
	}
}

func TestLogs_All(t *testing.T) {
	a := assert.New(t, false)

	l := newLogs(a)
	l.logs.All("abc")
	l.chk(LevelAll, "abc")
}

func TestLogs_Allf(t *testing.T) {
	a := assert.New(t, false)

	l := newLogs(a)
	l.logs.Allf("abc")
	l.chk(LevelAll, "abc")
}

func TestLogs_Logger(t *testing.T) {
	a := assert.New(t, false)

	l, err := New()
	a.NotError(err).NotNil(l)

	a.NotNil(l.Logger(LevelCritical))
	a.NotNil(l.CRITICAL())
	a.Nil(l.Logger(LevelAll))
}

func TestLogs_SetOutput(t *testing.T) {
	a := assert.New(t, false)
	l, err := New()
	a.NotError(err).NotNil(l)

	a.NotError(l.SetOutput(0, nil)) // 无任何操作发生

	a.NotError(l.SetOutput(LevelError, &bytes.Buffer{}))
	a.Equal(l.loggers[LevelError].container.Len(), 1)
	a.NotError(l.SetOutput(LevelError, nil))
	a.Equal(l.loggers[LevelError].container.Len(), 0)
}

func TestLogs_SetFlags(t *testing.T) {
	a := assert.New(t, false)
	l, err := New()
	a.NotError(err).NotNil(l)

	l.SetFlags(LevelAll, log.Ldate)
	for _, item := range l.loggers {
		a.Equal(item.Flags(), log.Ldate)
	}

	l.SetFlags(LevelAll, log.Lmsgprefix)
	for _, item := range l.loggers {
		a.Equal(item.Flags(), log.Lmsgprefix)
	}
}

func TestLogs_SetPrefix(t *testing.T) {
	a := assert.New(t, false)
	l, err := New()
	a.NotError(err).NotNil(l)

	l.SetPrefix(LevelAll, "p")
	for _, item := range l.loggers {
		a.Equal(item.Prefix(), "p")
	}

	l.SetPrefix(LevelAll, "")
	for _, item := range l.loggers {
		a.Equal(item.Prefix(), "")
	}
}

// NOTE:以下内容依赖所在的文件以及行号，有变动需要及时更改。
func TestLogs_LN(t *testing.T) {
	a := assert.New(t, false)
	l := newLogs(a)
	l.logs.SetFlags(LevelAll, log.Llongfile)

	l.logs.Trace("abc")
	l.chk(LevelTrace, "logs_test.go:139")
	l.reset()

	l.logs.Printf(LevelDebug, 0, "abc")
	l.chk(LevelDebug, "logs_test.go:143")
	l.reset()

	l.logs.Print(LevelDebug, 0, "abc")
	l.chk(LevelDebug, "logs_test.go:147")
	l.reset()

	a.Panic(func() {
		l.logs.Panicf(LevelError, "panic")
	})
	l.chk(LevelError, "logs_test.go:152")
	l.reset()

	l.logs.All("abc")
	l.chk(LevelDebug, "logs_test.go:157")
	l.reset()
}

func TestLogs_Panicf(t *testing.T) {
	a := assert.New(t, false)
	l := newLogs(a)

	a.Panic(func() {
		l.logs.Panicf(LevelError, "panic")
	})
	l.chk(LevelError, "panic")
	a.True(l.buffers[LevelWarn].Len() == 0)

	// panic

	l = newLogs(a)

	a.Panic(func() {
		l.logs.Panic(LevelCritical, "panic")
	})
	l.chk(LevelCritical, "panic")
	a.True(l.buffers[LevelError].Len() == 0)
}
