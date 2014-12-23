// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/issue9/conv"
	"github.com/issue9/logs/writer"
	"github.com/issue9/term"
)

// 本文件下声明一系列writer的注册函数。

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

var flagMap = map[string]int{
	"log.ldate":         log.Ldate,
	"log.ltime":         log.Ltime,
	"log.lmicroseconds": log.Lmicroseconds,
	"log.llongfile":     log.Llongfile,
	"log.lshortfile":    log.Lshortfile,
	"log.lstdflags":     log.LstdFlags,
}

func logWriterInitializer(args map[string]string) (io.Writer, error) {
	flagStr, found := args["flag"]
	if !found || (flagStr == "") {
		flagStr = "log.lstdflags"
	}

	flag, found := flagMap[strings.ToLower(flagStr)]
	if !found {
		return nil, fmt.Errorf("未知的Flag参数:[%v]", flagStr)
	}

	return newLogWriter(args["prefix"], flag), nil
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

	if !Register("info", logWriterInitializer) {
		panic("注册info时失败")
	}

	if !Register("debug", logWriterInitializer) {
		panic("注册debug时失败")
	}

	if !Register("trace", logWriterInitializer) {
		panic("注册trace时失败")
	}

	if !Register("warn", logWriterInitializer) {
		panic("注册warn时失败")
	}

	if !Register("error", logWriterInitializer) {
		panic("注册error时失败")
	}

	if !Register("critical", logWriterInitializer) {
		panic("注册critical时失败")
	}
}
