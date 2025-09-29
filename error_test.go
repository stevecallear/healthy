package healthy_test

import (
	"errors"
	"testing"

	"github.com/stevecallear/healthy"
)

func TestFatal(t *testing.T) {
	t.Run("should return a fatal error", func(t *testing.T) {
		err := healthy.Fatal(errors.New("error"))
		if !healthy.IsFatal(err) {
			t.Errorf("got %v, expected fatal error", err)
		}
	})
}

func TestFatalError_Error(t *testing.T) {
	t.Run("should return the inner error", func(t *testing.T) {
		const exp = "error"
		act := healthy.Fatal(errors.New(exp)).Error()
		if act != exp {
			t.Errorf("got %s, expected %s", act, exp)
		}
	})
}

func TestFatalError_Unwrap(t *testing.T) {
	t.Run("should return the inner error", func(t *testing.T) {
		exp := errors.New("error")
		act := healthy.Fatal(exp).(interface{ Unwrap() error }).Unwrap()
		if act != exp {
			t.Errorf("got %v, expected %v", act, exp)
		}
	})
}
