// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"strings"
	"testing"

	"github.com/issue9/assert/v2"
)

var (
	_ Logger = &entry{}
	_ Logger = &logger{}
	_ Logger = &emptyLogger{}
)

func TestEntry_Location(t *testing.T) {
	a := assert.New(t, false)

	e := NewEntry()
	a.NotNil(e)
	a.Empty(e.Path).Zero(e.Line)

	e.Location(1)
	a.True(strings.HasSuffix(e.Path, "logger_test.go")).Equal(e.Line, 26)
}

func TestLogger_location(t *testing.T) {
	a := assert.New(t, false)
	l := New()
	a.NotNil(l)

	buf := new(bytes.Buffer)
	w := NewWriter(TextFormat("2006-01-02"), buf)
	l.SetOutput(w, LevelInfo, LevelError)
	l.ERROR().Value("k1", "v1").
		Printf("pf") // 位置记录此行
	val := buf.String()
	a.Contains(val, "logger_test.go:39").Contains(val, "k1=v1")
}
