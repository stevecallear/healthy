package healthy_test

import (
	"context"
	"fmt"
	"maps"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stevecallear/healthy"
)

func TestHTTPCheck_Healthy(t *testing.T) {
	addr := fmt.Sprintf(":%d", getFreePort())
	url := "http://localhost" + addr

	t.Run("should return an error on failure", func(t *testing.T) {
		sut := healthy.HTTP(url).Timeout(time.Millisecond)
		err := sut.Healthy(context.Background())
		if err == nil {
			t.Error("got nil, expected error")
		}
	})

	t.Run("should return fatal error on invalid request", func(t *testing.T) {
		sut := healthy.HTTP("%")
		err := sut.Healthy(context.Background())
		if !healthy.IsFatal(err) {
			t.Errorf("got %v, expected fatal error", err)
		}
	})

	t.Run("should return error on incorrect status code", func(t *testing.T) {
		close := startHTTPImmediate(addr)
		defer close()

		sut := healthy.HTTP(url).Expect(http.StatusFound)
		err := sut.Healthy(context.Background())
		if err == nil {
			t.Error("got nil, expected error")
		}
	})

	t.Run("should return nil on success", func(t *testing.T) {
		close := startHTTPImmediate(addr)
		defer close()

		sut := healthy.HTTP(url)
		err := sut.Healthy(context.Background())
		if err != nil {
			t.Errorf("got %v, expected error", err)
		}
	})
}

func TestHTTP_Metadata(t *testing.T) {
	t.Run("should return the check metadata", func(t *testing.T) {
		const target = "http://localhost:8080s"
		exp := healthy.Metadata{"type": "http", "target": target, "timeout": "500ms"}
		act := healthy.HTTP(target).Timeout(500 * time.Millisecond).Metadata()
		if !maps.Equal(act, exp) {
			t.Errorf("got %v, expected %v", act, exp)
		}
	})
}

func startHTTPDelayed(addr string, delay time.Duration) func() error {
	s := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	go func() {
		<-time.After(delay)
		s.ListenAndServe()
	}()

	return s.Close
}

func startHTTPImmediate(addr string) func() error {
	s := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	go func() {
		s.ListenAndServe()
	}()

	// readiness check
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	for {
		l, err := net.DialTimeout("tcp", addr, 10*time.Millisecond)
		if err == nil {
			defer l.Close()
			break
		}

		select {
		case <-time.After(10 * time.Millisecond):
		case <-ctx.Done():
			panic("failed to start server")
		}
	}

	return s.Close
}
