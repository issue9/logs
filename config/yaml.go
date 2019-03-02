// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import yaml "gopkg.in/yaml.v2"

// YAMLUnmarshal YAML 的 Unmarhsal 接口
func YAMLUnmarshal(bs []byte, v interface{}) error {
	if err := yaml.Unmarshal(bs, v); err != nil {
		return err
	}

	conf, ok := v.(*Config)
	if !ok {
		panic("参数 v 不是 *Config 类型")
	}

	setName(conf)
	return nil
}

func setName(conf *Config) {
	for k, v := range conf.Items {
		if v.Name == "" {
			v.Name = k
		} else if v.Name != k {
			panic("名称不相符")
		}

		setName(v)
	}
}
