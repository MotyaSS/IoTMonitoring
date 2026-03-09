package logger

import (
	"context"
	"log/slog"
)

type prefixHandler struct {
	slog.Handler
	prefix string
}

func (h *prefixHandler) Handle(ctx context.Context, r slog.Record) error {
	r.Message = "[" + h.prefix + "] " + r.Message
	return h.Handler.Handle(ctx, r)
}

func (h *prefixHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &prefixHandler{Handler: h.Handler.WithAttrs(attrs), prefix: h.prefix}
}

func (h *prefixHandler) WithGroup(name string) slog.Handler {
	return &prefixHandler{Handler: h.Handler.WithGroup(name), prefix: h.prefix}
}

func NewLoggerWithPrefix(baseHandler slog.Handler, prefix string) *slog.Logger {
	handler := &prefixHandler{
		Handler: baseHandler,
		prefix:  prefix,
	}

	return slog.New(handler)
}
