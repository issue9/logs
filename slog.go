// SPDX-FileCopyrightText: 2014-2024 caixw
//
// SPDX-License-Identifier: MIT

package logs

import (
	"context"
	"log/slog"
	"runtime"
	"slices"
)

var slog2Logs = map[slog.Level]Level{
	slog.LevelInfo:  LevelInfo,
	slog.LevelDebug: LevelDebug,
	slog.LevelWarn:  LevelWarn,
	slog.LevelError: LevelError,
}

type slogHandler struct {
	l      *Logs
	attrs  []slog.Attr
	prefix string // groups 组成
}

// SLogHandler 将 logs 转换为 [slog.Handler] 接口
//
// 所有的 group 会作为普通 attr 的名称前缀，但是不影响 Level、Message 等字段。
func (l *Logs) SLogHandler() slog.Handler { return &slogHandler{l: l} }

func (h *slogHandler) Enabled(ctx context.Context, lv slog.Level) bool {
	return h.l.IsEnable(slog2Logs[lv])
}

func (h *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	rr := h.l.NewRecord()
	rr.AppendCreated = func(b *Buffer) { b.AppendTime(r.Time, h.l.createdFormat) }
	rr.AppendMessage = func(b *Buffer) { b.AppendString(r.Message) }

	for _, attr := range h.attrs {
		rr.with(h.prefix+attr.Key, attr.Value)
	}
	r.Attrs(func(attr slog.Attr) bool {
		rr.with(h.prefix+attr.Key, attr.Value)
		return true
	})

	if r.PC != 0 {
		f, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
		rr.AppendLocation = func(b *Buffer) {
			b.AppendString(f.File).AppendBytes(':').AppendInt(int64(f.Line), 10)
		}
	}

	rr.Output(h.l.Logger(slog2Logs[r.Level]))

	return nil
}

func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var a []slog.Attr
	if len(h.attrs) == 0 {
		a = attrs
	} else {
		a = append(slices.Clip(h.attrs), attrs...)
	}

	return &slogHandler{
		l:      h.l,
		attrs:  a,
		prefix: h.prefix,
	}
}

func (h *slogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	return &slogHandler{
		l:      h.l,
		attrs:  h.attrs,
		prefix: name + "." + h.prefix,
	}
}
