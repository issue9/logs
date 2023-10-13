// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
	"github.com/issue9/term/v3/colors"
)

type (
	marshalObject    string
	marshalErrObject string
)

func (o marshalObject) MarshalText() ([]byte, error) {
	return []byte(o), nil
}

func (o marshalObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(o))
}

func (o marshalErrObject) MarshalText() ([]byte, error) {
	return nil, errors.New("marshal text error")
}

func (o marshalErrObject) MarshalJSON() ([]byte, error) {
	return nil, errors.New("marshal json error")
}

func newRecord(a *assert.Assertion, logs *Logs, lv Level) *Record {
	e := logs.NewRecord(lv)
	a.NotNil(e)

	e.AppendMessage = func(bs []byte) []byte { return append(bs, "msg"...) }
	e.Path = "path.go:20"
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
	l := New(nil, WithCreated(layout), WithCaller())

	e := newRecord(a, l, LevelWarn)
	e.Created = now
	e.With("m1", marshalObject("m1"))

	a.PanicString(func() {
		NewTextHandler()
	}, "参数 w 不能为空")

	buf := new(bytes.Buffer)
	l.SetHandler(NewTextHandler(buf))
	e.output()
	a.Equal(buf.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1\n")

	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	b3 := new(bytes.Buffer)
	l.SetHandler(NewTextHandler(b1, b2, b3))
	e.output()
	a.Equal(b1.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1\n")
	a.Equal(b2.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1\n")
	a.Equal(b3.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1\n")

	// error

	e.With("m2", marshalErrObject("m2"))
	buf.Reset()
	l.SetHandler(NewTextHandler(buf))
	e.output()
	a.Equal(buf.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1 m2=Err(marshal text error)\n")
}

func TestJSONFormat(t *testing.T) {
	a := assert.New(t, false)
	layout := MilliLayout
	now := time.Now()
	l := New(nil, WithCreated(layout), WithCaller())

	e := newRecord(a, l, LevelWarn)
	e.Created = now
	e.With("m1", marshalObject("m1"))

	a.PanicString(func() {
		NewJSONHandler()
	}, "参数 w 不能为空")

	buf := new(bytes.Buffer)
	l.SetHandler(NewJSONHandler(buf))
	e.output()
	a.Equal(buf.String(), `{"level":"WARN","message":"msg","created":"`+now.Format(layout)+`","path":"path.go:20","params":[{"k1":"v1"},{"k2":"v2"},{"m1":"m1"}]}`)

	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	l.SetHandler(NewJSONHandler(b1, b2))
	e.output()
	a.Equal(b1.String(), `{"level":"WARN","message":"msg","created":"`+now.Format(layout)+`","path":"path.go:20","params":[{"k1":"v1"},{"k2":"v2"},{"m1":"m1"}]}`).
		Equal(b1.String(), b2.String())

	// error

	e.With("m2", marshalErrObject("m2"))
	buf.Reset()
	l.SetHandler(NewJSONHandler(buf))
	e.output()
	a.Equal(buf.String(), `{"level":"WARN","message":"msg","created":"`+now.Format(layout)+`","path":"path.go:20","params":[{"k1":"v1"},{"k2":"v2"},{"m1":"m1"},{"m2":"Err(json: error calling MarshalJSON for type logs.marshalErrObject: marshal json error)"}]}`)
}

func TestTermHandler(t *testing.T) {
	a := assert.New(t, false)

	t.Log("此测试将在终端输出一段带颜色的日志记录")

	l := New(nil)
	e := newRecord(a, l, LevelWarn)
	e.Created = time.Now()
	w := NewTermHandler(os.Stdout, map[Level]colors.Color{LevelError: colors.Green})
	w.Handle(e)

	l = New(nil, WithCaller(), WithCreated(MicroLayout))
	e = newRecord(a, l, LevelError)
	e.AppendMessage = func(bs []byte) []byte { return append(bs, "error message"...) }
	w = NewTermHandler(os.Stdout, map[Level]colors.Color{LevelError: colors.Green})
	w.Handle(e)
}

func TestDispatchHandler(t *testing.T) {
	a := assert.New(t, false)

	txtBuf := new(bytes.Buffer)
	jsonBuf := new(bytes.Buffer)

	w := NewDispatchHandler(map[Level]Handler{
		LevelInfo: NewTextHandler(txtBuf),
		LevelWarn: NewJSONHandler(jsonBuf),
	})
	l := New(w)

	e := l.NewRecord(LevelWarn)
	e.Created = time.Now()
	l.WARN().Printf("warnf test")
	a.Zero(txtBuf.Len()).Contains(jsonBuf.String(), "warnf test").True(json.Valid(jsonBuf.Bytes()))

	e.Level = LevelInfo
	jsonBuf.Reset()
	l.INFO().Print("info test")
	a.Zero(jsonBuf.Len()).Contains(txtBuf.String(), "info test")
}
