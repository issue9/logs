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
	logs    *Logs
	buffers map[int]*bytes.Buffer
	a       *assert.Assertion
}

func newLogs(a *assert.Assertion) *testLogs {
	l, err := New(nil)
	a.NotError(err).NotNil(l)

	bufs := make(map[int]*bytes.Buffer)
	for _, l := range levels {
		bufs[l] = new(bytes.Buffer)
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
	a := assert.New(t)

	l := newLogs(a)
	l.logs.All("abc")
	l.chk(LevelAll, "abc")
}

func TestLogs_Allf(t *testing.T) {
	a := assert.New(t)

	l := newLogs(a)
	l.logs.Allf("abc")
	l.chk(LevelAll, "abc")
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

	a.NotError(l.SetOutput(0, nil)) // 无任何操作发生

	a.NotError(l.SetOutput(LevelError, &bytes.Buffer{}))
	a.Equal(l.loggers[LevelError].container.Len(), 1)
	a.NotError(l.SetOutput(LevelError, nil))
	a.Equal(l.loggers[LevelError].container.Len(), 0)
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

// NOTE:以下内容依赖所在的文件以及行号，有变动需要及时更改。
func TestLogs_LN(t *testing.T) {
	a := assert.New(t)
	l := newLogs(a)
	l.logs.SetFlags(LevelAll, log.Llongfile)

	l.logs.Trace("abc")
	l.chk(LevelTrace, "logs_test.go:196")
	l.reset()

	l.logs.Printf(LevelDebug, "abc")
	l.chk(LevelDebug, "logs_test.go:200")
	l.reset()

	l.logs.Print(LevelDebug, "abc")
	l.chk(LevelDebug, "logs_test.go:204")
	l.reset()

	a.Panic(func() {
		l.logs.Panicf(LevelError, "panic")
	})
	l.chk(LevelError, "logs_test.go:209")
	l.reset()

	l.logs.All("abc")
	l.chk(LevelDebug, "logs_test.go:214")
	l.reset()
}

func TestLogs_Panicf(t *testing.T) {
	a := assert.New(t)
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
