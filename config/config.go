// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// logs包的配置文件处理。
package config

import (
	"errors"
	"fmt"
)

// 用于表示xml配置文件中的数据。
type Config struct {
	Parent *Config
	Name   string             // writer的名称，一般为节点名
	Attrs  map[string]string  // 参数列表
	Items  map[string]*Config // 若是容器，则还有子项
}

// 检测语法错误及基本的内容错误。
func check(cfg *Config) error {
	if cfg.Name != "logs" {
		return fmt.Errorf("check:顶级元素必须为logs，当前名称为[%v]", cfg.Name)
	}

	if len(cfg.Attrs) > 0 {
		return fmt.Errorf("check:logs元素不存在任何属性")
	}

	if len(cfg.Items) == 0 {
		return errors.New("check:空的logs元素")
	}

	if len(cfg.Items) > 6 {
		return errors.New("check:logs最多只有6个子元素")
	}

	for name, item := range cfg.Items {
		if len(item.Items) == 0 {
			return fmt.Errorf("check:[%v]并未指定子元素", name)
		}

		switch name {
		case "info":
		case "warn":
		case "debug":
		case "error":
		case "trace":
		case "critical":
		default:
			return fmt.Errorf("check:未知道的二级元素名称:[%v]", name)
		}
	}

	return nil
}
