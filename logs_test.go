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

var (
	debugW    = new(bytes.Buffer)
	infoW     = new(bytes.Buffer)
	errorW    = new(bytes.Buffer)
	traceW    = new(bytes.Buffer)
	warnW     = new(bytes.Buffer)
	criticalW = new(bytes.Buffer)
)

func initLogs(logs *Logs, a *assert.Assertion) {
	infoW.Reset()
	debugW.Reset()
	errorW.Reset()
	traceW.Reset()
	warnW.Reset()
	criticalW.Reset()

	a.True(infoW.Len() == 0)
	a.True(debugW.Len() == 0)
	a.True(errorW.Len() == 0)
	a.True(traceW.Len() == 0)
	a.True(warnW.Len() == 0)
	a.True(criticalW.Len() == 0)

	a.NotError(logs.SetOutput(LevelInfo, infoW))
	a.NotError(logs.SetOutput(LevelDebug, debugW))
	a.NotError(logs.SetOutput(LevelError, errorW))
	a.NotError(logs.SetOutput(LevelTrace, traceW))
	a.NotError(logs.SetOutput(LevelWarn, warnW))
	a.NotError(logs.SetOutput(LevelCritical, criticalW))
}

func checkLog(a *assert.Assertion) {
	a.True(infoW.Len() > 0)
	a.True(debugW.Len() > 0)
	a.True(errorW.Len() > 0)
	a.True(traceW.Len() > 0)
	a.True(warnW.Len() > 0)
	a.True(criticalW.Len() > 0)
}

func TestLogs_All(t *testing.T) {
	a := assert.New(t)

	l, err := New(nil)
	a.NotError(err).NotNil(l)

	initLogs(l, a)
	l.All("abc")
	checkLog(a)
}

func TestLogs_Allf(t *testing.T) {
	a := assert.New(t)

	l, err := New(nil)
	a.NotError(err).NotNil(l)

	initLogs(l, a)
	l.Allf("abc")
	checkLog(a)
}

func debugWInit(*config.Config) (io.Writer, error) {
	return debugW, nil
}

func TestNew(t *testing.T) {
	a := assert.New(t)
	l, err := New(nil)
	a.NotError(err).NotNil(l)

	// 重新注册以下用到的 writer
	clearInitializer()
	a.True(Register("debug", loggerInitializer(LevelDebug)), "注册 debug 时失败")
	a.True(Register("buffer", initfunc.Buffer), "注册 buffer 时失败")
	a.True(Register("debugW", debugWInit), "注册 debugW 时失败")

	data := `
<?xml version="1.0" encoding="utf-8" ?>
<logs>
	<debug prefix="[DEBUG]">
		<buffer size="10">
			<debugW />
		</buffer>
	</debug>
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

	// 测试 Flush
	a.NotError(l.Flush())
	a.True(debugW.Len() > 0)
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
	l, err := New(nil)
	a.NotError(err).NotNil(l)

	initLogs(l, a)

	l.Error("error")
	a.True(errorW.Len() > 0)
	a.Equal(debugW.Len(), 0)

	a.Panic(func() {
		l.Panicf(LevelError, "panic")
	})

	a.True(infoW.Len() == 0)
	a.True(errorW.Len() > 0)
}
