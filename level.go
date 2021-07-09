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

var levels = map[string]int{
	"info":     LevelInfo,
	"trace":    LevelTrace,
	"debug":    LevelDebug,
	"warn":     LevelWarn,
	"error":    LevelError,
	"critical": LevelCritical,
}

func (l *Logs) logs(level int) []*logger {
	logs := make([]*logger, 0, len(l.loggers))

	for key, item := range l.loggers {
		if key&level == key {
			logs = append(logs, item)
		}
	}

	return logs
}
