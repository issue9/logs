// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

/*
import (
	"testing"

	"github.com/issue9/assert"
)

func TestConfig_Sanitize(t *testing.T) {
	a := assert.New(t)

	// 错误的xml内容
	xml := `
<?xml version="1.0" encoding="utf-8"?>
<logs>
</logs>
`
	cfg, err := ParseXMLString(xml)
	a.Error(err).Nil(cfg)

	// 错误的xml内容,顶级只能为logs
	xml = `
<?xml version="1.0" encoding="utf-8"?>
<log>
	<debug></debug>
</log>
`
	cfg, err = ParseXMLString(xml)
	a.Error(err).Nil(cfg)

	// 错误的xml内容,未知的debug1元素
	xml = `
<?xml version="1.0" encoding="utf-8"?>
<logs>
	<debug1></debug1>
</logs>
`
	cfg, err = ParseXMLString(xml)
	a.Error(err).Nil(cfg)

	// 错误的xml内容,debug必须要有子元素。
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

*/
