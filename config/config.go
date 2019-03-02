// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package config 表示 logs 包的配置文件处理。
package config

import "github.com/issue9/config"

// Config 用于表示配置文件中的数据。
type Config struct {
	parent *Config // TODO 仅 xml 使用，考虑去掉

	Name  string             `yaml:"name"`  // writer 的名称，一般为节点名
	Attrs map[string]string  `yaml:"attrs"` // 参数列表
	Items map[string]*Config `yaml:"items"` // 若是容器，则还有子项
}

// Sanitize 检测语法错误及基本的内容错误。
//
// 同时也是实现 config.Sanitizer 接口。
func (cfg *Config) Sanitize() error {
	if cfg.Name != "logs" {
		return config.NewError("", "name", "顶级元素必须为 logs")
	}

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
