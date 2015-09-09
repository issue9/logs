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

// 将当前的config.Config转换成io.Writer
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

// writer的初始化函数。
// args参数为对应的xml节点的属性列表。
type WriterInitializer func(args map[string]string) (io.Writer, error)

type initMap struct {
	sync.Mutex
	funs  map[string]WriterInitializer
	names []string
}

var inits = &initMap{
	funs:  map[string]WriterInitializer{},
	names: []string{},
}

// 注册一个writer初始化函数。
// writer初始化函数原型可参考:WriterInitializer。
// 返回值反映是否注册成功。若已经存在相同名称的，则返回false
func Register(name string, init WriterInitializer) bool {
	inits.Lock()
	defer inits.Unlock()

	if _, found := inits.funs[name]; found {
		return false
	}

	inits.funs[name] = init
	inits.names = append(inits.names, name)
	return true
}

// 查询指定名称的Writer是否已经被注册
func IsRegisted(name string) bool {
	inits.Lock()
	defer inits.Unlock()

	_, found := inits.funs[name]
	return found
}

// 返回所有已注册的writer名称
func Registed() []string {
	return inits.names
}
