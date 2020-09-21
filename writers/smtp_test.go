// SPDX-License-Identifier: MIT

package writers

import (
	"io"
	"testing"
	"time"

	"github.com/issue9/assert"
)

var _ io.Writer = &SMTP{}

func testSMTP(t *testing.T) {
	smtp := NewSMTP("test@qq.com", "pwd", "test", "smtp.qq.com:25", []string{"test@gmail.com"})

	size, err := smtp.Write([]byte("test"))
	assert.NotError(t, err)
	assert.True(t, size > 0)

	time.Sleep(30 * time.Second)

	size, err = smtp.Write([]byte("test2"))
	assert.NotError(t, err)
	assert.True(t, size > 0)
}
