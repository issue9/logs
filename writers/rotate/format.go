// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rotate

import (
	"errors"
	"strings"
)

var ErrIndexNotExists = errors.New("必须存在 %i")

var dateRelpacer = strings.NewReplacer("%y", "06",
	"%Y", "2006",
	"%m", "01",
	"%d", "02",
	"%h", "03",
	"%H", "15")

const indexPlaceholder = "%i"

func parseFormat(format string) (prefix, suffix string, err error) {
	index := strings.Index(format, indexPlaceholder)
	if index < 0 {
		return "", "", ErrIndexNotExists
	}

	prefix = format[:index]
	suffix = format[index+len(indexPlaceholder):]

	return dateRelpacer.Replace(prefix), dateRelpacer.Replace(suffix), nil
}
