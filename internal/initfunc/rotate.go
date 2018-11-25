// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/issue9/logs/v2/config"
	"github.com/issue9/logs/v2/writers/rotate"
)

const (
	b int64 = 1 << (10 * iota)
	kb
	mb
	gb
)

// Rotate 的初始化函数。
func Rotate(cfg *config.Config) (io.Writer, error) {
	format, found := cfg.Attrs["filename"]
	if !found {
		return nil, argNotFoundErr("rotate", "filename")
	}

	dir, found := cfg.Attrs["dir"]
	if !found {
		return nil, argNotFoundErr("rotate", "dir")
	}

	sizeStr, found := cfg.Attrs["size"]
	if !found {
		return nil, argNotFoundErr("rotate", "size")
	}

	size, err := toByte(sizeStr)
	if err != nil {
		return nil, err
	}

	return rotate.New(format, dir, size)
}

// 将字符串转换成以字节为单位的数值。
// 粗略计算，并不 100% 正确，小数只取整数部分。
// 支持以下格式：
//  1024
//  1k
//  1M
//  1G
// 后缀单位只支持 k,g,m，不区分大小写。
func toByte(str string) (int64, error) {
	if len(str) == 0 {
		return -1, errors.New("不能传递空值")
	}

	str = strings.ToLower(str)

	scale := b
	unit := str[len(str)-1]
	switch {
	case unit >= '0' && unit <= '9':
		scale = b
	case unit == 'b':
		scale = b
	case unit == 'k':
		scale = kb
	case unit == 'm':
		scale = mb
	case unit == 'g':
		scale = gb
	default:
		return -1, fmt.Errorf("无法识别的单位:[%v]", unit)
	}

	if scale > 1 {
		str = str[:len(str)-1]
	}

	if len(str) == 0 {
		return -1, errors.New("传递了一个空值")
	}

	size, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return -1, err
	}

	if size <= 0 {
		return -1, fmt.Errorf("大小不能小于 0，当前值为:[%.4f]", size)
	}

	return int64(size) * scale, nil
}
