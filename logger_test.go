// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"log"
	"testing"

	"github.com/issue9/assert"
)

/*
func TestLogger_set(t *testing.T) {
	a := assert.New(t)

	l := &logger{
		log: log.New(ioutil.Discard, "", 0),
	}
	l.set(nil, "", 0)
	a.Nil(l.flush)

	cont := writers.NewContainer()
	l.set(cont, "", 0)
	a.Equal(cont, l.flush)

	l.set(cont, "abc", 2)
	a.Equal(cont, l.flush)
}
*/

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
