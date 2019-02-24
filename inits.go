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

// 将当前的 config.Config 转换成 io.Writer
func (logs *Logs) toWriter(c *config.Config) (io.Writer, error) {
	fun, found := logs.funs[c.Name]
	if !found {
		return nil, fmt.Errorf("未注册的初始化函数:[%v]", c.Name)
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
		return nil, fmt.Errorf("[%v]并未实现 writers.Adder 接口", c.Name)
	}

	for _, cfg := range c.Items {
		wr, err := logs.toWriter(cfg)
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
func (logs *Logs) Register(name string, init WriterInitializer) bool {
	if _, found := logs.funs[name]; found {
		return false
	}

	logs.funs[name] = init
	return true
}

// Register 注册一个 WriterInitializer。
func Register(name string, init WriterInitializer) bool {
	return defaultLogs.Register(name, init)
}

// IsRegisted 查询指定名称的 Writer 是否已经被注册
func (logs *Logs) IsRegisted(name string) bool {
	_, found := logs.funs[name]
	return found
}

// IsRegisted 查询指定名称的 Writer 是否已经被注册
func IsRegisted(name string) bool {
	return defaultLogs.IsRegisted(name)
}

// Registed 返回所有已注册的 writer 名称
func (logs *Logs) Registed() []string {
	names := make([]string, 0, len(logs.funs))
	for name := range logs.funs {
		names = append(names, name)
	}

	return names
}

// Registed 返回所有已注册的 writer 名称
func Registed() []string {
	return defaultLogs.Registed()
}

// 注册各类初始化函数
func (logs *Logs) initFuncs() {
	if !logs.Register("smtp", initfunc.SMTP) {
		panic("注册 smtp 时失败，已存在相同名称")
	}

	if !logs.Register("console", initfunc.Console) {
		panic("注册 console 时失败，已存在相同名称")
	}

	if !logs.Register("buffer", initfunc.Buffer) {
		panic("注册 buffer 时失败，已存在相同名称")
	}

	if !logs.Register("rotate", initfunc.Rotate) {
		panic("注册 rotate 时失败，已存在相同名称")
	}

	// logContInitializer
	for key := range levels {
		if !logs.Register(key, loggerInitializer) {
			panic(fmt.Sprintf("注册 %v 时失败，已存在相同名称", key))
		}
	}
}
