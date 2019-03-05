// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/issue9/assert"
)

func TestParseXMLString(t *testing.T) {
	a := assert.New(t)
	xmlCfg := `
<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <info>
		<console foreground="red" />
	</info>
</logs>
`
	cfg, err := ParseXMLString(xmlCfg)
	a.NotError(err).NotNil(cfg)
}

func TestConfig_yaml(t *testing.T) {
	a := assert.New(t)

	cfg, err := ParseYAMLFile("./config.yml")
	a.NotError(err).NotNil(cfg)
	a.Equal(5, len(cfg.Items))

	erro, found := cfg.Items["error"]
	a.True(found).NotNil(erro)
	a.Equal(3, len(erro.Items))

	console, found := erro.Items["console"]
	a.True(found).NotNil(console)
	a.Equal(console.Attrs["output"], "stderr")
}

func TestConfig_json(t *testing.T) {
	a := assert.New(t)

	cfg, err := ParseJSONFile("./config.json")
	a.NotError(err).NotNil(cfg)
	a.Equal(5, len(cfg.Items))

	erro, found := cfg.Items["error"]
	a.True(found).NotNil(erro)
	a.Equal(3, len(erro.Items))

	console, found := erro.Items["console"]
	a.True(found).NotNil(console)
	a.Equal(console.Attrs["output"], "stderr")
}

func TestConfig_xml(t *testing.T) {
	a := assert.New(t)

	cfg, err := ParseXMLFile("./config.xml")
	a.NotError(err).NotNil(cfg)
	a.Equal(6, len(cfg.Items))

	erro, found := cfg.Items["error"]
	a.True(found).NotNil(erro)
	a.Equal(3, len(erro.Items))

	console, found := erro.Items["console"]
	a.True(found).NotNil(console)
	a.Equal(console.Attrs["output"], "stderr")
}
