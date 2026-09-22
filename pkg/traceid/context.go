// Package traceid provides functions for associating a trace id to a request
// or retrieve it from a request if a traceid is already present in headers,
// and associating it with a Context.
package traceid

import (
	"context"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// traceIDKey is the context key for the trace id. Its value of zero is
// arbitrary. If this package defined other context keys, they would have
// different integer values.
const traceIDKey key = 0

// NewContext returns a new Context carrying trace id.
func NewContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// FromContext extracts the trace id from ctx, if present.
func FromContext(ctx context.Context) (string, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the string type assertion returns ok=false for nil.
	traceID, ok := ctx.Value(traceIDKey).(string)
	return traceID, ok
}
