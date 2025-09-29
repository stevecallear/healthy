package healthy_test

import (
	"context"
	"fmt"
	"maps"
	"net"
	"testing"
	"time"

	"github.com/stevecallear/healthy"
)

func TestTCPCheck_Healthy(t *testing.T) {
	addr := fmt.Sprintf("localhost:%d", getFreePort())

	t.Run("should return an error on failure", func(t *testing.T) {
		sut := healthy.TCP(addr).Timeout(time.Millisecond)
		err := sut.Healthy(context.Background())
		if err == nil {
			t.Error("got nil, expected error")
		}
	})

	t.Run("should return nil on success", func(t *testing.T) {
		l, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()

		sut := healthy.TCP(addr)
		err = sut.Healthy(context.Background())
		if err != nil {
			t.Errorf("got %v, expected error", err)
		}
	})
}

func TestTCP_Metadata(t *testing.T) {
	t.Run("should return the check metadata", func(t *testing.T) {
		const target = "localhost:8080"
		exp := healthy.Metadata{"type": "tcp", "target": target, "timeout": "500ms"}
		act := healthy.TCP(target).Timeout(500 * time.Millisecond).Metadata()
		if !maps.Equal(act, exp) {
			t.Errorf("got %v, expected %v", act, exp)
		}
	})
}

func getFreePort() int {
	a, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}
