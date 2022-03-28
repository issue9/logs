// SPDX-License-Identifier: MIT

package logs

type Level int8

// 目前支持的日志类型
const (
	LevelInfo Level = iota + 1
	LevelTrace
	LevelDebug
	LevelWarn
	LevelError
	LevelFatal
)

var levelStrings = map[Level]string{
	LevelInfo:  "INFO",
	LevelTrace: "TRAC",
	LevelDebug: "DBUG",
	LevelWarn:  "WARN",
	LevelError: "ERRO",
	LevelFatal: "FATL",
}

func (l Level) String() string { return levelStrings[l] }

func (l Level) MarshalText() ([]byte, error) { return []byte(l.String()), nil }
