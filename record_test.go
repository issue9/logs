// SPDX-License-Identifier: MIT

package logs

import (
	"strings"
	"testing"

	"github.com/issue9/assert/v3"
)

var _ Logger = &Record{}

func TestRecord_location(t *testing.T) {
	a := assert.New(t, false)
	l := New(nil, Caller, Created)

	e := l.NewRecord(LevelWarn)
	a.NotNil(e)
	a.Empty(e.Path).Zero(e.Line)

	e.setLocation(1)
	a.True(strings.HasSuffix(e.Path, "record_test.go")).Equal(e.Line, 22)
}
