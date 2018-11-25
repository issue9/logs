// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/issue9/logs/v2/config"
	"github.com/issue9/logs/v2/writers"
)

func TestSMTP(t *testing.T) {
	a := assert.New(t)
	cfg := &config.Config{Attrs: map[string]string{}}
	w, err := SMTP(cfg)
	a.Error(err).Nil(w)

	cfg.Attrs["username"] = "abc"
	cfg.Attrs["password"] = "abc"
	cfg.Attrs["subject"] = "subject"
	cfg.Attrs["host"] = "host"
	cfg.Attrs["sendTo"] = "sendTo"
	w, err = SMTP(cfg)
	a.NotError(err).NotNil(w)

	_, ok := w.(*writers.SMTP)
	a.True(ok)
}
