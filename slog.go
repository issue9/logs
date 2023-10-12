// SPDX-License-Identifier: MIT

//go:build go1.21

package logs

import (
	"context"
	"log/slog"
	"runtime"
	"slices"
	"strconv"
)

var slog2Logs = map[slog.Level]Level{
	slog.LevelInfo:  LevelInfo,
	slog.LevelDebug: LevelDebug,
	slog.LevelWarn:  LevelWarn,
	slog.LevelError: LevelError,
}

type logsHandler struct {
	l      *Logs
	attrs  []slog.Attr
	prefix string // groups 组成
}

// SLogHandler 将 logs 转换为 [slog.Handler] 接口
//
// 所有的 group 会作为普通 attr 的名称前缀，但是不影响 Level、Message 等字段。
func (l *Logs) SLogHandler() slog.Handler { return &logsHandler{l: l} }

// SLog 将 Logs 作为 [slog.Logger] 的后端
func (l *Logs) SLog() *slog.Logger { return slog.New(l.SLogHandler()) }

func (h *logsHandler) Enabled(ctx context.Context, lv slog.Level) bool {
	return h.l.IsEnable(slog2Logs[lv])
}

func (h *logsHandler) Handle(ctx context.Context, r slog.Record) error {
	rr := h.l.NewRecord(slog2Logs[r.Level])
	rr.Created = r.Time.Format(h.l.createdFormat)
	rr.Message = r.Message

	for _, attr := range h.attrs {
		rr.With(h.prefix+attr.Key, attr.Value)
	}
	r.Attrs(func(attr slog.Attr) bool {
		rr.With(h.prefix+attr.Key, attr.Value)
		return true
	})

	if r.PC != 0 {
		f, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
		rr.Path = f.File+":"+strconv.Itoa(f.Line)
	}

	rr.output()

	return nil
}

func (h *logsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var a []slog.Attr
	if len(h.attrs) == 0 {
		a = attrs
	} else {
		a = append(slices.Clip(h.attrs), attrs...)
	}

	return &logsHandler{
		l:      h.l,
		attrs:  a,
		prefix: h.prefix,
	}
}

func (h *logsHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	return &logsHandler{
		l:      h.l,
		attrs:  h.attrs,
		prefix: name + "." + h.prefix,
	}
}
