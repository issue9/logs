// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// logs包的配置文件处理。
// 只对语法错误负责，不负责配置内容的对错。
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

	return nil
}
