// SPDX-License-Identifier: MIT

package logs

import (
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/issue9/logs/v3/writers"
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
