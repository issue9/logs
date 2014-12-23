// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"encoding/xml"
	"fmt"
	"io"
	"sync"

	"github.com/issue9/logs/writer"
)

// 用于表示xml配置文件中的数据。
type config struct {
	parent *config
	name   string             // writer的名称，一般为节点名
	attrs  map[string]string  // 参数列表
	items  map[string]*config // 若是容器，则还有子项
}

// 从一个xml格式的reader初始化config
func loadFromXml(r io.Reader) (*config, error) {
	var cfg *config = nil
	var t xml.Token
	var err error

	d := xml.NewDecoder(r)
	for t, err = d.Token(); err == nil; t, err = d.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			c := &config{
				parent: cfg,
				name:   token.Name.Local,
				attrs:  make(map[string]string),
			}
			for _, v := range token.Attr {
				c.attrs[v.Name.Local] = v.Value
			}

			if cfg != nil {
				if cfg.items == nil {
					cfg.items = make(map[string]*config)
				}
				cfg.items[token.Name.Local] = c
			}
			cfg = c
		case xml.EndElement:
			if cfg.parent != nil {
				cfg = cfg.parent
			}
		default: // 可能还有ProcInst,CharData,Comment等用不到的标签
			continue
		}
	} // end for

	if err != io.EOF {
		return nil, err
	}

	return cfg, nil
}

// 将当前的config转换成io.Writer
func (c *config) toWriter() (io.Writer, error) {
	fun, found := inits.funs[c.name]
	if !found {
		return nil, fmt.Errorf("未注册的初始化函数:[%v]", c.name)
	}

	w, err := fun(c.attrs)
	if err != nil {
		return nil, err
	}

	if len(c.items) == 0 { // 没有子项
		return w, err
	}

	cont, ok := w.(writer.Adder)
	if !ok {
		return nil, fmt.Errorf("[%v]并未实现writer.Adder接口", c.name)
	}

	for _, cfg := range c.items {
		wr, err := cfg.toWriter()
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

// 清除已经注册的初始化函数。
func clearInitializer() {
	inits.Lock()
	defer inits.Unlock()

	inits.funs = make(map[string]WriterInitializer)
	inits.names = inits.names[:0]
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
