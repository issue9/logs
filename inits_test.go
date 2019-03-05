// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"io"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/logs/v2/config"
	"github.com/issue9/logs/v2/writers"
)

// configTestWriter.Write 的输入内容，写到此变量中
var configTestWriterContent []byte

type configTestWriter struct {
	ws []io.Writer
}

// 清除已经注册的初始化函数。
func clearInitializer() {
	funs = make(map[string]WriterInitializer)
}

func (t *configTestWriter) Write(bs []byte) (int, error) {
	configTestWriterContent = append(configTestWriterContent, bs...)
	return len(bs), nil
}

// 容器类初始化函数
func logsInit(cfg *config.Config) (io.Writer, error) {
	return writers.NewContainer(), nil
}

func debugInit(cfg *config.Config) (io.Writer, error) {
	return &configTestWriter{}, nil
}

func TestToWriter(t *testing.T) {
	a := assert.New(t)
	clearInitializer()

	a.True(Register("logs", logsInit))
	a.True(Register("debug", debugInit))

	// 构造一个类似于以下结构的 config.Config
	// 不使用 config.ParseXML，可以躲避错误检测
	// <logs>
	//     <debug></debug>
	// </logs>
	cfg := &config.Config{
		Items: map[string]*config.Config{
			"debug": &config.Config{},
		},
	}

	// 转换成 writer
	w, err := toWriter("logs", cfg)
	a.NotError(err).NotNil(w)

	// 转换成 writers.Container
	c, ok := w.(*writers.Container)
	a.True(ok).NotNil(c)

	// 写入 c，应该有内容输出到 configTestWriterContent
	c.Write([]byte("hello"))
	a.Equal(configTestWriterContent, []byte("hello"))

	c.Write([]byte(" world"))
	a.Equal(configTestWriterContent, []byte("hello world"))

	// 未注册的初始化函数
	w, err = toWriter("unregister", cfg)
	a.Error(err).Nil(w)
}

func TestInits(t *testing.T) {
	a := assert.New(t)

	// 清空，包的 init 函数有可能会初始化一些数据。
	clearInitializer()

	a.True(Register("init1", logsInit)).
		True(IsRegisted("init1")).
		Equal(Registed(), []string{"init1"})

	a.True(Register("init2", debugInit)).
		True(IsRegisted("init2")).
		True(IsRegisted("init1")).
		Equal(2, len(Registed()))

	a.False(IsRegisted("init3"))

	// 重复注册
	a.False(Register("init1", debugInit))
	a.True(IsRegisted("init1"))

	clearInitializer()
	a.Equal(0, len(funs))
}
