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

// ParseXMLFile 从一个 XML 文件初始化 Config 实例
func ParseXMLFile(path string) (*Config, error) {
	return ParseFile(path, xml.Unmarshal)
}

// ParseXMLFileFS 从一个 XML 文件初始化 Config 实例
func ParseXMLFileFS(f fs.FS, path string) (*Config, error) {
	return ParseFS(f, path, xml.Unmarshal)
}

// ParseXMLString 从一个 XML 字符串初始化 Config 实例
func ParseXMLString(data string) (*Config, error) {
	return ParseString(data, xml.Unmarshal)
}

// ParseJSONFile 从一个 JSON 文件初始化 Config 实例
func ParseJSONFile(path string) (*Config, error) {
	return ParseFile(path, json.Unmarshal)
}

// ParseJSONFileFS 从一个 JSON 文件初始化 Config 实例
func ParseJSONFileFS(f fs.FS, path string) (*Config, error) {
	return ParseFS(f, path, json.Unmarshal)
}

// ParseJSONString 从一个 JSON 字符串初始化 Config 实例
func ParseJSONString(data string) (*Config, error) {
	return ParseString(data, json.Unmarshal)
}

// ParseYAMLFile 从一个 YAML 文件初始化 Config 实例
func ParseYAMLFile(path string) (*Config, error) {
	return ParseFile(path, yaml.Unmarshal)
}

// ParseYAMLFileFS 从一个 YAML 文件初始化 Config 实例
func ParseYAMLFileFS(f fs.FS, path string) (*Config, error) {
	return ParseFS(f, path, yaml.Unmarshal)
}

// ParseYAMLString 从一个 YAML 字符串初始化 Config 实例
func ParseYAMLString(data string) (*Config, error) {
	return ParseString(data, yaml.Unmarshal)
}

// ParseFile 从文件中初始化 Config 对象
func ParseFile(path string, u func([]byte, interface{}) error) (*Config, error) {
	dir, base := filepath.Split(path)
	return ParseFS(os.DirFS(dir), filepath.ToSlash(base), u)
}

// ParseFS 从文件中初始化 Config 对象
func ParseFS(f fs.FS, path string, u func([]byte, interface{}) error) (*Config, error) {
	bs, err := fs.ReadFile(f, path)
	if err != nil {
		return nil, err
	}
	return Parse(bs, u)
}

// ParseFS 从字符串初始化 Config 对象
func ParseString(data string, u func([]byte, interface{}) error) (*Config, error) {
	return Parse([]byte(data), u)
}

// Parse 从 []byte 初始化 Config 对象
func Parse(data []byte, u func([]byte, interface{}) error) (*Config, error) {
	conf := &Config{}
	if err := u(data, conf); err != nil {
		return nil, err
	}

	if err := conf.Sanitize(); err != nil {
		return nil, err
	}
	return conf, nil
}
