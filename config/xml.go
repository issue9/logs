// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// ParseXMLFile 从一个 XML 文件初始化 Config 实例。
func ParseXMLFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	if err := XMLUnmarshal([]byte(bs), conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// ParseXMLString 从一个 XML 字符串初始化 Config 实例。
func ParseXMLString(xml string) (*Config, error) {
	conf := &Config{}
	if err := XMLUnmarshal([]byte(xml), conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// XMLUnmarshal XML 的 Unmarhsal 接口
func XMLUnmarshal(bs []byte, v interface{}) error {
	ret, ok := v.(*Config)
	if !ok {
		panic("参数 v 的类型不是 *Config")
	}

	var cfg *Config
	var t xml.Token
	var err error
	d := xml.NewDecoder(bytes.NewBuffer(bs))
	for t, err = d.Token(); err == nil; t, err = d.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			c := &Config{
				parent: cfg,
				Name:   token.Name.Local,
				Attrs:  make(map[string]string, len(token.Attr)),
			}
			for _, v := range token.Attr {
				c.Attrs[v.Name.Local] = v.Value
			}

			if cfg != nil {
				if cfg.Items == nil {
					cfg.Items = make(map[string]*Config)
				}

				if _, found := cfg.Items[token.Name.Local]; found {
					return fmt.Errorf("重复的元素名[%v]", token.Name.Local)
				}

				cfg.Items[token.Name.Local] = c
			}
			cfg = c
		case xml.EndElement:
			if cfg.parent != nil {
				cfg = cfg.parent
			}
		default: // 可能还有 ProcInst、CharData、Comment 等用不到的标签
			continue
		}
	} // end for

	if err != io.EOF {
		return err
	}

	if err = cfg.sanitize(); err != nil {
		return err
	}

	*ret = *cfg
	return nil
}
