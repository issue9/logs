// SPDX-License-Identifier: MIT

package writers

import (
	"io"
	"os"
	"testing"

	"github.com/issue9/term/v2/colors"
)

var _ io.Writer = &Console{}

func TestConsole(t *testing.T) {
	c := NewConsole(os.Stderr, colors.Cyan, colors.Default)
	c.Write([]byte("is cyan\n"))

	c.SetColor(colors.Blue, colors.Default)
	c.Write([]byte("is blue\n"))

	os.Stderr.WriteString("Reset\n")
}
