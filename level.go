// SPDX-License-Identifier: MIT

package logs

// 目前支持的日志类型
const (
	LevelInfo = 1 << iota
	LevelTrace
	LevelDebug
	LevelWarn
	LevelError
	LevelCritical
	LevelAll = LevelInfo + LevelTrace + LevelDebug + LevelWarn + LevelError + LevelCritical
)

func (l *Logs) walk(level int, walk func(l *logger) error) error {
	for key, item := range l.loggers {
		if key&level == key {
			if err := walk(item); err != nil {
				return err
			}
		}
	}
	return nil
}
