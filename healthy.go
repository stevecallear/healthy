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
	// Option represents an execution option.
	Option func(*options)

	options struct {
		ctx      context.Context
		timeout  time.Duration
		delay    time.Duration
		jitter   time.Duration
		callback CallbackFunc
	}

	// CallbackFunc represents an execution callback function.
	// The function is invoked on each health check invocation.
	CallbackFunc func(ctx context.Context, err error)
)

const mdKeyAttempt = "attempt"

var defaultOptions = options{
	timeout: 30 * time.Second,
	delay:   time.Second,
}

// Wait executes the check using the supplied options.
// Checks are retried until successful execution or option limits are reached.
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

// JoinOptions joins the specified options to simplify re-use.
func JoinOptions(opts ...Option) Option {
	return func(o *options) {
		for _, opt := range opts {
			opt(o)
		}
	}
}

// WithContext uses the supplied context for check execution.
// This allows alternative context cancellations to be specified.
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithTimeout specifies the retry execution timeout.
// Retry execution will be cancelled either on the cancellation of a
// supplied context or after the timeout has elapsed.
// The default value is 30 seconds.
func WithTimeout(t time.Duration) Option {
	return func(o *options) {
		o.timeout = t
	}
}

// WithDelay specifies the retry delay between check executions.
// The default value is one second.
func WithDelay(d time.Duration) Option {
	return func(o *options) {
		o.delay = d
	}
}

// WithJitter specifies an optional maximum jitter to apply to the delay.
// The default value is no jitter.
func WithJitter(j time.Duration) Option {
	return func(o *options) {
		o.jitter = j
	}
}

// WithCallback specifies the callback function to be invoked after check execution.
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
