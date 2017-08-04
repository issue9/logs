// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package config 表示 logs 包的配置文件处理。
package config

import (
	"errors"
	"fmt"
)

// Config 用于表示配置文件中的数据。
type Config struct {
	Parent *Config
	Name   string             // writer 的名称，一般为节点名
	Attrs  map[string]string  // 参数列表
	Items  map[string]*Config // 若是容器，则还有子项
}

// 检测语法错误及基本的内容错误。
func check(cfg *Config) error {
	if cfg.Name != "logs" {
		return fmt.Errorf("顶级元素必须为 logs，当前名称为 %s", cfg.Name)
	}

	if len(cfg.Attrs) > 0 {
		return fmt.Errorf("根元素不能存在任何属性")
	}

	if len(cfg.Items) == 0 {
		return errors.New("空的 logs 元素")
	}

	for name, item := range cfg.Items {
		if len(item.Items) == 0 {
			return fmt.Errorf("%s 并未指定子元素", name)
		}
	}

	return nil
}
