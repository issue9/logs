// SPDX-License-Identifier: MIT

package logs

import (
	"io"
	"log"
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/issue9/logs/v3/writers"
)

var (
	_ io.Writer     = &logger{}
	_ writers.Adder = &logger{}
)

func TestLogger_SetOutput(t *testing.T) {
	a := assert.New(t, false)

	l := newLogger("", 0)
	a.Equal(l.container.Len(), 0)

	cont := writers.NewContainer()
	a.NotError(l.SetOutput(cont))
	a.Equal(l.container.Len(), 1)

	// setOutput 会替换旧有的 writer
	a.NotError(l.SetOutput(cont))
	a.Equal(l.container.Len(), 1)

	a.NotError(l.SetOutput(nil))
	a.Equal(l.container.Len(), 0)
}

func TestParseFlag(t *testing.T) {
	a := assert.New(t, false)

	eq := func(str string, v int) {
		ret, err := parseFlag(str)
		a.NotError(err).Equal(ret, v)
	}

	eq("log.Ldate|log.ltime", log.Ldate|log.Ltime)
	eq("log.Ldate| log.Ltime", log.Ldate|log.Ltime)
	eq(" ", 0)
	eq("", 0)
}
