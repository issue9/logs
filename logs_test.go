// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/issue9/assert"

	"github.com/issue9/logs/v2/config"
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

	logs.loggers[LevelInfo].setOutput(infoW, "[INFO]", log.LstdFlags)
	logs.loggers[LevelDebug].setOutput(debugW, "[DEBUG]", log.LstdFlags)
	logs.loggers[LevelError].setOutput(errorW, "[ERROR]", log.LstdFlags)
	logs.loggers[LevelTrace].setOutput(traceW, "[TRACE]", log.LstdFlags)
	logs.loggers[LevelWarn].setOutput(warnW, "[WARN]", log.LstdFlags)
	logs.loggers[LevelCritical].setOutput(criticalW, "[CRITICAL]", log.LstdFlags)
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

func TestSetWriter(t *testing.T) {
	a := assert.New(t)

	a.NotError(defaultLogs.SetOutput(LevelError, nil, "", 0))

	a.Error(defaultLogs.SetOutput(100, nil, "", 0))
}

func debugWInit(cfg *config.Config) (io.Writer, error) {
	return debugW, nil
}
