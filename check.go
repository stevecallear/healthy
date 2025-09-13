package healthy

import (
	"context"
)

type (
	Check interface {
		Healthy(ctx context.Context) error
	}

	CheckFunc func(ctx context.Context) error
)

func (c CheckFunc) Healthy(ctx context.Context) error {
	return c(ctx)
}
