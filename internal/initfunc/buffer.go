// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"io"
	"strconv"

	"github.com/issue9/logs/writers"
)

// Buffer 是 writers.Buffer 的初始化函数
func Buffer(args map[string]string) (io.Writer, error) {
	size, found := args["size"]
	if !found {
		return nil, argNotFoundErr("buffer", "size")
	}

	num, err := strconv.Atoi(size)
	if err != nil {
		return nil, err
	}

	return writers.NewBuffer(num), nil
}
