// SPDX-License-Identifier: MIT

package logs

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/issue9/assert/v3"
	"golang.org/x/xerrors"
)

var _ Logger = &Record{}

type err struct {
	err error
}

func (e *err) Error() string { return e.err.Error() }

func (e *err) FormatError(p xerrors.Printer) error {
	p.Print("root\n")
	return e.err
}

func TestRecord_location(t *testing.T) {
	a := assert.New(t, false)
	l := New(nil, WithCaller(), WithCreated(MicroLayout))

	e := l.NewRecord(LevelWarn)
	a.NotNil(e)
	a.Empty(e.Path)

	e.setLocation(1)
	a.True(strings.HasSuffix(e.Path, "record_test.go:36"), e.Path) // 上一行
}

func TestRecord_Error(t *testing.T) {
	a := assert.New(t, false)
	buf := &bytes.Buffer{}
	l := New(NewTextHandler(buf), WithCaller(), WithCreated(MicroLayout))
	err1 := errors.New("error")

	e := l.NewRecord(LevelWarn)
	a.NotNil(e)
	a.Empty(e.Path)
	e.Error(err1)
	a.True(strings.HasSuffix(e.Path, "record_test.go:49"), e.Path) // 依赖上一行

	buf.Reset()
	err2 := &err{err: err1}
	e = l.NewRecord(LevelWarn)
	a.NotNil(e)
	a.Empty(e.Path)
	e.Error(err2)
	a.True(strings.HasSuffix(e.Path, "record_test.go:57"), e.Path) // 依赖上一行
	a.True(strings.HasSuffix(buf.String(), "root\nerror\n"), buf.String())

	buf.Reset()
	e = l.NewRecord(LevelWarn)
	a.NotNil(e)
	a.Empty(e.Path)
	e.Error(&err{err: err2})
	a.True(strings.HasSuffix(e.Path, "record_test.go:65"), e.Path) // 依赖上一行
	a.True(strings.HasSuffix(buf.String(), "root\nroot\nerror\n"), buf.String())
}
