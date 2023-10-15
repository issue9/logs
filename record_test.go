// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/issue9/assert/v3"
	"github.com/issue9/localeutil"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"golang.org/x/xerrors"
)

var _ Logger = &Record{}

type err struct {
	err error
}

func (e *err) Error() string { return e.err.Error() }

func (e *err) FormatError(p xerrors.Printer) error {
	p.Print("root\n")
	if p.Detail() {
		return e.err
	}
	return nil
}

func TestRecord_location(t *testing.T) {
	a := assert.New(t, false)
	l := New(nil, WithLocation(true), WithCreated(MicroLayout))

	e := l.NewRecord(LevelWarn)
	a.NotNil(e)
	a.Empty(e.AppendLocation)

	e.initLocationCreated(1) // 输出定位
	b := NewBuffer(false)
	defer b.Free()
	e.AppendLocation(b)
	a.True(strings.HasSuffix(string(b.data), "record_test.go:43"), string(b.data))
}

func TestRecord_Error(t *testing.T) {
	a := assert.New(t, false)
	buf := &bytes.Buffer{}
	l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout))
	a.NotNil(l)
	err1 := errors.New("error")

	t.Run("error", func(t *testing.T) {
		a := assert.New(t, false)
		e := l.NewRecord(LevelWarn)
		a.Empty(e.AppendLocation)
		e.Error(err1) // 输出定位
		a.NotEmpty(e.AppendLocation)

		b := NewBuffer(false)
		defer b.Free()
		e.AppendLocation(b)
		a.True(strings.HasSuffix(string(b.data), "record_test.go:61"), string(b.data))
	})

	err2 := &err{err: err1}
	buf.Reset()
	t.Run("xerrors>error", func(t *testing.T) {
		a := assert.New(t, false)
		e := l.NewRecord(LevelWarn)
		a.Empty(e.AppendLocation)
		e.Error(err2)
		a.True(strings.HasSuffix(buf.String(), "root\n\n"), buf.String())
	})

	l.detail = true
	buf.Reset()
	t.Run("xerrors>error", func(t *testing.T) {
		a := assert.New(t, false)
		e := l.NewRecord(LevelWarn)
		a.Empty(e.AppendLocation)
		e.Error(err2)
		a.True(strings.HasSuffix(buf.String(), "root\nerror\n"), buf.String())
	})

	buf.Reset()
	t.Run("2*xerrors>error", func(t *testing.T) {
		a := assert.New(t, false)
		e := l.NewRecord(LevelWarn)
		a.Empty(e.AppendLocation)
		e.Error(&err{err: err2})
		a.True(strings.HasSuffix(buf.String(), "root\nroot\nerror\n"), buf.String())
	})

	lerr1 := localeutil.Error("loc err")
	buf.Reset()
	t.Run("locale without catalog", func(t *testing.T) {
		a := assert.New(t, false)
		e := l.NewRecord(LevelWarn)
		a.Empty(e.AppendLocation)
		e.Error(lerr1)
		a.True(strings.HasSuffix(buf.String(), "loc err\n"), buf.String())
	})

	lerr2 := &err{err: lerr1}
	buf.Reset()
	t.Run("xerrors>locale without catalog", func(t *testing.T) {
		a := assert.New(t, false)
		e := l.NewRecord(LevelWarn)
		a.Empty(e.AppendLocation)
		e.Error(lerr2)
		a.True(strings.HasSuffix(buf.String(), "loc err\n"), buf.String())
	})

	c := catalog.NewBuilder()
	a.NotError(c.SetString(language.SimplifiedChinese, "loc err", "cn"))
	l.printer = message.NewPrinter(language.SimplifiedChinese, message.Catalog(c))

	buf.Reset()
	t.Run("locale with catalog", func(t *testing.T) {
		a := assert.New(t, false)
		e := l.NewRecord(LevelWarn)
		a.Empty(e.AppendLocation)
		e.Error(lerr1)
		a.True(strings.HasSuffix(buf.String(), "cn\n"), buf.String())
	})

	buf.Reset()
	t.Run("xerrors>locale with catalog", func(t *testing.T) {
		a := assert.New(t, false)
		e := l.NewRecord(LevelWarn)
		a.Empty(e.AppendLocation)
		e.Error(lerr2)
		a.True(strings.HasSuffix(buf.String(), "root\ncn\n"), buf.String())
	})

	buf.Reset()
	t.Run("2*xerrors>locale with catalog", func(t *testing.T) {
		a := assert.New(t, false)
		e := l.NewRecord(LevelWarn)
		a.Empty(e.AppendLocation)
		e.Error(&err{err: lerr2})
		a.True(strings.HasSuffix(buf.String(), "root\nroot\ncn\n"), buf.String())
	})
}

func TestRecord_Println(t *testing.T) {
	a := assert.New(t, false)

	c := catalog.NewBuilder()
	a.NotError(c.SetString(language.SimplifiedChinese, "abc", "cn"))
	p := message.NewPrinter(language.SimplifiedChinese, message.Catalog(c))
	a.NotNil(p)
	buf := &bytes.Buffer{}
	l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout), WithLocale(p))
	a.NotNil(l)

	e := l.NewRecord(LevelWarn)
	e.Println(localeutil.Phrase("abc"))
	a.True(strings.HasSuffix(buf.String(), "cn\n\n"), buf.String()) // Println 本身包含了一个回车符
}
