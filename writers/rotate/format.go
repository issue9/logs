// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package rotate

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

const indexPlaceholder = "%i"

// ErrIndexNotExists 格式化字符串 %i 不存在
var ErrIndexNotExists = errors.New("必须存在 %i")

var dateRelpacer = strings.NewReplacer(
	"%y", "06",
	"%Y", "2006",
	"%m", "01",
	"%d", "02",
	"%h", "03",
	"%H", "15")

func parseFormat(format string) (prefix, suffix string, err error) {
	index := strings.Index(format, indexPlaceholder)
	if index < 0 {
		return "", "", ErrIndexNotExists
	}

	prefix = format[:index]
	suffix = format[index+len(indexPlaceholder):]

	return dateRelpacer.Replace(prefix), dateRelpacer.Replace(suffix), nil
}

// 获取指定目录下，去掉前后缀之后，最大的索引值。
func getIndex(dir, prefix, suffix string) (int, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return 0, err
	}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	var index int
	for _, f := range fs {
		name := f.Name()

		if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, suffix) {
			continue
		}

		istr := strings.TrimSuffix(strings.TrimPrefix(f.Name(), prefix), suffix)
		i, err := strconv.Atoi(istr)
		if err != nil {
			continue
		}

		if i > index {
			index = i
		}
	}

	return index, nil
}
