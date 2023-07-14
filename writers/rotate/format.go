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

func cutString(format string) (prefix, suffix string, err error) {
	if strings.ContainsAny(format, "/\\") {
		return "", "", errors.New("不能包含路径分隔符 / 或 \\")
	}

	prefix, suffix, ok := strings.Cut(format, indexPlaceholder)
	if !ok {
		return "", "", ErrIndexNotExists
	}
	return prefix, suffix, nil
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
