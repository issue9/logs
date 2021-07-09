// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"testing"

	"github.com/issue9/assert"

	"github.com/issue9/logs/v3/config"
	"github.com/issue9/logs/v3/internal/initfunc"
)

type testLogs struct {
	logs                                     *Logs
	info, debug, erro, trace, warn, critical *bytes.Buffer
}

func newLogs(a *assert.Assertion) *testLogs {
	l, err := New(nil)
	a.NotError(err).NotNil(l)

	debug := new(bytes.Buffer)
	info := new(bytes.Buffer)
	erro := new(bytes.Buffer)
	trace := new(bytes.Buffer)
	warn := new(bytes.Buffer)
	critical := new(bytes.Buffer)

	a.NotError(l.SetOutput(LevelInfo, info))
	a.NotError(l.SetOutput(LevelDebug, debug))
	a.NotError(l.SetOutput(LevelError, erro))
	a.NotError(l.SetOutput(LevelTrace, trace))
	a.NotError(l.SetOutput(LevelWarn, warn))
	a.NotError(l.SetOutput(LevelCritical, critical))

	return &testLogs{
		logs:     l,
		info:     info,
		debug:    debug,
		erro:     erro,
		trace:    trace,
		warn:     warn,
		critical: critical,
	}
}

func (l *testLogs) checkLog(a *assert.Assertion) {
	a.True(l.info.Len() > 0)
	a.True(l.debug.Len() > 0)
	a.True(l.erro.Len() > 0)
	a.True(l.trace.Len() > 0)
	a.True(l.warn.Len() > 0)
	a.True(l.critical.Len() > 0)
}

func TestLogs_All(t *testing.T) {
	a := assert.New(t)

	l := newLogs(a)
	l.logs.All("abc")
	l.checkLog(a)
}

func TestLogs_Allf(t *testing.T) {
	a := assert.New(t)

	l := newLogs(a)
	l.logs.Allf("abc")
	l.checkLog(a)
}

func TestNew(t *testing.T) {
	a := assert.New(t)
	l, err := New(nil)
	a.NotError(err).NotNil(l)
	debugW := new(bytes.Buffer)
	warnW := new(bytes.Buffer)

	debugWInit := func(*config.Config) (io.Writer, error) {
		return debugW, nil
	}

	warnWInit := func(*config.Config) (io.Writer, error) {
		return warnW, nil
	}

	// 重新注册以下用到的 writer
	clearInitializer()
	a.True(Register("debug", loggerInitializer), "注册 debug 时失败")
	a.True(Register("warn", loggerInitializer), "注册 warn 时失败")
	a.True(Register("buffer", initfunc.Buffer), "注册 buffer 时失败")
	a.True(Register("debugW", debugWInit), "注册 debugW 时失败")
	a.True(Register("warnW", warnWInit), "注册 warnW 时失败")

	data := `
<?xml version="1.0" encoding="utf-8" ?>
<logs>
	<debug prefix="[DEBUG]">
		<buffer size="10">
			<debugW />
		</buffer>
	</debug>
	<warn prefix="[WARN]">
		<warnW />
	</warn>
</logs>
`
	debugW.Reset()

	cfg := &config.Config{}
	a.NotError(xml.Unmarshal([]byte(data), cfg))
	l, err = New(cfg)
	a.NotError(err).NotNil(l)

	l.Debug("abc")
	a.True(debugW.Len() == 0, "assert.True 失败，实际值为%d", debugW.Len()) // 缓存未达 10，依然为空
	l.Allf("def\n")
	a.True(debugW.Len() == 0, "assert.True 失败，实际值为%d", debugW.Len()) // 缓存未达 10，依然为空

	l.Warn("warn")
	a.Contains(warnW.String(), "warn")

	// 测试 Flush
	a.NotError(l.Flush())
	a.True(debugW.Len() > 0)
	a.True(warnW.Len() > 0)
}

func TestLogs_Logger(t *testing.T) {
	a := assert.New(t)

	l, err := New(nil)
	a.NotError(err).NotNil(l)

	a.NotNil(l.Logger(LevelCritical))
	a.NotNil(l.CRITICAL())
	a.Nil(l.Logger(LevelAll))
}

func TestLogs_SetOutput(t *testing.T) {
	a := assert.New(t)
	l, err := New(nil)
	a.NotError(err).NotNil(l)

	l.SetOutput(0, nil) // 无任何操作发生

	a.NotError(l.SetOutput(LevelError, &bytes.Buffer{}))
	a.Equal(l.logs(LevelError)[0].container.Len(), 1)
	a.NotError(l.SetOutput(LevelError, nil))
	a.Equal(l.logs(LevelError)[0].container.Len(), 0)
}

func TestLogs_SetFlags(t *testing.T) {
	a := assert.New(t)
	l, err := New(nil)
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
	a := assert.New(t)
	l, err := New(nil)
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

func TestLogs_Panicf(t *testing.T) {
	a := assert.New(t)

	l := newLogs(a)

	l.logs.Error("error")
	a.True(l.erro.Len() > 0)
	a.Equal(l.debug.Len(), 0)

	a.Panic(func() {
		l.logs.Panicf(LevelError, "panic")
	})

	a.True(l.info.Len() == 0)
	a.True(l.erro.Len() > 0)

	// panic

	l = newLogs(a)

	a.Panic(func() {
		l.logs.Panic(LevelError, "panic")
	})

	a.True(l.info.Len() == 0)
	a.Contains(l.erro.String(), "panic")
}
