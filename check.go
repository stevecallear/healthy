package healthy

import (
	"context"
	"net/http"
	"time"
)

type (
	Check interface {
		Healthy(ctx context.Context) error
	}

	InfoCheck interface {
		Check
		Info() Info
	}

	Info map[string]string

	check struct {
		info Info
		fn   CheckFunc
	}

	CheckFunc func(ctx context.Context) error

	TCPCheck struct {
		addr    string
		timeout time.Duration
	}

	HTTPCheck struct {
		client     *http.Client
		url        string
		statusCode int
	}
)

func NewCheck(fn CheckFunc, info ...string) InfoCheck {
	if fn == nil {
		panic("fn must not be nil")
	}
	if len(info)%2 != 0 {
		panic("info must be key value pairs")
	}

	c := &check{
		info: map[string]string{},
		fn:   fn,
	}

	for i := 0; i < len(info); i += 2 {
		c.info[info[i]] = info[i+1]
	}

	return c
}

func (c CheckFunc) Healthy(ctx context.Context) error {
	return c(ctx)
}

func (c *check) Healthy(ctx context.Context) error {
	return c.fn(ctx)
}

func (c *check) Info() Info {
	return c.info
}
