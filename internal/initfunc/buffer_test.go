// SPDX-License-Identifier: MIT

package initfunc

import (
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/issue9/logs/v3/config"
	"github.com/issue9/logs/v3/writers"
)

func TestBuffer(t *testing.T) {
	a := assert.New(t, false)
	cfg := &config.Config{Attrs: map[string]string{}}

	w, err := Buffer(cfg)
	a.Error(err).Nil(w)

	cfg.Attrs["size"] = "5"
	w, err = Buffer(cfg)
	a.NotError(err).NotNil(w)
	_, ok := w.(*writers.Buffer)
	a.True(ok)

	// 无法解析的 size 参数
	cfg.Attrs["size"] = "5l"
	w, err = Buffer(cfg)
	a.Error(err).Nil(w)
}
