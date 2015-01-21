// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

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
    <info></info>
</logs>
`
	cfg, err := ParseXMLString(xmlCfg)
	a.NotError(err).NotNil(cfg)
}

func TestParseXML(t *testing.T) {
	var xmlCfg = `
<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <!-- comment -->
    <info>
        <buffer size="5">
            <file dir="/var/logs/info" />
        </buffer>
        <console color="\033[12m" />
    </info>

    <!-- 中文注释 -->
    <debug>
        <file dir="/var/logs/debug" />
        <console color="\033[12" />
    </debug>
</logs>
`
	a := assert.New(t)
	r := bytes.NewReader([]byte(xmlCfg))
	a.NotNil(r)

	cfg, err := parseXML(r)
	a.NotError(err).NotNil(cfg)
	a.Equal(2, len(cfg.Items)) // info debug

	info, found := cfg.Items["info"]
	a.True(found).NotNil(info).Equal(info.Name, "info")
	a.Equal(2, len(info.Items)) // buffer,console

	buf, found := info.Items["buffer"]
	a.True(found).NotNil(buf)
	a.Equal(buf.Attrs["size"], "5")

	// 测试错误的xml配置文件，重复info元素名
	xmlCfg = `
<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <info></info>
    <info></info>
</logs>
`
	r = bytes.NewReader([]byte(xmlCfg))
	a.NotNil(r)

	cfg, err = parseXML(r)
	a.Error(err).Nil(cfg)
}
