// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"encoding/xml"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// ParseXMLFile 从一个 XML 文件初始化 Config 实例。
func ParseXMLFile(path string) (*Config, error) {
	return ParseFile(path, xml.Unmarshal)
}

// ParseXMLString 从一个 XML 字符串初始化 Config 实例。
func ParseXMLString(data string) (*Config, error) {
	return ParseString(data, xml.Unmarshal)
}

// ParseJSONFile 从一个 JSON 文件初始化 Config 实例。
func ParseJSONFile(path string) (*Config, error) {
	return ParseFile(path, json.Unmarshal)
}

// ParseJSONString 从一个 JSON 字符串初始化 Config 实例。
func ParseJSONString(data string) (*Config, error) {
	return ParseString(data, json.Unmarshal)
}

// ParseYAMLFile 从一个 YAML 文件初始化 Config 实例。
func ParseYAMLFile(path string) (*Config, error) {
	return ParseFile(path, yaml.Unmarshal)
}

// ParseYAMLString 从一个 YAML 字符串初始化 Config 实例。
func ParseYAMLString(data string) (*Config, error) {
	return ParseString(data, yaml.Unmarshal)
}

// ParseFile 从文件中初始化 Config 对象，由 unmarshal 决定解析方式
func ParseFile(path string, unmarshal func([]byte, interface{}) error) (*Config, error) {
	dir, base := filepath.Split(path)
	return ParseFS(os.DirFS(dir), filepath.ToSlash(base), unmarshal)
}

// ParseFS 从文件中初始化 Config 对象，由 unmarshal 决定解析方式
func ParseFS(f fs.FS, path string, unmarshal func([]byte, interface{}) error) (*Config, error) {
	bs, err := fs.ReadFile(f, path)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	if err := unmarshal(bs, conf); err != nil {
		return nil, err
	}

	if err := conf.Sanitize(); err != nil {
		return nil, err
	}
	return conf, nil
}

// ParseString 从字符串中初始化 Config 对象，由 unmarshal 决定解析方式
func ParseString(data string, unmarshal func([]byte, interface{}) error) (*Config, error) {
	conf := &Config{}
	if err := unmarshal([]byte(data), conf); err != nil {
		return nil, err
	}

	if err := conf.Sanitize(); err != nil {
		return nil, err
	}
	return conf, nil
}
