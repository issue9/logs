// SPDX-License-Identifier: MIT

//go:build go1.21

package logs

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestLogs_SLog(t *testing.T) {
	a := assert.New(t, false)

	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MilliLayout))
	slog.SetDefault(l.SLog())

	slog.Error("error", "a1", "v1")
	a.Contains(buf.String(), "error").Contains(buf.String(), "a1=v1")

	l2 := slog.With("attr1", "val1")
	l2.Warn("warn")
	a.Contains(buf.String(), "warn").Contains(buf.String(), "attr1=val1")

	l3 := l2.WithGroup("g1")
	l3.Warn("group")
	a.Contains(buf.String(), "g1.attr1=val1")
}
