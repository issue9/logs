// SPDX-License-Identifier: MIT

package logs

import (
	"encoding"
	"fmt"
)

var (
	_ encoding.TextMarshaler = LevelDebug
	_ fmt.Stringer           = LevelError
)
