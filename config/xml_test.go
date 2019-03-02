// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/config"
)

var _ config.UnmarshalFunc = XMLUnmarshal

func TestParseXMLFile(t *testing.T) {
	a := assert.New(t)
	cfg, err := ParseXMLFile("config.xml")
	a.NotError(err).NotNil(cfg)
}

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

func TestXMLUnmarshal(t *testing.T) {
	var xmlCfg = []byte(`
<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <!-- comment -->
    <info>
        <buffer size="5">
            <file dir="/var/logs/info" />
        </buffer>
        <console foreground="red" />
    </info>

    <!-- 中文注释 -->
    <debug>
        <file dir="/var/logs/debug" />
        <console foreground="red" />
    </debug>
</logs>
`)
	a := assert.New(t)

	cfg := &Config{}
	a.NotError(XMLUnmarshal(xmlCfg, cfg))
	a.Equal(2, len(cfg.Items)) // info debug

	info, found := cfg.Items["info"]
	a.True(found).NotNil(info).Equal(info.Name, "info")
	a.Equal(2, len(info.Items)) // buffer,console

	buf, found := info.Items["buffer"]
	a.True(found).NotNil(buf)
	a.Equal(buf.Attrs["size"], "5")

	// 测试错误的 xml 配置文件，重复 info 元素名
	xmlCfg = []byte(`
<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <info></info>
    <info></info>
</logs>
`)

	cfg = &Config{}
	a.Error(XMLUnmarshal(xmlCfg, cfg))
}
