// SPDX-License-Identifier: MIT

package writers

import (
	"bytes"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestNew(t *testing.T) {
	a := assert.New(t, false)

	a.PanicString(func() {
		New()
	}, "参数 w 不能为空")

	w1 := new(bytes.Buffer)
	w2 := new(bytes.Buffer)

	w := New(w1)
	a.Equal(w1, w)

	w = New(w1, w2)
	a.Length(w, 2)

	w.Write([]byte("123"))

	a.Equal(w1.String(), "123").
		Equal(w2.String(), "123")
}
