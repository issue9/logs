// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
	"github.com/issue9/term/v3/colors"
)

func newRecord(a *assert.Assertion, logs *Logs, lv Level) *Record {
	e := logs.NewRecord(lv)
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

func TestTextHandler(t *testing.T) {
	a := assert.New(t, false)
	layout := MilliLayout
	now := time.Now()
	l := New(nil, Created, Caller)

	e := newRecord(a, l, LevelWarn)
	e.Created = now

	a.PanicString(func() {
		NewTextHandler(layout)
	}, "参数 w 不能为空")

	buf := new(bytes.Buffer)
	l.SetHandler(NewTextHandler(layout, buf))
	e.output()
	a.Equal(buf.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2\n")

	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	b3 := new(bytes.Buffer)
	l.SetHandler(NewTextHandler(layout, b1, b2, b3))
	e.output()
	a.Equal(b1.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2\n")
	a.Equal(b2.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2\n")
	a.Equal(b3.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2\n")
}

func TestJSONFormat(t *testing.T) {
	a := assert.New(t, false)
	now := time.Now()

	e := newRecord(a, New(nil), LevelWarn)
	e.Created = now

	a.PanicString(func() {
		NewJSONHandler(MicroLayout)
	}, "参数 w 不能为空")

	buf := new(bytes.Buffer)
	NewJSONHandler(MicroLayout, buf).Handle(e)
	a.True(json.Valid(buf.Bytes())).
		Contains(buf.String(), LevelWarn.String()).
		Contains(buf.String(), "k1")

	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	NewJSONHandler(MicroLayout, b1, b2).Handle(e)
	a.True(json.Valid(b1.Bytes())).
		Contains(b1.String(), LevelWarn.String()).
		Contains(b1.String(), "k1")
	a.Equal(b1.String(), b2.String())
}

func TestTermHandler(t *testing.T) {
	a := assert.New(t, false)

	t.Log("此测试将在终端输出一段带颜色的日志记录")

	l := New(nil)
	e := newRecord(a, l, LevelWarn)
	e.Created = time.Now()
	w := NewTermHandler(MilliLayout, colors.Blue, os.Stdout)
	w.Handle(e)

	l = New(nil, Caller, Created)
	e = newRecord(a, l, LevelError)
	e.Message = "error message"
	w = NewTermHandler(MicroLayout, colors.Red, os.Stdout)
	w.Handle(e)
}

func TestDispatchHandler(t *testing.T) {
	a := assert.New(t, false)

	txtBuf := new(bytes.Buffer)
	jsonBuf := new(bytes.Buffer)

	w := NewDispatchHandler(map[Level]Handler{
		LevelInfo: NewTextHandler(NanoLayout, txtBuf),
		LevelWarn: NewJSONHandler(MicroLayout, jsonBuf),
	})
	l := New(w)

	e := l.NewRecord(LevelWarn)
	e.Created = time.Now()
	l.Warnf("warnf test")
	a.Zero(txtBuf.Len()).Contains(jsonBuf.String(), "warnf test").True(json.Valid(jsonBuf.Bytes()))

	e.Level = LevelInfo
	jsonBuf.Reset()
	l.Info("info test")
	a.Zero(jsonBuf.Len()).Contains(txtBuf.String(), "info test")
}
