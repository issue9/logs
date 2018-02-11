// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rotate

import (
	"testing"

	"github.com/issue9/assert"
)

func TestParseFormat(t *testing.T) {
	a := assert.New(t)

	p, s, err := parseFormat("")
	a.Empty(p).Empty(s).Equal(err, ErrIndexNotExists)

	p, s, err = parseFormat("%i")
	a.NotError(err).
		Empty(p).
		Empty(s)

	p, s, err = parseFormat("test%i")
	a.NotError(err).
		Equal(p, "test").
		Empty(s)

	p, s, err = parseFormat("test-%Y%d%i%yy%m-%H")
	a.NotError(err).
		Equal(p, "test-200602").
		Equal(s, "06y01-15")
}
