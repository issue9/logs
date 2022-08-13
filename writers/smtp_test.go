// SPDX-License-Identifier: MIT

package writers

import (
	"io"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
)

var _ io.Writer = &SMTP{}

func testSMTP(t *testing.T) {
	a := assert.New(t, false)
	smtp := NewSMTP("test@qq.com", "pwd", "test", "smtp.qq.com:25", []string{"test@gmail.com"})

	size, err := smtp.Write([]byte("test"))
	a.NotError(err)
	a.True(size > 0)

	time.Sleep(30 * time.Second)

	size, err = smtp.Write([]byte("test2"))
	a.NotError(err)
	a.True(size > 0)
}
