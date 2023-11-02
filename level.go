// SPDX-License-Identifier: MIT

package logs

import (
	"fmt"
	"strings"
)

type Level int8

// 目前支持的日志类型
const (
	levelStart Level = iota
	LevelInfo
	LevelTrace
	LevelDebug
	LevelWarn
	LevelError
	LevelFatal
	levelEnd
)

var levelStrings = map[Level]string{
	LevelInfo:  "INFO",
	LevelTrace: "TRAC",
	LevelDebug: "DBUG",
	LevelWarn:  "WARN",
	LevelError: "ERRO",
	LevelFatal: "FATL",
}

func IsValidLevel(l Level) bool { return l > levelStart && l < levelEnd }

func (l Level) String() string { return levelStrings[l] }

func (l Level) MarshalText() ([]byte, error) { return []byte(l.String()), nil }

func (l *Level) UnmarshalText(data []byte) error {
	lv, err := ParseLevel(string(data))
	if err != nil {
		return err
	}
	*l = lv
	return nil
}

func ParseLevel(s string) (Level, error) {
	s = strings.ToUpper(s)

	for level, name := range levelStrings {
		if s == name {
			return level, nil
		}
	}

	return -1, fmt.Errorf("无效的值 %s", s)
}
