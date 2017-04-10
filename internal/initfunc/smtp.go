// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"io"
	"strings"

	"github.com/issue9/logs/writers"
)

// SMTP 是 writers.SMTP 的初始化函数
func SMTP(args map[string]string) (io.Writer, error) {
	username, found := args["username"]
	if !found {
		return nil, argNotFoundErr("stmp", "username")
	}

	password, found := args["password"]
	if !found {
		return nil, argNotFoundErr("stmp", "password")
	}

	subject, found := args["subject"]
	if !found {
		return nil, argNotFoundErr("stmp", "subject")
	}

	host, found := args["host"]
	if !found {
		return nil, argNotFoundErr("stmp", "host")
	}

	sendToStr, found := args["sendTo"]
	if !found {
		return nil, argNotFoundErr("stmp", "sendTo")
	}

	sendTo := strings.Split(sendToStr, ";")

	return writers.NewSMTP(username, password, subject, host, sendTo), nil
}
