package healthy_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stevecallear/healthy"
)

func TestWait(t *testing.T) {
	t.Run("should return nil for nil check", func(t *testing.T) {
		err := healthy.Wait(nil)
		if err != nil {
			t.Errorf("got %v, expected nil", err)
		}
	})

	t.Run("should error on timeout", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fn, close := tempFile()
			close() // remove immediately

			err := healthy.Wait(healthy.File(fn))
			if err == nil {
				t.Errorf("got nil, expected error")
			}
		})
	})

	t.Run("should abort on fatal error", func(t *testing.T) {
		err := healthy.Wait(healthy.CheckFunc(func(ctx context.Context) error {
			return healthy.Fatal(errors.New("fatal"))
		}))
		if !healthy.IsFatal(err) {
			t.Errorf("got %v, expected fatal error", err)
		}
	})

	t.Run("should wait for check", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var n int32
			err := healthy.Wait(
				healthy.CheckFunc(func(ctx context.Context) error {
					if atomic.LoadInt32(&n) < 1 {
						return errors.New("error")
					}
					return nil
				}),
				healthy.WithCallback(func(ctx context.Context, err error) {
					atomic.AddInt32(&n, 1)
				}))
			if err != nil {
				t.Errorf("got %v, expected nil", err)
			}
			if act, exp := n, int32(2); act != exp {
				t.Errorf("got %d attempts expected %d", act, exp)
			}
		})
	})
}

func TestJoinOptions(t *testing.T) {
	t.Run("should join the options", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		err := errors.New("error")
		var cberr error
		opts := healthy.JoinOptions(
			healthy.WithContext(ctx),
			healthy.WithCallback(func(ctx context.Context, err error) {
				cberr = err
			}),
		)

		synctest.Test(t, func(t *testing.T) {
			act := healthy.Wait(healthy.CheckFunc(func(ctx context.Context) error {
				return err
			}), opts)

			if !errors.Is(act, context.Canceled) || !errors.Is(act, err) {
				t.Errorf("got %v, expected %v and %v", act, context.Canceled, err)
			}

			if cberr != err {
				t.Errorf("got %v, expected %v", cberr, err)
			}
		})
	})
}

func TestWithContext(t *testing.T) {
	check := healthy.CheckFunc(func(ctx context.Context) error {
		return errors.New("error")
	})

	t.Run("should use the context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		err := healthy.Wait(check, healthy.WithContext(ctx), healthy.WithTimeout(time.Second))
		if err == nil {
			t.Error("got nil, expected error")
		}
	})

	t.Run("should apply the timeout to the supplied context", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx := context.Background()
			err := healthy.Wait(check, healthy.WithContext(ctx), healthy.WithTimeout(time.Millisecond))
			if err == nil {
				t.Error("got nil, expected error")
			}
		})
	})

	t.Run("should accept zero timeout", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err := healthy.Wait(check, healthy.WithContext(ctx), healthy.WithTimeout(0)) // use supplied context to control timeout
			if err == nil {
				t.Error("got nil, expected error")
			}
		})
	})
}

func TestWithCallback(t *testing.T) {
	t.Run("should invoke the callback", func(t *testing.T) {
		const ctype = "test"

		exp := errors.New("error")
		check := healthy.WithMetadata(func(ctx context.Context) error {
			return exp
		}, "type", ctype)

		synctest.Test(t, func(t *testing.T) {
			var md healthy.Metadata
			var cerr error

			mu := new(sync.Mutex)
			healthy.Wait(check, healthy.WithCallback(func(ctx context.Context, err error) {
				mu.Lock()
				defer mu.Unlock()
				md = healthy.GetContextMetadata(ctx)
				cerr = err
			}))

			synctest.Wait()

			if act, exp := md["type"], ctype; act != exp {
				t.Errorf("got %s, expected %s", act, exp)
			}
			if act, exp := cerr, exp; act != exp {
				t.Errorf("got %v, expected %v", act, exp)
			}
		})
	})
}
