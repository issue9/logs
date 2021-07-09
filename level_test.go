// SPDX-License-Identifier: MIT

package logs

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLogs_logs(t *testing.T) {
	a := assert.New(t)

	l, err := New(nil)
	a.NotError(err).NotNil(l)

	logs := l.logs(LevelDebug)
	a.Equal(1, len(logs)).Equal(logs[0].level, LevelDebug)

	logs = l.logs(LevelDebug | LevelWarn)
	a.Equal(2, len(logs)).
		Equal(logs[0].level, LevelDebug).
		Equal(logs[1].level, LevelWarn)

	logs = l.logs(LevelAll)
	a.Equal(len(levels), len(logs))
}
