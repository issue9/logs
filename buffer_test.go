// SPDX-License-Identifier: MIT

package logs

import (
	"io"

	"golang.org/x/xerrors"
)

var (
	_ xerrors.Printer = NewBuffer()
	_ io.Writer       = NewBuffer()
)
