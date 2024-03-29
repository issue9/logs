// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestLogs_WithStd(t *testing.T) {
	a := assert.New(t, false)

	buf := new(bytes.Buffer)
	New(NewTextHandler(buf), WithLocation(true), WithCreated(MilliLayout), WithStd())

	slog.Error("error", "a1", "v1")
	a.Contains(buf.String(), "error").Contains(buf.String(), "a1=v1")

	l2 := slog.With("attr1", "val1")
	l2.Warn("warn")
	a.Contains(buf.String(), "warn").Contains(buf.String(), "attr1=val1")

	l3 := l2.WithGroup("g1")
	l3.Warn("group")
	a.Contains(buf.String(), "g1.attr1=val1")
}
