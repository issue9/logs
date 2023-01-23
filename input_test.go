// SPDX-License-Identifier: MIT

package logs

import (
	"strings"
	"testing"

	"github.com/issue9/assert/v3"
)

var (
	_ Input = &emptyInput{}
	_ Input = &Entry{}
)

func TestEntry_location(t *testing.T) {
	a := assert.New(t, false)
	l := New(nil, Caller, Created)

	e := l.NewEntry(LevelWarn)
	a.NotNil(e)
	a.Empty(e.Path).Zero(e.Line)

	e.Location(1)
	a.True(strings.HasSuffix(e.Path, "input_test.go")).Equal(e.Line, 25)
}
