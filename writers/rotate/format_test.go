// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rotate

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

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

func TestGetIndex(t *testing.T) {
	a := assert.New(t)
	now := time.Now()
	prefixValue := now.Format("2006.")
	suffixValue := now.Format(".01")

	i, err := getIndex("./testdata", prefixValue, suffixValue)
	a.NotError(err).Equal(i, 0)

	w := func(i int) {
		name := "./testdata/" + prefixValue + strconv.Itoa(i) + suffixValue
		ioutil.WriteFile(name, []byte("123"), os.ModePerm)
	}

	w(5)
	i, err = getIndex("./testdata", prefixValue, suffixValue)
	a.NotError(err).Equal(i, 5)

	w(8)
	i, err = getIndex("./testdata", prefixValue, suffixValue)
	a.NotError(err).Equal(i, 8)
}
