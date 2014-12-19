// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package logs

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/issue9/logs/writer"
)

// 用于表示config.xml中的配置数据。
type config struct {
	parent *config
	name   string             // writer的名称，一般为节点名
	attrs  map[string]string  // 参数列表
	items  map[string]*config // 若是容器，则还有子项
}

// 从一个xml reader初始化config
func loadFromXml(r io.Reader) (*config, error) {
	var cfg *config = nil //&config{parent: nil}
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
	initializer, found := regInitializer[c.name]
	if !found {
		return nil, fmt.Errorf("未注册的初始化函数:[%v]", c.name)
	}

	w, err := initializer(c.attrs)
	if err != nil {
		return nil, err
	}

	if len(c.items) == 0 {
		return w, err
	}

	cont, ok := w.(writer.FlushAdder)
	if !ok {
		return nil, fmt.Errorf("[%v]并未实现writer.FlushAdder接口", c.name)
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
