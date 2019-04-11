// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package config 表示 logs 包的配置文件处理。
package config

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/issue9/config"
)

// Config 用于表示配置文件中的数据。
//
// 提供了对 JSON、XML 和 YAML 的支持
type Config struct {
	parent *Config

	Attrs map[string]string  `yaml:"attrs" json:"attrs"` // 参数列表
	Items map[string]*Config `yaml:"items" json:"items"` // 若是容器，则还有子项
}

// Sanitize 检测语法错误及基本的内容错误。
//
// 同时也是实现 config.Sanitizer 接口。
func (cfg *Config) Sanitize() error {
	if len(cfg.Attrs) > 0 {
		return config.NewError("", "attrs", "根元素不能存在任何属性")
	}

	if len(cfg.Items) == 0 {
		return config.NewError("", "items", "不能为空")
	}

	for name, item := range cfg.Items {
		if len(item.Items) == 0 {
			return config.NewError("", name+".items", "不能为空")
		}
	}

	return nil
}

// UnmarshalXML xml.Unmarshaler 接口实现
func (cfg *Config) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for t, err := d.Token(); ; t, err = d.Token() {
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}

		switch token := t.(type) {
		case xml.StartElement:
			c := &Config{
				parent: cfg,
				Attrs:  make(map[string]string, len(token.Attr)),
			}
			for _, v := range token.Attr {
				c.Attrs[v.Name.Local] = v.Value
			}

			if cfg.Items == nil {
				cfg.Items = make(map[string]*Config)
			}
			if _, found := cfg.Items[token.Name.Local]; found {
				return fmt.Errorf("重复的元素名[%v]", token.Name.Local)
			}
			cfg.Items[token.Name.Local] = c

			cfg = c
		case xml.EndElement:
			if cfg.parent != nil {
				cfg = cfg.parent
			}
		} // end switch
	} // end for
}

// MarshalXML xml.Unmarshaler 接口实现
func (cfg *Config) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return cfg.marshalXML(e, xml.StartElement{Name: xml.Name{Local: "logs"}})
}

func (cfg *Config) marshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for k, v := range cfg.Items {
		s := xml.StartElement{
			Name: xml.Name{Local: k},
			Attr: make([]xml.Attr, 0, len(v.Attrs)),
		}
		for name, val := range v.Attrs {
			s.Attr = append(s.Attr, xml.Attr{
				Name:  xml.Name{Local: name},
				Value: val,
			})
		}
		if err := v.marshalXML(e, s); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
