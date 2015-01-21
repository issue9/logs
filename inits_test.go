// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"io"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/logs/config"
	"github.com/issue9/logs/writer"
)

// test config.toWriter

// configTestWriter.Write的输入内容，写到此变量中
var configTestWriterContent []byte

type configTestWriter struct {
	ws []io.Writer
}

func (t *configTestWriter) Write(bs []byte) (int, error) {
	configTestWriterContent = append(configTestWriterContent, bs...)
	return len(bs), nil
}

func (t *configTestWriter) AddWriter(w io.Writer) error {
	t.ws = append(t.ws, w)
	return nil
}

var xmlCfg = `
<?xml version="1.0" encoding="utf-8" ?>
<logs>
    <debug>
    </debug>
</logs>
`

func TestToWriter(t *testing.T) {
	a := assert.New(t)
	clearInitializer()

	logs := func(args map[string]string) (io.Writer, error) {
		return writer.NewContainer(), nil
	}
	a.True(Register("logs", logs))

	debug := func(args map[string]string) (io.Writer, error) {
		return &configTestWriter{}, nil
	}
	a.True(Register("debug", debug))

	cfg, err := config.ParseXMLString(xmlCfg)
	a.NotError(err).NotNil(cfg)

	// 转换成writer
	w, err := toWriter(cfg)
	a.NotError(err).NotNil(w)

	// 转换成writer.Container
	c, ok := w.(*writer.Container)
	a.True(ok).NotNil(c)

	// 写入c，应该有内容输出到configTestWriterContent
	c.Write([]byte("hello"))
	a.Equal(configTestWriterContent, []byte("hello"))

	c.Write([]byte(" world"))
	a.Equal(configTestWriterContent, []byte("hello world"))
}

func init1(a1 map[string]string) (io.Writer, error) {
	return nil, nil
}

func init2(a1 map[string]string) (io.Writer, error) {
	return nil, nil
}

func TestInit(t *testing.T) {
	a := assert.New(t)

	// 清空，包的init函数有可能会初始化一些数据。
	clearInitializer()

	a.True(Register("init1", init1)).
		True(IsRegisted("init1")).
		Equal(Registed(), []string{"init1"})

	a.True(Register("init2", init2)).
		True(IsRegisted("init2")).
		True(IsRegisted("init1")).
		Equal(Registed(), []string{"init1", "init2"})

	a.False(IsRegisted("init3"))

	a.False(Register("init1", init2)) // 重复注册
	a.True(IsRegisted("init1"))

	clearInitializer()
	a.Equal(0, len(inits.names)).Equal(0, len(inits.funs))
}