// SPDX-License-Identifier: MIT

package config

import (
	"os"
	"testing"

	"github.com/issue9/assert/v2"
)

func TestParseXMLString(t *testing.T) {
	a := assert.New(t, false)
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
	a := assert.New(t, false)

	cfg, err := ParseYAMLFile("./config.yml")
	a.NotError(err).NotNil(cfg)
	a.Equal(5, len(cfg.Items))

	erro, found := cfg.Items["error"]
	a.True(found).NotNil(erro)
	a.Equal(3, len(erro.Items))

	console, found := erro.Items["console"]
	a.True(found).NotNil(console)
	a.Equal(console.Attrs["output"], "stderr")

	cfg2, err := ParseYAMLFileFS(os.DirFS("./"), "config.yml")
	a.NotError(err).NotNil(cfg)
	a.Equal(cfg2, cfg)
}

func TestConfig_json(t *testing.T) {
	a := assert.New(t, false)

	cfg, err := ParseJSONFile("./config.json")
	a.NotError(err).NotNil(cfg)
	a.Equal(5, len(cfg.Items))

	erro, found := cfg.Items["error"]
	a.True(found).NotNil(erro)
	a.Equal(3, len(erro.Items))

	console, found := erro.Items["console"]
	a.True(found).NotNil(console)
	a.Equal(console.Attrs["output"], "stderr")

	cfg2, err := ParseJSONFileFS(os.DirFS("./"), "config.json")
	a.NotError(err).NotNil(cfg)
	a.Equal(cfg2, cfg)
}

func TestConfig_xml(t *testing.T) {
	a := assert.New(t, false)

	cfg, err := ParseXMLFile("./config.xml")
	a.NotError(err).NotNil(cfg)
	a.Equal(6, len(cfg.Items))

	erro, found := cfg.Items["error"]
	a.True(found).NotNil(erro)
	a.Equal(3, len(erro.Items))

	console, found := erro.Items["console"]
	a.True(found).NotNil(console)
	a.Equal(console.Attrs["output"], "stderr")

	cfg2, err := ParseXMLFileFS(os.DirFS("./"), "config.xml")
	a.NotError(err).NotNil(cfg)
	a.Equal(cfg2, cfg)
}
