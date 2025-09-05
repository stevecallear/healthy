package healthy_test

import (
	"context"
	"maps"
	"testing"

	"github.com/stevecallear/healthy"
)

func TestNewCheck(t *testing.T) {
	tests := []struct {
		name  string
		fn    healthy.CheckFunc
		info  []string
		exp   healthy.Info
		panic bool
	}{
		{
			name:  "should panic on nil fn",
			fn:    nil,
			panic: true,
		},
		{
			name:  "should panic on invalid info",
			fn:    healthy.CheckFunc(func(ctx context.Context) error { return nil }),
			info:  []string{"a", "b", "c"},
			panic: true,
		},
		{
			name: "should use the info",
			fn:   healthy.CheckFunc(func(ctx context.Context) error { return nil }),
			info: []string{"a", "b", "c", "d"},
			exp:  healthy.Info{"a": "b", "c": "d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				err := recover()
				if err != nil && !tt.panic {
					t.Errorf("got panic %v, expected nil", err)
				}
				if err == nil && tt.panic {
					t.Errorf("got nil, expected panic")
				}
			}()

			c := healthy.NewCheck(tt.fn, tt.info...)
			if tt.panic {
				return
			}

			if act := c.Info(); !maps.Equal(act, tt.exp) {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
			if err := c.Healthy(context.Background()); err != nil {
				t.Errorf("got %v, expected nil", err)
			}
		})
	}
}
