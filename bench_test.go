// SPDX-License-Identifier: MIT

package logs

import (
	"testing"

	"github.com/issue9/assert/v2"
)

func BenchmarkEntry_Printf(b *testing.B) {
	a := assert.New(b, false)
	l := New()
	a.NotNil(l)

	err := l.ERROR()
	for i := 0; i < b.N; i++ {
		err.Value("k1", "v1").Printf("p1")
	}
}
