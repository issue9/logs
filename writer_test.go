// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert/v2"
	"github.com/issue9/term/v3/colors"
)

func TestTextFormat(t *testing.T) {
	a := assert.New(t, false)
	layout := "15:04:05"
	now := time.Now()

	e := &Entry{
		Level:   LevelWarn,
		Created: now,
		Message: "msg",
		Path:    "path.go",
		Line:    20,
		Pairs: []Pair{
			{K: "k1", V: "v1"},
			{K: "k2", V: "v2"},
		},
	}
	data := TextFormat(layout)(e)
	a.Equal(string(data), "[WARN] "+now.Format(layout)+" msg\tpath.go:20 k1=v1 k2=v2")
}

func TestJSONFormat(t *testing.T) {
	a := assert.New(t, false)
	now := time.Now()

	e := &Entry{
		Level:   LevelWarn,
		Created: now,
		Message: "msg",
		Path:    "path.go",
		Line:    20,
		Pairs: []Pair{
			{K: "k1", V: "v1"},
			{K: "k2", V: "v2"},
		},
	}
	data := JSONFormat(e)
	a.True(json.Valid(data)).
		Contains(string(data), LevelWarn.String()).
		Contains(string(data), "k1")
}

func TestNewWriter(t *testing.T) {
	a := assert.New(t, false)
	layout := "15:04:05"
	e := &Entry{
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

	b1 := new(bytes.Buffer)
	w1 := NewWriter(TextFormat(layout), b1)
	w1.WriteEntry(e)
	a.True(b1.Len() > 0).Contains(b1.String(), "path.go")

	b1.Reset()
	b2 := new(bytes.Buffer)
	w2 := NewWriter(TextFormat(layout), b1, b2)
	w2.WriteEntry(e)
	a.True(b1.Len() > 0).
		Equal(b1.String(), b2.String()).
		Contains(b1.String(), "path.go")
}

func TestTermWriter(t *testing.T) {
	t.Log("此测试将在终端输出一段带颜色的日志记录")

	layout := "15:04:05"
	e := &Entry{
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

	w := NewTermWriter(layout, colors.Blue, os.Stdout)
	w.WriteEntry(e)

	e.Level = LevelError
	e.Message = "error message"
	w = NewTermWriter(layout, colors.Red, os.Stdout)
	w.WriteEntry(e)
}
