// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"testing"
	"time"

	"github.com/issue9/assert/v2"
	"github.com/issue9/term/v3/colors"
)

var benchEntry = &Entry{
	Level:   LevelWarn,
	Created: time.Now(),
	Message: "msg",
	Path:    "path.go",
	Line:    20,
	Pairs: []Pair{
		{K: "k1", V: "v1"},
		{K: "k2", V: "v2"},
	},
}

func BenchmarkTextWriter(b *testing.B) {
	buf := new(bytes.Buffer)
	w := NewTextWriter("2006-01-02", buf)

	for i := 0; i < b.N; i++ {
		w.WriteEntry(benchEntry)
	}
}

func BenchmarkJSONWriter(b *testing.B) {
	buf := new(bytes.Buffer)
	w := NewJSONWriter(buf)

	for i := 0; i < b.N; i++ {
		w.WriteEntry(benchEntry)
	}
}

func BenchmarkTermWriter(b *testing.B) {
	buf := new(bytes.Buffer)
	w := NewTermWriter("2006-01-02", colors.Red, buf)

	for i := 0; i < b.N; i++ {
		w.WriteEntry(benchEntry)
	}
}

func BenchmarkEntry_Printf(b *testing.B) {
	a := assert.New(b, false)
	buf := new(bytes.Buffer)
	l := New(NewTextWriter("2006-01-02", buf))
	a.NotNil(l)
	l.Enable(LevelError)

	err := l.ERROR()
	for i := 0; i < b.N; i++ {
		err.Value("k1", "v1").Printf("p1")
	}
}

func BenchmarkLogs_disableLogger(b *testing.B) {
	a := assert.New(b, false)
	buf := new(bytes.Buffer)
	w := NewTextWriter("2006-01-02", buf)
	l := New(w)
	a.NotNil(l)
	l.Enable(LevelInfo)

	err := l.ERROR()
	for i := 0; i < b.N; i++ {
		err.Value("k1", "v1").Printf("p1")
	}
}

func BenchmarkLogs_nop(b *testing.B) {
	a := assert.New(b, false)
	l := New(nil)
	a.NotNil(l)
	l.Enable(LevelError)

	err := l.ERROR()
	for i := 0; i < b.N; i++ {
		err.Value("k1", "v1").Printf("p1")
	}
}
