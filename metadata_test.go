package healthy_test

import (
	"context"
	"maps"
	"testing"

	"github.com/stevecallear/healthy"
)

func TestWithMetadata(t *testing.T) {
	tests := []struct {
		name  string
		fn    healthy.CheckFunc
		md    []any
		exp   healthy.Metadata
		panic bool
	}{
		{
			name:  "should panic on nil fn",
			fn:    nil,
			panic: true,
		},
		{
			name:  "should panic on invalid pairs",
			fn:    healthy.CheckFunc(func(ctx context.Context) error { return nil }),
			md:    []any{"a", true, "c"},
			panic: true,
		},
		{
			name:  "should panic on invalid key type",
			fn:    healthy.CheckFunc(func(ctx context.Context) error { return nil }),
			md:    []any{1, "b"},
			panic: true,
		},
		{
			name: "should use the metadata",
			fn:   healthy.CheckFunc(func(ctx context.Context) error { return nil }),
			md:   []any{"a", "b", "c", 1},
			exp:  healthy.Metadata{"a": "b", "c": 1},
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

			c := healthy.WithMetadata(tt.fn, tt.md...)
			if tt.panic {
				return
			}

			if act := c.Metadata(); !maps.Equal(act, tt.exp) {
				t.Errorf("got %v, expected %v", act, tt.exp)
			}
			if err := c.Healthy(context.Background()); err != nil {
				t.Errorf("got %v, expected nil", err)
			}
		})
	}
}

func TestMetadata_Set(t *testing.T) {
	tests := []struct {
		name  string
		sut   healthy.Metadata
		key   string
		value any
		exp   any
	}{
		{
			name:  "should do nothing if the metadata is nil",
			sut:   nil,
			key:   "key",
			value: "value",
			exp:   nil,
		},
		{
			name:  "should store the value",
			sut:   healthy.Metadata{},
			key:   "key",
			value: "value",
			exp:   "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sut.Set(tt.key, tt.value)
			act := tt.sut[tt.key]
			if act != tt.exp {
				t.Errorf("got %s, expected %s", act, tt.exp)
			}
		})
	}
}

func TestMetadata_Get(t *testing.T) {
	tests := []struct {
		name string
		sut  healthy.Metadata
		key  string
		exp  any
	}{
		{
			name: "should return nil if the metadata is nil",
			sut:  nil,
			key:  "key",
			exp:  nil,
		},
		{
			name: "should store the value",
			sut:  healthy.Metadata{"key": "value"},
			key:  "key",
			exp:  "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := tt.sut.Get(tt.key)
			if act != tt.exp {
				t.Errorf("got %s, expected %s", act, tt.exp)
			}
		})
	}
}
