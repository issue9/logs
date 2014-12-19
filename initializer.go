// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/issue9/conv"
	"github.com/issue9/logs/writer"
	"github.com/issue9/term"
)

// writer的初始化函数。
// args参数为对应的xml节点的属性列表。
type WriterInitializer func(args map[string]string) (io.Writer, error)

// 注册的writer，所有注册的writer，都可以通过配置文件配置。
var regInitializer = map[string]WriterInitializer{}
var regNames []string

// 清除已经注册的初始化函数。
func clearInitializer() {
	regInitializer = make(map[string]WriterInitializer)
	regNames = regNames[:0]
}

// 注册一个initizlizer
// 返回值反映是否注册成功。若已经存在相同名称的，则返回false
func Register(name string, init WriterInitializer) bool {
	if IsRegisted(name) {
		return false
	}

	regInitializer[name] = init
	regNames = append(regNames, name)
	return true
}

// 查询指定名称的Writer是否已经被注册
func IsRegisted(name string) bool {
	_, found := regInitializer[name]
	return found
}

// 返回所有已注册的writer名称
func Registed() []string {
	return regNames
}

func argNotFoundErr(wname, argName string) error {
	return fmt.Errorf("[%v]配置文件中未指定参数:[%v]", wname, argName)
}

// writer.Rotate的初始化函数。
func rotateInitializer(args map[string]string) (io.Writer, error) {
	dir, found := args["dir"]
	if !found {
		return nil, argNotFoundErr("rotate", "dir")
	}

	sizeStr, found := args["size"]
	if !found {
		return nil, argNotFoundErr("rotate", "size")
	}

	size, err := conv.ToByte(sizeStr)
	if err != nil {
		return nil, err
	}

	return writer.NewRotate(dir, int(size))
}

// writer.Buffer的初始化函数
func bufferInitializer(args map[string]string) (io.Writer, error) {
	size, found := args["size"]
	if !found {
		return nil, argNotFoundErr("buffer", "size")
	}

	num, err := strconv.Atoi(size)
	if err != nil {
		return nil, err
	}

	return writer.NewBuffer(nil, num), nil
}

var consoleOutputMap = map[string]io.Writer{
	"os.stderr": os.Stderr,
	"os.stdin":  os.Stdin,
	"os.stdout": os.Stdout,
}

// writer.Console的初始化函数
func consoleInitializer(args map[string]string) (io.Writer, error) {
	outputIndex, found := args["output"]
	if !found {
		outputIndex = "os.stderr"
	}

	output, found := consoleOutputMap[outputIndex]
	if !found {
		return nil, fmt.Errorf("[%v]不是一个有效的控制台输出项", outputIndex)
	}

	color, found := args["color"]
	if !found {
		color = term.FRed
	}

	if color[0] != '\033' && color[len(color)-1] != 'm' {
		return nil, fmt.Errorf("color的值[%v]必须为一个ansi color值", color)
	}

	return writer.NewConsole(output, color), nil
}

// writer.Stmp的初始化函数
func stmpInitializer(args map[string]string) (io.Writer, error) {
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

	return writer.NewSmtp(username, password, subject, host, sendTo), nil
}

func init() {
	if !Register("stmp", stmpInitializer) {
		panic("注册stmp时失败")
	}

	if !Register("console", consoleInitializer) {
		panic("注册console时失败")
	}

	if !Register("buffer", bufferInitializer) {
		panic("注册buffer时失败")
	}

	if !Register("rotate", rotateInitializer) {
		panic("注册rotate时失败")
	}
}
