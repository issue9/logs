// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"log"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/logs/writers"
)

func TestToByte(t *testing.T) {
	a := assert.New(t)

	eq := func(str string, val int64) {
		size, err := toByte(str)
		a.NotError(err).Equal(size, val)
	}

	e := func(str string) {
		size, err := toByte(str)
		a.Error(err).Equal(size, -1)
	}

	eq("1m", 1024*1024)
	eq("100G", 100*1024*1024*1024)
	eq("10.2k", 10*1024)
	eq("10.9K", 10*1024)

	e("")
	e("M")
	e("-1M")
	e("-1.0G")
	e("1P")
	e("10MB")
}

func TestParseFlag(t *testing.T) {
	a := assert.New(t)

	eq := func(str string, v int) {
		ret, err := parseFlag(str)
		a.NotError(err).Equal(ret, v)
	}

	eq("log.Ldate|log.ltime", log.Ldate|log.Ltime)
	eq("log.Ldate| log.Ltime", log.Ldate|log.Ltime)
}

func TestRotateInitializer(t *testing.T) {
	a := assert.New(t)
	args := map[string]string{}

	w, err := rotateInitializer(args)
	a.Error(err).Nil(w)

	// 缺少size
	args["dir"] = "./testdata"
	w, err = rotateInitializer(args)
	a.Error(err).Nil(w)

	// 错误的size参数
	args["size"] = "12P"
	w, err = rotateInitializer(args)
	a.Error(err).Nil(w)

	// 正常
	args["size"] = "12"
	w, err = rotateInitializer(args)
	a.NotError(err).NotNil(w)

	_, ok := w.(*writers.Rotate)
	a.True(ok)
}

func TestBufferInitializer(t *testing.T) {
	a := assert.New(t)
	args := map[string]string{}

	w, err := bufferInitializer(args)
	a.Error(err).Nil(w)

	args["size"] = "5"
	w, err = bufferInitializer(args)
	a.NotError(err).NotNil(w)
	_, ok := w.(*writers.Buffer)
	a.True(ok)

	// 无法解析的size参数
	args["size"] = "5l"
	w, err = bufferInitializer(args)
	a.Error(err).Nil(w)
}

func TestConsoleInitializer(t *testing.T) {
	a := assert.New(t)
	args := map[string]string{}

	// 可以接受空参数，consoleInitializer的args都有默认值
	w, err := consoleInitializer(args)
	a.NotError(err).NotNil(w)

	// 无效的output
	args["output"] = "file"
	w, err = consoleInitializer(args)
	a.Error(err).Nil(w)
	args["output"] = "stderr"

	// 无效的foreground
	args["foreground"] = "red1"
	w, err = consoleInitializer(args)
	a.Error(err).Nil(w)
	args["foreground"] = "red"

	// 无效的background
	args["background"] = "red1"
	w, err = consoleInitializer(args)
	a.Error(err).Nil(w)

	args["background"] = "blue"
	w, err = consoleInitializer(args)
	a.NotError(err).NotNil(w)

	_, ok := w.(*writers.Console)
	a.True(ok)
}

func TestStmpInitializer(t *testing.T) {
	a := assert.New(t)
	args := map[string]string{}
	w, err := stmpInitializer(args)
	a.Error(err).Nil(w)

	args["username"] = "abc"
	args["password"] = "abc"
	args["subject"] = "subject"
	args["host"] = "host"
	args["sendTo"] = "sendTo"
	w, err = stmpInitializer(args)
	a.NotError(err).NotNil(w)

	_, ok := w.(*writers.SMTP)
	a.True(ok)
}
