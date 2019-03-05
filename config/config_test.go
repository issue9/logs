// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/config"
)

var (
	_ config.Sanitizer = &Config{}
	_ xml.Unmarshaler  = &Config{}
)

func TestConfig_sanitize(t *testing.T) {
	a := assert.New(t)

	// 错误的 xml 内容
	xml := `
<?xml version="1.0" encoding="utf-8"?>
<logs>
</logs>
`
	cfg, err := ParseXMLString(xml)
	a.Error(err).Nil(cfg)

	// 错误的 xml 内容,顶级只能为 logs
	xml = `
<?xml version="1.0" encoding="utf-8"?>
<log>
	<debug></debug>
</log>
`
	cfg, err = ParseXMLString(xml)
	a.Error(err).Nil(cfg)

	// 错误的 xml 内容,未知的 debug1 元素
	xml = `
<?xml version="1.0" encoding="utf-8"?>
<logs>
	<debug1></debug1>
</logs>
`
	cfg, err = ParseXMLString(xml)
	a.Error(err).Nil(cfg)

	// 错误的 xml 内容,debug 必须要有子元素。
	xml = `
<?xml version="1.0" encoding="utf-8"?>
<logs>
    <debug>
    </debug>
</logs>
`
	cfg, err = ParseXMLString(xml)
	a.Error(err).Nil(cfg)

	// 正确内容
	xml = `
<?xml version="1.0" encoding="utf-8"?>
<logs>
    <debug>
		<buffer size="10"></buffer>
    </debug>
</logs>
`
	cfg, err = ParseXMLString(xml)
	a.NotError(err).NotNil(cfg)
}
