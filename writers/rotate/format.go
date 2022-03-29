// SPDX-License-Identifier: MIT

package rotate

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const indexPlaceholder = "%i"

// ErrIndexNotExists 格式化字符串 %i 不存在
var ErrIndexNotExists = errors.New("必须存在 %i")

// TODO(go1.18): 可以采用 strings.Cut 代替
func cutString(format string) (prefix, suffix string, err error) {
	if index := strings.Index(format, indexPlaceholder); index >= 0 {
		return format[:index], format[index+len(indexPlaceholder):], nil
	}
	return "", "", ErrIndexNotExists
}

// 获取指定目录下，去掉前后缀之后，最大的索引值。
func getIndex(dir, prefix, suffix string) (int, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return 0, err
	}

	fs, err := os.ReadDir(dir)
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
