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

func TestTextWriter(t *testing.T) {
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
	buf := new(bytes.Buffer)
	NewTextWriter(layout, buf).WriteEntry(e)

	a.Equal(buf.String(), "[WARN] "+now.Format(layout)+" msg\tpath.go:20 k1=v1 k2=v2\n")
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
	buf := new(bytes.Buffer)
	NewJSONWriter(buf).WriteEntry(e)

	a.True(json.Valid(buf.Bytes())).
		Contains(buf.String(), LevelWarn.String()).
		Contains(buf.String(), "k1")
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
