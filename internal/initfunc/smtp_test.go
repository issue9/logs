// SPDX-License-Identifier: MIT

package initfunc

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/issue9/logs/v3/config"
	"github.com/issue9/logs/v3/writers"
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
