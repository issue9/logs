// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"fmt"
	"io"
	"sync"

	"github.com/issue9/logs/internal/config"
	"github.com/issue9/logs/writers"
)

// 将当前的 config.Config 转换成 io.Writer
func toWriter(c *config.Config) (io.Writer, error) {
	fun, found := inits.funs[c.Name]
	if !found {
		return nil, fmt.Errorf("toWriter:未注册的初始化函数:[%v]", c.Name)
	}

	w, err := fun(c.Attrs)
	if err != nil {
		return nil, err
	}

	if len(c.Items) == 0 { // 没有子项
		return w, err
	}

	cont, ok := w.(writers.Adder)
	if !ok {
		return nil, fmt.Errorf("toWriter:[%v]并未实现writers.Adder接口", c.Name)
	}

	for _, cfg := range c.Items {
		wr, err := toWriter(cfg)
		if err != nil {
			return nil, err
		}
		cont.Add(wr)
	}

	return w, nil
}

// writer 的初始化函数。
// args 参数为对应的 XML 节点的属性列表。
type WriterInitializer func(args map[string]string) (io.Writer, error)

type initMap struct {
	sync.Mutex
	funs map[string]WriterInitializer
}

var inits = &initMap{
	funs: map[string]WriterInitializer{},
}

// 注册一个 writer 初始化函数。
// writer 初始化函数原型可参考: WriterInitializer。
// 返回值反映是否注册成功。若已经存在相同名称的，则返回false
func Register(name string, init WriterInitializer) bool {
	inits.Lock()
	defer inits.Unlock()

	if _, found := inits.funs[name]; found {
		return false
	}

	inits.funs[name] = init
	return true
}

// 查询指定名称的 Writer 是否已经被注册
func IsRegisted(name string) bool {
	inits.Lock()
	defer inits.Unlock()

	_, found := inits.funs[name]
	return found
}

// 返回所有已注册的 writer 名称
func Registed() []string {
	inits.Lock()
	defer inits.Unlock()

	names := make([]string, 0, len(inits.funs))
	for name := range inits.funs {
		names = append(names, name)
	}

	return names
}
