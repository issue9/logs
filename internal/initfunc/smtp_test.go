// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/logs/writers"
)

func TestSMTP(t *testing.T) {
	a := assert.New(t)
	args := map[string]string{}
	w, err := SMTP(args)
	a.Error(err).Nil(w)

	args["username"] = "abc"
	args["password"] = "abc"
	args["subject"] = "subject"
	args["host"] = "host"
	args["sendTo"] = "sendTo"
	w, err = SMTP(args)
	a.NotError(err).NotNil(w)

	_, ok := w.(*writers.SMTP)
	a.True(ok)
}
