// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package logs

import (
	"io"

	"golang.org/x/xerrors"
)

var (
	_ xerrors.Printer = NewBuffer(false)
	_ io.Writer       = NewBuffer(false)
)
