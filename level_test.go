// SPDX-License-Identifier: MIT

package logs

import (
	"testing"

	"github.com/issue9/assert/v2"
)

func TestLogs_walk(t *testing.T) {
	a := assert.New(t, false)

	l, err := New(nil)
	a.NotError(err).NotNil(l)

	num := 0
	a.NotError(l.walk(LevelDebug, func(ll *logger) error {
		num++
		return nil
	}))
	a.Equal(1, num)

	num = 0
	a.NotError(l.walk(LevelDebug|LevelWarn, func(ll *logger) error {
		num++
		return nil
	}))
	a.Equal(2, num)

	num = 0
	a.NotError(l.walk(LevelAll, func(ll *logger) error {
		num++
		return nil
	}))
	a.Equal(len(levels), num)
}
