// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/issue9/assert/v2"
)

var (
	_ Logger    = &logger{}
	_ io.Writer = &logger{}
	_ Logger    = &emptyLogger{}
	_ io.Writer = &emptyLogger{}
)

func TestEntry_Location(t *testing.T) {
	a := assert.New(t, false)
	l := New(nil, Caller, Created)

	e := l.NewEntry(LevelWarn)
	a.NotNil(e)
	a.Empty(e.Path).Zero(e.Line)

	e.Location(1)
	a.True(strings.HasSuffix(e.Path, "logger_test.go")).Equal(e.Line, 29)
}

func TestLogger_location(t *testing.T) {
	a := assert.New(t, false)

	buf := new(bytes.Buffer)
	l := New(NewTextWriter("2006-01-02", buf), Caller, Created)
	a.NotNil(l)
	l.Enable(LevelError)
	l.ERROR().Value("k1", "v1").
		Printf("pf") // 位置记录此行
	val := buf.String()
	a.Contains(val, "logger_test.go:41", val).
		Contains(val, "k1=v1")
}
