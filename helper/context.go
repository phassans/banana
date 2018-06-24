package helper

import (
	"context"
)

func NewContext() context.Context {
	return context.Background()
}

type (
	fieldKey   struct{ string }
	contextKey struct{}
)

// EmbedContext embeds one context inside another.
func EmbedContext(embedding, embedded context.Context) context.Context {
	return context.WithValue(embedding, contextKey{}, embedded)
}

// GetContext returns the embedded context if exists, otherwise return the embedding context.
func GetContext(ctx context.Context) context.Context {
	if embedded := ctx.Value(contextKey{}); embedded != nil {
		return embedded.(context.Context)
	}
	return ctx
}

// WithValue adds a new value into the context.
func WithValue(ctx context.Context, field string, val interface{}) context.Context {
	embedded := GetContext(ctx)
	newEmbedded := context.WithValue(embedded, fieldKey{field}, val)
	if embedded != ctx {
		return EmbedContext(ctx, newEmbedded)
	}
	return newEmbedded
}

// GetValue extracts a value from the context.
func GetValue(ctx context.Context, field string) interface{} {
	return GetContext(ctx).Value(fieldKey{field})
}

// GetAPIVersion extracts the API version from a context.
func GetContextValue(ctx context.Context, field string) string {
	if v := GetValue(ctx, field); v != nil {
		return v.(string)
	}
	return ""
}
