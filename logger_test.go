// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"log"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/logs/writers"
)

func TestLogger_set(t *testing.T) {
	a := assert.New(t)

	l := &logger{}
	l.set(nil, "", 0)
	a.Nil(l.flush).Equal(l.log, discard)

	cont := writers.NewContainer()
	l.set(cont, "", 0)
	a.Equal(cont, l.flush).NotEqual(l.log, discard)
}

func TestParseFlag(t *testing.T) {
	a := assert.New(t)

	eq := func(str string, v int) {
		ret, err := parseFlag(str)
		a.NotError(err).Equal(ret, v)
	}

	eq("log.Ldate|log.ltime", log.Ldate|log.Ltime)
	eq("log.Ldate| log.Ltime", log.Ldate|log.Ltime)
	eq(" ", 0)
	eq("", 0)
}
