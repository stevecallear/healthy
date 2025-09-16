package healthy

import (
	"context"
)

type (
	// Metadata represents a check that has metadata
	MetadataCheck interface {
		Check
		Metadata() Metadata
	}

	metadataCheck struct {
		md Metadata
		fn CheckFunc
	}

	// Metadata represents check metadata
	Metadata map[string]any

	metadataContextKey struct{}
)

// WithMetadata wraps the CheckFunc with the supplied metadata
// Metadata should be specified in key/value pairs. The function will
// panic if the pairs are invalid.
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

// Healthy returns nil if the health check is successful
func (c *metadataCheck) Healthy(ctx context.Context) error {
	return c.fn(ctx)
}

// Metadata returns the check metadata
func (c *metadataCheck) Metadata() Metadata {
	return c.md
}

// GetContextMetadata returns the metadata stored in the context
func GetContextMetadata(ctx context.Context) Metadata {
	if m, ok := ctx.Value(metadataContextKey{}).(Metadata); ok {
		return m
	}
	return Metadata{}
}

// SetContextMetadata stores the supplied metadata in the context
func SetContextMetadata(ctx context.Context, m Metadata) context.Context {
	return context.WithValue(ctx, metadataContextKey{}, m)
}

// Set sets the metadata key/value
func (m Metadata) Set(key string, value any) {
	if m == nil {
		return
	}
	m[key] = value
}

// Get gets the metadata value
func (m Metadata) Get(key string) any {
	if m == nil {
		return nil
	}
	return m[key]
}
