package logger

import (
	"context"
	"log/slog"

	"github.com/yoannduc/go-project-template/pkg/traceid"
	"github.com/yoannduc/go-project-template/pkg/userip"
)

// A ContextHandler wraps a Handler with a Handle method
// that adds new attributes traceid and userip from context.
type ContextHandler struct {
	handler slog.Handler
}

// NewContextHandler returns a ContextHandler.
// All methods delegate to h with added attributes in Handle.
func NewContextHandler(h slog.Handler) *ContextHandler {
	// Optimization: avoid chains of ContextHandler.
	if lh, ok := h.(*ContextHandler); ok {
		h = lh.Handler()
	}

	return &ContextHandler{h}
}

// Enabled implements Handler.Enabled.
func (h *ContextHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.handler.Enabled(ctx, lvl)
}

// Handle implements Handler.Handle by adding new attributes
// traceid and userip from context.
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if tID, ok := traceid.FromContext(ctx); ok {
		r.AddAttrs(slog.String("traceId", tID))
	}
	if ip, ok := userip.FromContext(ctx); ok {
		r.AddAttrs(slog.String("ip", ip.String()))
	}

	return h.handler.Handle(ctx, r)
}

// WithAttrs implements Handler.WithAttrs.
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

// WithGroup implements Handler.WithGroup.
func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

// Handler returns the Handler wrapped by h.
func (h *ContextHandler) Handler() slog.Handler {
	return h.handler
}
