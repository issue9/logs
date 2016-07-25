// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/issue9/logs/writers"
	"github.com/issue9/term/colors"
)

// 本文件下声明一系列writer的注册函数。

const (
	b int64 = 1 << (10 * iota)
	kb
	mb
	gb
)

// 将字符串转换成以字节为单位的数值。
// 粗略计算，并不100%正确，小数只取整数部分。
// 支持以下格式：
//  1024
//  1k
//  1M
//  1G
// 后缀单位只支持k,g,m，不区分大小写。
func toByte(str string) (int64, error) {
	if len(str) == 0 {
		return -1, errors.New("不能传递空值")
	}

	str = strings.ToLower(str)

	scale := b
	unit := str[len(str)-1]
	switch {
	case unit >= '0' && unit <= '9':
		scale = b
	case unit == 'b':
		scale = b
	case unit == 'k':
		scale = kb
	case unit == 'm':
		scale = mb
	case unit == 'g':
		scale = gb
	default:
		return -1, fmt.Errorf("无法识别的单位:[%v]", unit)
	}

	if scale > 1 {
		str = str[:len(str)-1]
	}

	if len(str) == 0 {
		return -1, errors.New("传递了一个空值")
	}

	size, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return -1, err
	}

	if size <= 0 {
		return -1, fmt.Errorf("大小不能小于0，当前值为:[%v]", size)
	}

	return int64(size) * scale, nil
}

func argNotFoundErr(wname, argName string) error {
	return fmt.Errorf("[%v]配置文件中未指定参数:[%v]", wname, argName)
}

// writers.Rotate 的初始化函数。
func rotateInitializer(args map[string]string) (io.Writer, error) {
	prefix, found := args["prefix"]
	if !found {
		prefix = ""
	}

	dir, found := args["dir"]
	if !found {
		return nil, argNotFoundErr("rotate", "dir")
	}

	sizeStr, found := args["size"]
	if !found {
		return nil, argNotFoundErr("rotate", "size")
	}

	size, err := toByte(sizeStr)
	if err != nil {
		return nil, err
	}

	return writers.NewRotate(prefix, dir, int(size))
}

// writers.Buffer 的初始化函数
func bufferInitializer(args map[string]string) (io.Writer, error) {
	size, found := args["size"]
	if !found {
		return nil, argNotFoundErr("buffer", "size")
	}

	num, err := strconv.Atoi(size)
	if err != nil {
		return nil, err
	}

	return writers.NewBuffer(num), nil
}

var consoleOutputMap = map[string]*os.File{
	"stderr": os.Stderr,
	"stdout": os.Stdout,
	"stdin":  os.Stdin,
}

var consoleColorMap = map[string]colors.Color{
	"default": colors.Default,
	"black":   colors.Black,
	"red":     colors.Red,
	"green":   colors.Green,
	"yellow":  colors.Yellow,
	"blue":    colors.Blue,
	"magenta": colors.Magenta,
	"cyan":    colors.Cyan,
	"white":   colors.White,
}

// writers.Console 的初始化函数
func consoleInitializer(args map[string]string) (io.Writer, error) {
	outputIndex, found := args["output"]
	if !found {
		outputIndex = "stderr"
	}

	output, found := consoleOutputMap[outputIndex]
	if !found {
		return nil, fmt.Errorf("[%v]不是一个有效的控制台输出项", outputIndex)
	}

	fcIndex, found := args["foreground"]
	if !found { // 默认用红色前景色
		fcIndex = "red"
	}

	fc, found := consoleColorMap[fcIndex]
	if !found {
		return nil, fmt.Errorf("无效的前景色[%v]", fcIndex)
	}

	bcIndex, found := args["background"]
	if !found {
		bcIndex = "default"
	}

	bc, found := consoleColorMap[bcIndex]
	if !found {
		return nil, fmt.Errorf("无效的背景色[%v]", bcIndex)
	}

	return writers.NewConsole(output, fc, bc), nil
}

// writers.Stmp 的初始化函数
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

	return writers.NewSmtp(username, password, subject, host, sendTo), nil
}

var flagMap = map[string]int{
	"none":              0,
	"log.ldate":         log.Ldate,
	"log.ltime":         log.Ltime,
	"log.lmicroseconds": log.Lmicroseconds,
	"log.llongfile":     log.Llongfile,
	"log.lshortfile":    log.Lshortfile,
	"log.lstdflags":     log.LstdFlags,
}

func logContInitializer(args map[string]string) (io.Writer, error) {
	return writers.NewContainer(), nil
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

	// logWriter

	if !Register("info", logContInitializer) {
		panic("注册info时失败")
	}

	if !Register("debug", logContInitializer) {
		panic("注册debug时失败")
	}

	if !Register("trace", logContInitializer) {
		panic("注册trace时失败")
	}

	if !Register("warn", logContInitializer) {
		panic("注册warn时失败")
	}

	if !Register("error", logContInitializer) {
		panic("注册error时失败")
	}

	if !Register("critical", logContInitializer) {
		panic("注册critical时失败")
	}
}
