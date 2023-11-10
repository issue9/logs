// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"errors"
	"testing"

	"github.com/issue9/assert/v3"
	"github.com/issue9/localeutil"
)

func BenchmarkTextHandler(b *testing.B) {
	benchHandler(NewTextHandler(new(bytes.Buffer)), b)
}

func BenchmarkJSONHandler(b *testing.B) {
	benchHandler(NewJSONHandler(new(bytes.Buffer)), b)
}

func BenchmarkTermHandler(b *testing.B) {
	benchHandler(NewTermHandler(new(bytes.Buffer), nil), b)
}

func benchHandler(h Handler, b *testing.B) {
	b.Run("default", func(b *testing.B) {
		l := New(h)
		e := l.ERROR().With("k1", "v1").With("k2", 2).With("k3", localeutil.Phrase("lang"))

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			e.String("err")
		}
	})

	b.Run("withAttr", func(b *testing.B) {
		l := New(h)
		err := l.ERROR().New(map[string]any{"k1": "v1", "k2": 2, "k3": localeutil.Phrase("lang")})

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err.String("err")
		}
	})
}

func BenchmarkWithRecorder(b *testing.B) {
	a := assert.New(b, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf))
	a.NotNil(l)
	l.Enable(LevelError)

	err := l.ERROR()
	e := errors.New("err")

	b.Run("print", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err.With("k1", "v1").Print(e)
		}
	})

	b.Run("printf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err.With("k1", "v1").Printf("%v", e)
		}
	})

	b.Run("error", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err.With("k1", "v1").Error(e)
		}
	})

	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err.With("k1", "v1").String("err")
		}
	})
}

func BenchmarkLogs_disableRecorder(b *testing.B) {
	a := assert.New(b, false)
	buf := new(bytes.Buffer)
	w := NewTextHandler(buf)
	l := New(w)
	a.NotNil(l)
	l.Enable(LevelInfo)

	err := l.ERROR()
	for i := 0; i < b.N; i++ {
		err.With("k1", "v1").Printf("p1")
	}
}

func BenchmarkLogs_nop(b *testing.B) {
	a := assert.New(b, false)
	l := New(nil)
	a.NotNil(l)
	l.Enable(LevelError)

	err := l.ERROR()
	for i := 0; i < b.N; i++ {
		err.With("k1", "v1").Printf("p1")
	}
}

func BenchmarkLogger_LogLogger(b *testing.B) {
	a := assert.New(b, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf))
	a.NotNil(l)
	l.Enable(LevelError)
	err := l.ERROR()

	for i := 0; i < b.N; i++ {
		err.LogLogger().Printf("std log")
	}
}

func BenchmarkLogs_LogLogger_withDisable(b *testing.B) {
	a := assert.New(b, false)
	buf := new(bytes.Buffer)
	l := New(NewTextHandler(buf))
	a.NotNil(l)
	l.Enable(LevelInfo)
	err := l.ERROR()

	for i := 0; i < b.N; i++ {
		err.LogLogger().Printf("std log")
	}
}
