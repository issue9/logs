// SPDX-License-Identifier: MIT

package initfunc

import (
	"io"
	"strconv"

	"github.com/issue9/logs/v2/config"
	"github.com/issue9/logs/v2/writers"
)

// Buffer 是 writers.Buffer 的初始化函数
func Buffer(cfg *config.Config) (io.Writer, error) {
	size, found := cfg.Attrs["size"]
	if !found {
		return nil, argNotFoundErr("buffer", "size")
	}

	num, err := strconv.Atoi(size)
	if err != nil {
		return nil, err
	}

	return writers.NewBuffer(num), nil
}
