package healthy

import (
	"context"
)

type (
	MetadataCheck interface {
		Check
		Metadata() Metadata
	}

	metadataCheck struct {
		md Metadata
		fn CheckFunc
	}

	Metadata map[string]any

	metadataContextKey struct{}
)

func WithMetadata(fn CheckFunc, pairs ...any) MetadataCheck {
	if fn == nil {
		panic("fn must not be nil")
	}
	if len(pairs)%2 != 0 {
		panic("info must be key value pairs")
	}

	c := &metadataCheck{
		md: map[string]any{},
		fn: fn,
	}

	for i := 0; i < len(pairs); i += 2 {
		k, ok := pairs[i].(string)
		if !ok {
			panic("key must be a string")
		}
		c.md[k] = pairs[i+1]
	}

	return c
}

func (c *metadataCheck) Healthy(ctx context.Context) error {
	return c.fn(ctx)
}

func (c *metadataCheck) Metadata() Metadata {
	return c.md
}

func GetContextMetadata(ctx context.Context) Metadata {
	if m, ok := ctx.Value(metadataContextKey{}).(Metadata); ok {
		return m
	}
	return Metadata{}
}

func SetContextMetadata(ctx context.Context, m Metadata) context.Context {
	return context.WithValue(ctx, metadataContextKey{}, m)
}

func (m Metadata) Set(key string, value any) {
	if m == nil {
		return
	}
	m[key] = value
}

func (m Metadata) Get(key string) any {
	if m == nil {
		return nil
	}
	return m[key]
}
