package healthy

import (
	"context"
)

type (
	// Check represents a health check.
	Check interface {
		Healthy(ctx context.Context) error
	}

	// CheckFunc represents a health check function.
	CheckFunc func(ctx context.Context) error
)

// Healthy should return nil if the health check is successful.
func (c CheckFunc) Healthy(ctx context.Context) error {
	return c(ctx)
}
