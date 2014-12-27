// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"testing"

	"github.com/issue9/assert"
)

func TestLogs(t *testing.T) {
	a := assert.New(t)
	xml := `
<?xml version="1.0" encoding="utf-8"?>
<logs>
</logs>
`
	// 错误的xml内容
	a.Error(InitFromXml(xml))

	xml = `
<?xml version="1.0" encoding="utf-8"?>
<logs>
    <debug>
    </debug>
</logs>
`
	// 错误的xml内容
	a.Error(InitFromXml(xml))
}
