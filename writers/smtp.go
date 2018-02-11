// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writers

import (
	"bytes"
	"net/smtp"
	"strings"
)

// SMTP 实现 io.Writer 接口的邮件发送。
type SMTP struct {
	username string   // smtp账号
	password string   // smtp密码
	host     string   // smtp主机，需要带上端口
	sendTo   []string // 接收者。
	subject  string   // 邮件主题。

	// 邮件内容的缓存
	cache *bytes.Buffer
	// 邮件头部分的长度
	headerLen int

	auth smtp.Auth
}

// NewSMTP 新建 SMTP 对象。
// username 为smtp 的账号；
// password 为 smtp 对应的密码；
// subject 为发送邮件的主题；
// host 为 smtp 的主机地址，需要带上端口号；
// sendTo 为接收者的地址。
func NewSMTP(username, password, subject, host string, sendTo []string) *SMTP {
	ret := &SMTP{
		username: username,
		password: password,
		subject:  subject,
		host:     host,
		sendTo:   sendTo,
	}
	ret.init()

	return ret
}

// 初始化一些基本内容。
//
// 像To,From这些内容都是固定的，可以先写入到缓存中，这样
// 这后就不需要再次构造这些内容。
func (s *SMTP) init() {
	s.cache = bytes.NewBufferString("")
	s.cache.Grow(1024)

	// to
	s.cache.WriteString("To: ")
	s.cache.WriteString(strings.Join(s.sendTo, ";"))
	s.cache.WriteString("\r\n")

	// from
	s.cache.WriteString("From: ")
	s.cache.WriteString(s.username) // <...>有需要吗？
	s.cache.WriteString("\r\n")

	// subject
	s.cache.WriteString("Subject: ")
	s.cache.WriteString(s.subject)
	s.cache.WriteString("\r\n")

	// mime-version
	s.cache.WriteString("MIME-Version: ")
	s.cache.WriteString("1.0\r\n")

	// contentType
	s.cache.WriteString(`Content-Type: text/plain; charset="utf-8"`)
	s.cache.WriteString("\r\n\r\n")

	s.headerLen = s.cache.Len()

	// 去掉端口部分
	h := strings.Split(s.host, ":")[0]
	s.auth = smtp.PlainAuth("", s.username, s.password, h)
}

func (s *SMTP) Write(msg []byte) (int, error) {
	s.cache.Write(msg)

	err := smtp.SendMail(
		s.host,
		s.auth,
		s.username,
		s.sendTo,
		s.cache.Bytes(),
	)
	l := s.cache.Len()

	s.cache.Truncate(s.headerLen)

	return l, err
}
