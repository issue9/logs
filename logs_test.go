// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"bytes"
	"log"
	"testing"

	"github.com/issue9/assert"
)

var debugW = bytes.NewBufferString("")
var infoW = bytes.NewBufferString("")
var errorW = bytes.NewBufferString("")
var traceW = bytes.NewBufferString("")
var warnW = bytes.NewBufferString("")
var criticalW = bytes.NewBufferString("")

func resetLog() {
	infoW.Reset()
	debugW.Reset()
	errorW.Reset()
	traceW.Reset()
	warnW.Reset()
	criticalW.Reset()

	INFO = log.New(infoW, "[INFO]", log.LstdFlags)
	DEBUG = log.New(debugW, "[DEBUG]", log.LstdFlags)
	ERROR = log.New(errorW, "[ERROR]", log.LstdFlags)
	TRACE = log.New(traceW, "[TRACE]", log.LstdFlags)
	WARN = log.New(warnW, "[WARN]", log.LstdFlags)
	CRITICAL = log.New(criticalW, "[CRITICAL]", log.LstdFlags)
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
	resetLog()
	All("abc")
	checkLog(t)
}

func TestAllf(t *testing.T) {
	resetLog()
	Allf("abc")
	checkLog(t)
}
