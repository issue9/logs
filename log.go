// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

//go:build !go1.21

package logs

import (
	"log"

	"github.com/issue9/logs/v7/writers"
)

// TODO(go1.21): go1.21 之后可删除，slog 默认会接管 log
func withStd(l *Logs) {
	log.SetOutput(writers.WriteFunc(func(b []byte) (int, error) {
		l.INFO().String(string(b)) // 参考 slog 默认输出至 INFO
		return len(b), nil
	}))
}
