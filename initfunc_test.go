// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"log"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/logs/writer"
)

func TestRotateInitializer(t *testing.T) {
	a := assert.New(t)
	args := map[string]string{}

	w, err := rotateInitializer(args)
	a.Error(err).Nil(w)

	args["size"] = "12"
	args["dir"] = "c:/"
	w, err = rotateInitializer(args)
	a.NotError(err).NotNil(w)

	_, ok := w.(*writer.Rotate)
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
	_, ok := w.(*writer.Buffer)
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

	args["output"] = "stdin"
	w, err = consoleInitializer(args)
	a.Error(err).Nil(w)

	args["output"] = "stderr"
	args["foreground"] = "red"
	args["background"] = "blue"
	w, err = consoleInitializer(args)
	a.NotError(err).NotNil(w)

	_, ok := w.(*writer.Console)
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

	_, ok := w.(*writer.Smtp)
	a.True(ok)
}

func TestLogWriterInitializer(t *testing.T) {
	a := assert.New(t)

	args := map[string]string{
		"prefix": "[INFO]",
		"flag":   "log.lstdflags",
		"misc":   "misc",
	}
	w, err := logWriterInitializer(args)
	a.NotError(err).NotNil(w)

	lw, ok := w.(*logWriter)
	a.True(ok).NotNil(lw)

	a.Equal(lw.prefix, "[INFO]").
		Equal(lw.flag, log.LstdFlags)
}
