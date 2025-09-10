package healthy

import (
	"context"
	"errors"
	"maps"
	"math/rand/v2"
	"time"

	"golang.org/x/sync/errgroup"
)

type (
	Waiter struct {
		checks []Check
	}

	Result struct {
		Attempt int
		Info    Info
		Err     error
	}

	CallbackFunc func(ctx context.Context, r Result)

	Option func(*options)

	options struct {
		ctx      context.Context
		timeout  time.Duration
		delay    time.Duration
		jitter   time.Duration
		callback CallbackFunc
	}
)

var defaultOptions = options{
	timeout: 30 * time.Second,
	delay:   time.Second,
	jitter:  100 * time.Millisecond,
}

func New(c ...Check) *Waiter {
	return &Waiter{
		checks: c,
	}
}

func (w *Waiter) Wait(opts ...Option) error {
	if len(w.checks) < 1 {
		return nil
	}

	o := w.getOptions(opts)
	ctx, cancel := o.contextWithCancel()
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	for _, cc := range w.checks {
		g.Go(func() error {
			attempt := 1

			var info Info
			if ic, ok := cc.(InfoCheck); ok {
				info = maps.Clone(ic.Info())
			}

			for {
				err := cc.Healthy(ctx)
				if o.callback != nil {
					// callback must be synchronous to avoid data race on err
					o.callback(ctx, Result{
						Err:     err,
						Info:    info,
						Attempt: attempt,
					})
				}
				if err == nil {
					return nil
				}

				if IsFatal(err) {
					return err
				}

				select {
				case <-time.After(o.calculateDelay()):
				case <-ctx.Done():
					return errors.Join(context.Cause(ctx), err)
				}
				attempt++
			}
		})
	}

	return g.Wait()
}

func (w *Waiter) getOptions(opts []Option) options {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	return o
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
