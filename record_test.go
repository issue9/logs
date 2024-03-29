// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/issue9/assert/v4"
	"github.com/issue9/localeutil"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"golang.org/x/xerrors"
)

var (
	_ Recorder = &withRecorder{}
	_ Recorder = &disableRecorder{}
	_ Recorder = &Logger{}
)

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
	l := New(NewTextHandler(io.Discard), WithLocation(true), WithCreated(MicroLayout))

	e := l.NewRecord()
	a.NotNil(e)
	a.Empty(e.AppendLocation)

	e.initLocationCreated(1) // 输出定位
	b := NewBuffer(false)
	defer b.Free()
	e.AppendLocation(b)
	a.True(strings.HasSuffix(string(b.data), "record_test.go:50"), string(b.data))
}

func TestRecord_Error(t *testing.T) {
	a := assert.New(t, false)
	err1 := errors.New("error")

	t.Run("error", func(*testing.T) {
		buf := &bytes.Buffer{}
		l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout))
		a.NotNil(l)
		l.WARN().Error(err1) // 输出定位
		a.True(strings.Contains(buf.String(), "record_test.go:65"), buf.String())
	})

	err2 := &err{err: err1}
	t.Run("xerrors>error", func(*testing.T) {
		buf := &bytes.Buffer{}
		l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout))
		l.WARN().Error(err2)
		a.True(strings.HasSuffix(buf.String(), "root\n\n"), buf.String())
	})

	t.Run("xerrors>error", func(*testing.T) {
		buf := &bytes.Buffer{}
		l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout), WithDetail(true))
		l.WARN().Error(err2)
		a.True(strings.HasSuffix(buf.String(), "root\nerror\n"), buf.String())
	})

	t.Run("2*xerrors>error", func(*testing.T) {
		buf := &bytes.Buffer{}
		l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout), WithDetail(true))
		l.WARN().Error(&err{err: err2})
		a.True(strings.HasSuffix(buf.String(), "root\nroot\nerror\n"), buf.String())
	})

	lerr1 := localeutil.Error("loc err")
	t.Run("locale without catalog", func(*testing.T) {
		buf := &bytes.Buffer{}
		l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout), WithDetail(true))
		l.WARN().Error(lerr1)
		a.True(strings.HasSuffix(buf.String(), "loc err\n"), buf.String())
	})

	lerr2 := &err{err: lerr1}
	t.Run("xerrors>locale without catalog", func(*testing.T) {
		buf := &bytes.Buffer{}
		l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout), WithDetail(true))
		l.WARN().Error(lerr2)
		a.True(strings.HasSuffix(buf.String(), "loc err\n"), buf.String())
	})

	t.Run("locale with catalog", func(*testing.T) {
		c := catalog.NewBuilder()
		a.NotError(c.SetString(language.SimplifiedChinese, "loc err", "cn"))
		p := message.NewPrinter(language.SimplifiedChinese, message.Catalog(c))

		buf := &bytes.Buffer{}
		l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout), WithDetail(true), WithLocale(p))
		l.WARN().Error(lerr1)
		a.True(strings.HasSuffix(buf.String(), "cn\n"), buf.String())
	})

	t.Run("xerrors>locale with catalog", func(*testing.T) {
		c := catalog.NewBuilder()
		a.NotError(c.SetString(language.SimplifiedChinese, "loc err", "cn"))
		p := message.NewPrinter(language.SimplifiedChinese, message.Catalog(c))

		buf := &bytes.Buffer{}
		l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout), WithDetail(true), WithLocale(p))
		l.WARN().Error(lerr2)
		a.True(strings.HasSuffix(buf.String(), "root\ncn\n"), buf.String())
	})

	t.Run("2*xerrors>locale with catalog", func(*testing.T) {
		c := catalog.NewBuilder()
		a.NotError(c.SetString(language.SimplifiedChinese, "loc err", "cn"))
		p := message.NewPrinter(language.SimplifiedChinese, message.Catalog(c))

		buf := &bytes.Buffer{}
		l := New(NewTextHandler(buf), WithLocation(true), WithCreated(MicroLayout), WithDetail(true), WithLocale(p))
		l.WARN().Error(&err{err: lerr2})
		a.True(strings.HasSuffix(buf.String(), "root\nroot\ncn\n"), buf.String())
	})
}
