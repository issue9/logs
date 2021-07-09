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

func resetLog(logs *Logs, t *testing.T) {
	a := assert.New(t)

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

	a.NotError(logs.loggers[LevelInfo].SetOutput(infoW))
	a.NotError(logs.loggers[LevelDebug].SetOutput(debugW))
	a.NotError(logs.loggers[LevelError].SetOutput(errorW))
	a.NotError(logs.loggers[LevelTrace].SetOutput(traceW))
	a.NotError(logs.loggers[LevelWarn].SetOutput(warnW))
	a.NotError(logs.loggers[LevelCritical].SetOutput(criticalW))
}

func checkLog(t *testing.T) {
	a := assert.New(t)

	a.True(infoW.Len() > 0)
	a.True(debugW.Len() > 0)
	a.True(errorW.Len() > 0)
	a.True(traceW.Len() > 0)
	a.True(warnW.Len() > 0)
	a.True(criticalW.Len() > 0)
}

func TestAll(t *testing.T) {
	resetLog(defaultLogs, t)
	All("abc")
	checkLog(t)
}

func TestAllf(t *testing.T) {
	resetLog(defaultLogs, t)
	Allf("abc")
	checkLog(t)
}

func debugWInit(*config.Config) (io.Writer, error) {
	return debugW, nil
}

func TestInit(t *testing.T) {
	a := assert.New(t)

	// 重新注册以下用到的 writer
	clearInitializer()
	a.True(Register("debug", loggerInitializer), "注册 debug 时失败")
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
	a.NotError(Init(cfg))

	Debug("abc")
	a.True(debugW.Len() == 0, "assert.True 失败，实际值为%d", debugW.Len()) // 缓存未达 10，依然为空
	Allf("def\n")
	a.True(debugW.Len() == 0, "assert.True 失败，实际值为%d", debugW.Len()) // 缓存未达 10，依然为空

	// 测试 Flush
	a.NotError(Flush())
	a.True(debugW.Len() > 0)
}

func TestLogs_SetOutput(t *testing.T) {
	a := assert.New(t)

	a.Panic(func() {
		Default().SetOutput(-1, nil)
	})

	a.NotError(Default().SetOutput(LevelError, &bytes.Buffer{}))
	a.Equal(Default().loggers[LevelError].container.Len(), 1)
	a.NotError(Default().SetOutput(LevelError, nil))
	a.Equal(Default().loggers[LevelError].container.Len(), 0)
}

func TestLogs_SetFlags(t *testing.T) {
	a := assert.New(t)

	Default().SetFlags(log.Ldate)
	for _, l := range Default().loggers {
		a.Equal(l.Flags(), log.Ldate)
	}

	Default().SetFlags(log.Lmsgprefix)
	for _, l := range Default().loggers {
		a.Equal(l.Flags(), log.Lmsgprefix)
	}
}

func TestLogs_SetPrefix(t *testing.T) {
	a := assert.New(t)

	Default().SetPrefix("p")
	for _, l := range Default().loggers {
		a.Equal(l.Prefix(), "p")
	}

	Default().SetPrefix("")
	for _, l := range Default().loggers {
		a.Equal(l.Prefix(), "")
	}
}

func TestPanicf(t *testing.T) {
	a := assert.New(t)
	resetLog(defaultLogs, t)

	Error("error")
	a.True(errorW.Len() > 0)
	a.Equal(debugW.Len(), 0)

	a.Panic(func() {
		Panicf(LevelError, "panic")
	})

	a.True(infoW.Len() == 0)
	a.True(errorW.Len() > 0)
}
