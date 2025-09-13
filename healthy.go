package healthy

import (
	"context"
	"errors"
	"maps"
	"math/rand/v2"
	"strconv"
	"time"
)

type (
	Option func(*options)

	options struct {
		ctx      context.Context
		timeout  time.Duration
		delay    time.Duration
		jitter   time.Duration
		callback CallbackFunc
	}

	CallbackFunc func(ctx context.Context, err error)
)

const mdKeyAttempt = "attempt"

var defaultOptions = options{
	timeout: 30 * time.Second,
	delay:   time.Second,
	jitter:  100 * time.Millisecond,
}

func Wait(c Check, opts ...Option) error {
	if c == nil {
		return nil
	}

	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	ctx, cancel := o.contextWithCancel()
	defer cancel()

	attempt := 1
	md := GetContextMetadata(ctx)
	if mc, ok := c.(MetadataCheck); ok {
		maps.Copy(md, mc.Metadata())
	}

	for {
		md.Set(mdKeyAttempt, strconv.Itoa(attempt))
		mdctx := SetContextMetadata(ctx, md)

		err := c.Healthy(mdctx)
		if o.callback != nil {
			o.callback(mdctx, err)
		}
		if IsFatal(err) {
			return err
		}
		if err == nil {
			return nil
		}

		select {
		case <-time.After(o.calculateDelay()):
		case <-ctx.Done():
			return errors.Join(context.Cause(ctx), err)
		}
		attempt++
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithDelay(d time.Duration) Option {
	return func(o *options) {
		o.delay = d
	}
}

func WithTimeout(t time.Duration) Option {
	return func(o *options) {
		o.timeout = t
	}
}

func WithJitter(j time.Duration) Option {
	return func(o *options) {
		o.jitter = j
	}
}

func WithCallback(fn CallbackFunc) Option {
	return func(o *options) {
		o.callback = fn
	}
}

func (o options) contextWithCancel() (context.Context, func()) {
	ctx := o.ctx
	cancel := func() {}

	if ctx == nil {
		ctx = context.Background()
	}
	if o.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, o.timeout)
	}

	return ctx, cancel
}

func (o options) calculateDelay() time.Duration {
	if o.jitter < 1 {
		return o.delay
	}

	j := rand.Int64N(int64(o.jitter))
	return o.delay + time.Duration(j)
}
