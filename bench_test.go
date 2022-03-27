// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"testing"

	"github.com/issue9/assert/v2"
)

func BenchmarkEntry_Printf(b *testing.B) {
	a := assert.New(b, false)
	l := New()
	a.NotNil(l)
	buf := new(bytes.Buffer)
	l.SetOutput(NewWriter(TextFormat("2006-01-02"), buf), LevelError)

	err := l.ERROR()
	for i := 0; i < b.N; i++ {
		err.Value("k1", "v1").Printf("p1")
	}
}

func BenchmarkLogs_empty(b *testing.B) {
	a := assert.New(b, false)
	l := New()
	a.NotNil(l)
	buf := new(bytes.Buffer)
	l.SetOutput(NewWriter(TextFormat("2006-01-02"), buf), LevelInfo)
	l.Enable(LevelInfo)

	err := l.ERROR()
	for i := 0; i < b.N; i++ {
		err.Value("k1", "v1").Printf("p1")
	}
}
