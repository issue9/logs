// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"fmt"
	"io"

	"github.com/issue9/logs/v2/config"
	"github.com/issue9/logs/v2/internal/initfunc"
	"github.com/issue9/logs/v2/writers"
)

// WriterInitializer io.Writer 实例的初始化函数。
// args 参数为对应的 XML 节点的属性列表。
type WriterInitializer func(cfg *config.Config) (io.Writer, error)

var funs = map[string]WriterInitializer{}

// 将当前的 config.Config 转换成 io.Writer
func toWriter(name string, c *config.Config) (io.Writer, error) {
	fun, found := funs[name]
	if !found {
		return nil, fmt.Errorf("未注册的初始化函数:[%v]", name)
	}

	w, err := fun(c)
	if err != nil {
		return nil, err
	}

	if len(c.Items) == 0 { // 没有子项
		return w, err
	}

	cont, ok := w.(writers.Adder)
	if !ok {
		return nil, fmt.Errorf("[%v]并未实现 writers.Adder 接口", name)
	}

	for name, cfg := range c.Items {
		wr, err := toWriter(name, cfg)
		if err != nil {
			return nil, err
		}
		cont.Add(wr)
	}

	return w, nil
}

// Register 注册一个 WriterInitializer。
//
// writer 初始化函数原型可参考: WriterInitializer。
// 返回值反映是否注册成功。若已经存在相同名称的，则返回 false
func Register(name string, init WriterInitializer) bool {
	if _, found := funs[name]; found {
		return false
	}

	funs[name] = init
	return true
}

// IsRegisted 查询指定名称的 Writer 是否已经被注册
func IsRegisted(name string) bool {
	_, found := funs[name]
	return found
}

// Registed 返回所有已注册的 writer 名称
func Registed() []string {
	names := make([]string, 0, len(funs))
	for name := range funs {
		names = append(names, name)
	}

	return names
}

// 注册各类初始化函数
func init() {
	if !Register("smtp", initfunc.SMTP) {
		panic("注册 smtp 时失败，已存在相同名称")
	}

	if !Register("console", initfunc.Console) {
		panic("注册 console 时失败，已存在相同名称")
	}

	if !Register("buffer", initfunc.Buffer) {
		panic("注册 buffer 时失败，已存在相同名称")
	}

	if !Register("rotate", initfunc.Rotate) {
		panic("注册 rotate 时失败，已存在相同名称")
	}

	// logContInitializer
	for key := range levels {
		if !Register(key, loggerInitializer) {
			panic(fmt.Sprintf("注册 %v 时失败，已存在相同名称", key))
		}
	}
}
