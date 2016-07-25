// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writers

import (
	"io"
	"os"
	"testing"

	"github.com/issue9/term/colors"
)

var _ io.Writer = &Console{}

func TestConsole(t *testing.T) {
	c := NewConsole(os.Stderr, colors.Cyan, colors.Default)
	c.Write([]byte("is cyan\n"))

	c.SetColor(colors.Blue, colors.Default)
	c.Write([]byte("is blue\n"))

	os.Stderr.WriteString("Reset\n")
}
