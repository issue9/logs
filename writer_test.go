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

func newEntry(a *assert.Assertion, logs *Logs, lv Level) *Entry {
	e := logs.NewEntry(lv)
	a.NotNil(e)

	e.Message = "msg"
	e.Path = "path.go"
	e.Line = 20
	e.Params = []Pair{
		{K: "k1", V: "v1"},
		{K: "k2", V: "v2"},
	}

	return e
}

func TestTextWriter(t *testing.T) {
	a := assert.New(t, false)
	layout := "15:04:05"
	now := time.Now()
	l := New(nil, Created, Caller)

	e := newEntry(a, l, LevelWarn)
	e.Created = now

	a.PanicString(func() {
		NewTextWriter(layout)
	}, "参数 w 不能为空")

	buf := new(bytes.Buffer)
	l.SetOutput(NewTextWriter(layout, buf))
	l.Output(e)
	a.Equal(buf.String(), "[WARN] "+now.Format(layout)+" msg\tpath.go:20 k1=v1 k2=v2\n")

	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	b3 := new(bytes.Buffer)
	l.SetOutput(NewTextWriter(layout, b1, b2, b3))
	l.Output(e)
	a.Equal(b1.String(), "[WARN] "+now.Format(layout)+" msg\tpath.go:20 k1=v1 k2=v2\n")
	a.Equal(b2.String(), "[WARN] "+now.Format(layout)+" msg\tpath.go:20 k1=v1 k2=v2\n")
	a.Equal(b3.String(), "[WARN] "+now.Format(layout)+" msg\tpath.go:20 k1=v1 k2=v2\n")
}

func TestJSONFormat(t *testing.T) {
	a := assert.New(t, false)
	now := time.Now()

	e := newEntry(a, New(nil), LevelWarn)
	e.Created = now

	a.PanicString(func() {
		NewJSONWriter(false)
	}, "参数 w 不能为空")

	buf := new(bytes.Buffer)
	NewJSONWriter(true, buf).WriteEntry(e)
	a.True(json.Valid(buf.Bytes())).
		Contains(buf.String(), LevelWarn.String()).
		Contains(buf.String(), "k1")

	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	NewJSONWriter(true, b1, b2).WriteEntry(e)
	a.True(json.Valid(b1.Bytes())).
		Contains(b1.String(), LevelWarn.String()).
		Contains(b1.String(), "k1")
	a.Equal(b1.String(), b2.String())
}

func TestTermWriter(t *testing.T) {
	a := assert.New(t, false)
	layout := "15:04:05"

	t.Log("此测试将在终端输出一段带颜色的日志记录")

	l := New(nil)
	e := newEntry(a, l, LevelWarn)
	e.Created = time.Now()
	w := NewTermWriter(layout, colors.Blue, os.Stdout)
	w.WriteEntry(e)

	l = New(nil, Caller, Created)
	e = newEntry(a, l, LevelError)
	e.Message = "error message"
	w = NewTermWriter(layout, colors.Red, os.Stdout)
	w.WriteEntry(e)
}

func TestDispatchWriter(t *testing.T) {
	a := assert.New(t, false)
	layout := "15:04:05"

	txtBuf := new(bytes.Buffer)
	jsonBuf := new(bytes.Buffer)

	w := NewDispatchWriter(map[Level]Writer{
		LevelInfo: NewTextWriter(layout, txtBuf),
		LevelWarn: NewJSONWriter(true, jsonBuf),
	})
	l := New(w)

	e := l.NewEntry(LevelWarn)
	e.Created = time.Now()
	l.Warnf("warnf test")
	a.Zero(txtBuf.Len()).Contains(jsonBuf.String(), "warnf test").True(json.Valid(jsonBuf.Bytes()))

	e.Level = LevelInfo
	jsonBuf.Reset()
	l.Info("info test")
	a.Zero(jsonBuf.Len()).Contains(txtBuf.String(), "info test")
}
