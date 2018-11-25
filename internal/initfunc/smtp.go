// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"io"
	"strings"

	"github.com/issue9/logs/v2/config"
	"github.com/issue9/logs/v2/writers"
)

// SMTP 是 writers.SMTP 的初始化函数
func SMTP(cfg *config.Config) (io.Writer, error) {
	username, found := cfg.Attrs["username"]
	if !found {
		return nil, argNotFoundErr("stmp", "username")
	}

	password, found := cfg.Attrs["password"]
	if !found {
		return nil, argNotFoundErr("stmp", "password")
	}

	subject, found := cfg.Attrs["subject"]
	if !found {
		return nil, argNotFoundErr("stmp", "subject")
	}

	host, found := cfg.Attrs["host"]
	if !found {
		return nil, argNotFoundErr("stmp", "host")
	}

	sendToStr, found := cfg.Attrs["sendTo"]
	if !found {
		return nil, argNotFoundErr("stmp", "sendTo")
	}

	sendTo := strings.Split(sendToStr, ";")

	return writers.NewSMTP(username, password, subject, host, sendTo), nil
}
