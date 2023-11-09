// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"encoding/json"
	"errors"
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

func newRecord(a *assert.Assertion) *Record {
	e := NewRecord()
	a.NotNil(e)

	e.AppendMessage = func(b *Buffer) { b.AppendString("msg") }
	e.AppendLocation = func(b *Buffer) { b.AppendString("path.go:20") }
	e.Attrs = []Attr{
		{K: "k1", V: "v1"},
		{K: "k2", V: "v2"},
	}

	return e
}

func TestTextHandler(t *testing.T) {
	a := assert.New(t, false)

	a.Panic(func() { NewTextHandler(nil) })

	layout := MilliLayout
	now := time.Now()

	l := New(nil, WithCreated(layout), WithLocation(true))
	e := newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1"))

	buf := new(bytes.Buffer)
	l = New(NewTextHandler(buf), WithCreated(layout), WithLocation(true))
	e = newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1"))
	e.output(l.detail, l.WARN().Handler())
	a.Equal(buf.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1\n")

	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	b3 := new(bytes.Buffer)
	l = New(NewTextHandler(b1, b2, b3), WithCreated(layout), WithLocation(true))
	e = newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1"))
	e.output(l.detail, l.WARN().Handler())
	a.Equal(b1.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1\n")
	a.Equal(b2.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1\n")
	a.Equal(b3.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1\n")

	// error

	buf.Reset()
	l = New(NewTextHandler(buf), WithCreated(layout), WithLocation(true))
	e = newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1")).with(l, "m2", marshalErrObject("m2"))
	e.output(l.detail, l.WARN().Handler())
	a.Equal(buf.String(), "[WARN] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1 m2=Err(marshal text error)\n")

	// Handler.New

	buf.Reset()
	h := NewTextHandler(buf)
	l = New(h, WithLocation(true), WithDetail(true))
	e = newRecord(a)
	h = h.New(l.Detail(), LevelWarn, []Attr{{K: "attr1", V: 3.51}})
	h.Handle(e)
	a.Equal(buf.String(), "[WARN] path.go:20\tmsg attr1=3.51 k1=v1 k2=v2\n")

	// Handler.New().New()
	buf.Reset()
	e = newRecord(a)
	h.New(l.Detail(), LevelWarn, []Attr{{K: "a1", V: int8(5)}, {K: "a2", V: uint(8)}}).Handle(e)
	a.Equal(buf.String(), "[WARN] path.go:20\tmsg attr1=3.51 a1=5 a2=8 k1=v1 k2=v2\n")
}

func TestJSONFormat(t *testing.T) {
	a := assert.New(t, false)
	layout := MilliLayout
	now := time.Now()

	l := New(nil, WithCreated(layout), WithLocation(true))
	e := newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1"))

	a.Panic(func() { NewJSONHandler() })

	buf := new(bytes.Buffer)
	l = New(NewJSONHandler(buf), WithCreated(layout), WithLocation(true))
	e = newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1"))
	e.output(l.detail, l.WARN().Handler())
	a.Equal(buf.String(), `{"level":"WARN","message":"msg","created":"`+now.Format(layout)+`","path":"path.go:20","attrs":[{"k1":"v1"},{"k2":"v2"},{"m1":"m1"}]}`)

	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)
	l = New(NewJSONHandler(b1, b2), WithCreated(layout), WithLocation(true))
	e = newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1"))
	e.output(l.detail, l.WARN().Handler())
	a.Equal(b1.String(), `{"level":"WARN","message":"msg","created":"`+now.Format(layout)+`","path":"path.go:20","attrs":[{"k1":"v1"},{"k2":"v2"},{"m1":"m1"}]}`).
		Equal(b1.String(), b2.String())

	// error

	buf.Reset()
	l = New(NewJSONHandler(buf), WithCreated(layout), WithLocation(true))
	e = newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1")).with(l, "m2", marshalErrObject("m2"))
	e.output(l.detail, l.WARN().Handler())
	a.Equal(buf.String(), `{"level":"WARN","message":"msg","created":"`+now.Format(layout)+`","path":"path.go:20","attrs":[{"k1":"v1"},{"k2":"v2"},{"m1":"m1"},{"m2":"Err(json: error calling MarshalJSON for type logs.marshalErrObject: marshal json error)"}]}`)

	// Handler.New

	buf.Reset()
	h := NewJSONHandler(buf)
	l = New(h, WithLocation(true), WithDetail(true))
	e = newRecord(a)
	h = h.New(l.Detail(), LevelWarn, []Attr{{K: "attr1", V: 3.5}})
	h.Handle(e)
	a.Equal(buf.String(), `{"level":"WARN","message":"msg","path":"path.go:20","attrs":[{"attr1":3.5},{"k1":"v1"},{"k2":"v2"}]}`)

	// Handler.New().New()
	buf.Reset()
	e = newRecord(a)
	h.New(l.Detail(), LevelWarn, []Attr{{K: "a1", V: int8(5)}, {K: "a2", V: uint(8)}}).Handle(e)
	a.Equal(buf.String(), `{"level":"WARN","message":"msg","path":"path.go:20","attrs":[{"attr1":3.5},{"a1":5},{"a2":8},{"k1":"v1"},{"k2":"v2"}]}`)
}

func TestTermHandler(t *testing.T) {
	a := assert.New(t, false)

	a.PanicString(func() { NewTermHandler(nil, nil) }, "参数 w 不能为空")

	layout := MilliLayout
	now := time.Now()

	l := New(nil, WithCreated(layout), WithLocation(true))
	e := newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1"))

	buf := new(bytes.Buffer)
	l = New(NewTermHandler(buf, nil), WithCreated(layout), WithLocation(true))
	e = newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1"))
	e.output(l.detail, l.WARN().Handler())
	a.Equal(buf.String(), "[\033[33;49mWARN\033[0m] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1\n")

	// error

	buf.Reset()
	l = New(NewTermHandler(buf, map[Level]colors.Color{LevelWarn: colors.Red}), WithCreated(layout), WithLocation(true))
	e = newRecord(a)
	e.AppendCreated = func(b *Buffer) { b.AppendTime(now, l.createdFormat) }
	e.with(l, "m1", marshalObject("m1")).with(l, "m2", marshalErrObject("m2"))
	e.output(l.detail, l.WARN().Handler())
	a.Equal(buf.String(), "[\033[31;49mWARN\033[0m] "+now.Format(layout)+" path.go:20\tmsg k1=v1 k2=v2 m1=m1 m2=Err(marshal text error)\n")

	// Handler.New

	buf.Reset()
	h := NewTermHandler(buf, nil)
	l = New(h, WithLocation(true), WithDetail(true))
	e = newRecord(a)
	h = h.New(l.Detail(), LevelWarn, []Attr{{K: "attr1", V: 3.51}})
	h.Handle(e)
	a.Equal(buf.String(), "[\033[33;49mWARN\033[0m] path.go:20\tmsg attr1=3.51 k1=v1 k2=v2\n")

	// Handler.New().New()
	buf.Reset()
	e = newRecord(a)
	h.New(l.Detail(), LevelWarn, []Attr{{K: "a1", V: int8(5)}, {K: "a2", V: uint(8)}}).Handle(e)
	a.Equal(buf.String(), "[\033[33;49mWARN\033[0m] path.go:20\tmsg attr1=3.51 a1=5 a2=8 k1=v1 k2=v2\n")
}

func TestDispatchHandler(t *testing.T) {
	a := assert.New(t, false)

	textBuf := new(bytes.Buffer)
	jsonBuf := new(bytes.Buffer)

	a.PanicString(func() {
		NewDispatchHandler(map[Level]Handler{
			LevelInfo: NewTextHandler(textBuf),
			LevelWarn: NewJSONHandler(jsonBuf),
		})
	}, "需指定所有 Level 对应的对象")

	w := NewDispatchHandler(map[Level]Handler{
		LevelInfo:  NewTextHandler(textBuf),
		LevelWarn:  NewJSONHandler(jsonBuf),
		LevelDebug: NewTextHandler(textBuf),
		LevelError: NewTextHandler(textBuf),
		LevelFatal: NewTextHandler(textBuf),
		LevelTrace: NewTextHandler(textBuf),
	})

	l := New(w)
	l.WARN().Printf("warnf test")
	l.INFO().Print("info test")

	a.Equal(jsonBuf.String(), `{"level":"WARN","message":"warnf test"}`).
		True(json.Valid(jsonBuf.Bytes())).
		Equal(textBuf.String(), "[INFO] info test\n")
}

func TestMergeHandler(t *testing.T) {
	a := assert.New(t, false)

	textBuf := new(bytes.Buffer)
	jsonBuf := new(bytes.Buffer)
	w := MergeHandler(NewTextHandler(textBuf), NewJSONHandler(jsonBuf))

	l := New(w)
	l.WARN().Printf("warnf test")

	a.Equal(jsonBuf.String(), `{"level":"WARN","message":"warnf test"}`).
		True(json.Valid(jsonBuf.Bytes())).
		Equal(textBuf.String(), "[WARN] warnf test\n")

	// Handler.New()

	textBuf.Reset()
	jsonBuf.Reset()

	w = w.New(l.Detail(), LevelWarn, []Attr{{K: "a1", V: "v1"}})
	l = New(w)
	l.WARN().Printf("warnf test")

	a.Equal(jsonBuf.String(), `{"level":"WARN","message":"warnf test","attrs":[{"a1":"v1"}]}`).
		True(json.Valid(jsonBuf.Bytes()), jsonBuf.String()).
		Equal(textBuf.String(), "[WARN] warnf test a1=v1\n")

	// Handler.New().New()

	textBuf.Reset()
	jsonBuf.Reset()

	w = w.New(l.Detail(), LevelWarn, []Attr{{K: "a2", V: uint8(3)}})
	l = New(w)
	l.WARN().Printf("warnf test")

	a.Equal(jsonBuf.String(), `{"level":"WARN","message":"warnf test","attrs":[{"a1":"v1"},{"a2":3}]}`).
		True(json.Valid(jsonBuf.Bytes()), jsonBuf.String()).
		Equal(textBuf.String(), "[WARN] warnf test a1=v1 a2=3\n")
}
