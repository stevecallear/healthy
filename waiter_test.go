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

func TestExecutor_Wait(t *testing.T) {
	t.Run("should return nil for zero checks", func(t *testing.T) {
		err := healthy.New().Wait()
		if err != nil {
			t.Errorf("got %v, expected nil", err)
		}
	})

	t.Run("should error on timeout", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			fn, close := tempFile()
			close() // remove immediately

			err := healthy.New(healthy.File(fn)).Wait()
			if err == nil {
				t.Errorf("got nil, expected error")
			}
		})
	})

	t.Run("should abort on fatal error", func(t *testing.T) {
		err := healthy.New(healthy.CheckFunc(func(ctx context.Context) error {
			return healthy.Fatal(errors.New("fatal"))
		})).Wait()
		if !healthy.IsFatal(err) {
			t.Errorf("got %v, expected fatal error", err)
		}
	})

	t.Run("should require all checks", func(t *testing.T) {
		f1, c1 := tempFile()
		defer c1()

		f2, c2 := tempFile()
		c2() // remove immediately

		synctest.Test(t, func(t *testing.T) {
			err := healthy.New(healthy.File(f1), healthy.File(f2)).Wait()
			if err == nil {
				t.Errorf("got nil, expected error")
			}
		})
	})

	t.Run("should wait for check", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var n int32
			err := healthy.New(healthy.CheckFunc(func(ctx context.Context) error {
				if atomic.LoadInt32(&n) < 1 {
					return errors.New("error")
				}
				return nil
			}),
			).Wait(healthy.WithCallback(func(ctx context.Context, r healthy.Result) {
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

func TestWithContext(t *testing.T) {
	sut := healthy.New(healthy.CheckFunc(func(ctx context.Context) error {
		return errors.New("error")
	}))

	t.Run("should use the context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately

		err := sut.Wait(healthy.WithContext(ctx), healthy.WithTimeout(time.Second))
		if err == nil {
			t.Error("got nil, expected error")
		}
	})

	t.Run("should apply the timeout to the supplied context", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			ctx := context.Background()
			err := sut.Wait(healthy.WithContext(ctx), healthy.WithTimeout(time.Millisecond))
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
		sut := healthy.New(healthy.NewCheck(func(ctx context.Context) error {
			return exp
		}, "type", ctype))

		synctest.Test(t, func(t *testing.T) {
			var res healthy.Result
			mu := new(sync.Mutex)

			sut.Wait(healthy.WithCallback(func(ctx context.Context, r healthy.Result) {
				mu.Lock()
				defer mu.Unlock()
				res = r
			}))

			synctest.Wait()

			if act, exp := res.Info["type"], ctype; act != exp {
				t.Errorf("got %s, expected %s", act, exp)
			}
			if act, exp := res.Err, exp; act != exp {
				t.Errorf("got %v, expected %v", act, exp)
			}
		})
	})
}
