// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"testing"

	"github.com/issue9/assert/v2"
	"github.com/issue9/term/v3/colors"
)

func BenchmarkTextWriter(b *testing.B) {
	a := assert.New(b, false)

	buf := new(bytes.Buffer)
	w := NewTextWriter("2006-01-02", buf)
	l := New(w, Caller, Created)
	e := newEntry(a, l, LevelWarn)

	for i := 0; i < b.N; i++ {
		w.WriteEntry(e)
	}
}

func BenchmarkJSONWriter(b *testing.B) {
	a := assert.New(b, false)

	buf := new(bytes.Buffer)
	w := NewJSONWriter(true, buf)
	l := New(w, Caller, Created)
	e := newEntry(a, l, LevelWarn)

	for i := 0; i < b.N; i++ {
		w.WriteEntry(e)
	}
}

func BenchmarkTermWriter(b *testing.B) {
	a := assert.New(b, false)

	buf := new(bytes.Buffer)
	w := NewTermWriter("2006-01-02", colors.Red, buf)
	l := New(w, Caller, Created)
	e := newEntry(a, l, LevelWarn)

	for i := 0; i < b.N; i++ {
		w.WriteEntry(e)
	}
}

func BenchmarkEntry_Printf(b *testing.B) {
	a := assert.New(b, false)
	buf := new(bytes.Buffer)
	l := New(NewTextWriter("2006-01-02", buf), Created, Caller)
	a.NotNil(l)
	l.Enable(LevelError)

	err := l.ERROR()
	for i := 0; i < b.N; i++ {
		err.Value("k1", "v1").Printf("p1")
	}
}

func BenchmarkEntry_Printf_withCallerAndCreated(b *testing.B) {
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

func BenchmarkLogs_StdLogger(b *testing.B) {
	a := assert.New(b, false)
	buf := new(bytes.Buffer)
	l := New(NewTextWriter("2006-01-02", buf))
	a.NotNil(l)
	l.Enable(LevelError)

	err := l.StdLogger(LevelError)
	for i := 0; i < b.N; i++ {
		err.Printf("std log")
	}
}
