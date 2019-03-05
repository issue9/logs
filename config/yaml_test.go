// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"io/ioutil"
	"testing"

	"github.com/issue9/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestConfig_MarshalYAML(t *testing.T) {
	a := assert.New(t)

	data, err := ioutil.ReadFile("./config.yml")
	a.NotError(err).NotEmpty(data)

	cfg := &Config{}
	a.NotError(yaml.Unmarshal(data, cfg))
	a.Equal(5, len(cfg.Items))

	erro, found := cfg.Items["error"]
	a.True(found).NotNil(erro)
	a.Equal(3, len(erro.Items))

	console, found := erro.Items["console"]
	a.True(found).NotNil(console)
	a.Equal(console.Attrs["output"], "stderr")
}
