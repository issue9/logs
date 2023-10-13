// SPDX-License-Identifier: MIT

package logs

import "golang.org/x/xerrors"

var _ xerrors.Printer = NewBuffer()
