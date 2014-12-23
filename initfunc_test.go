// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"log"
	"testing"

	"github.com/issue9/assert"
)

func TestLogWriterInitializer(t *testing.T) {
	a := assert.New(t)

	args := map[string]string{
		"prefix": "[INFO]",
		"flag":   "log.lstdflags",
		"misc":   "misc",
	}
	w, err := logWriterInitializer(args)
	a.NotError(err).NotNil(w)

	lw, ok := w.(*logWriter)
	a.True(ok).NotNil(lw)

	a.Equal(lw.prefix, "[INFO]").
		Equal(lw.flag, log.LstdFlags)
}
