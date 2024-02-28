// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package writers

import (
	"bytes"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	a.Panic(func() {
		New()
	})

	a.Panic(func() {
		New(nil)
	})

	b1 := &bytes.Buffer{}
	b2 := &bytes.Buffer{}
	w := New(b1, nil, b2)
	a.NotNil(w)
	w.Write([]byte("123"))
	a.Equal(b1.String(), b2.String()).Equal(b1.String(), "123")
}
