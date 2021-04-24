// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"encoding/xml"
	"io"
	"testing"

	"github.com/issue9/assert"

	"github.com/issue9/logs/v2/config"
	"github.com/issue9/logs/v2/internal/initfunc"
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

	a.NotError(logs.loggers[LevelInfo].setOutput(infoW))
	a.NotError(logs.loggers[LevelDebug].setOutput(debugW))
	a.NotError(logs.loggers[LevelError].setOutput(errorW))
	a.NotError(logs.loggers[LevelTrace].setOutput(traceW))
	a.NotError(logs.loggers[LevelWarn].setOutput(warnW))
	a.NotError(logs.loggers[LevelCritical].setOutput(criticalW))
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
