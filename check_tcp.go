package healthy

import (
	"context"
	"net"
	"time"
)

func TCP(addr string) *TCPCheck {
	return &TCPCheck{
		addr:    addr,
		timeout: time.Second,
	}
}

func (c *TCPCheck) Timeout(t time.Duration) *TCPCheck {
	c.timeout = t
	return c
}

func (c *TCPCheck) Healthy(ctx context.Context) error {
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return err
	}

	defer conn.Close()
	return nil
}

func (c *TCPCheck) Info() Info {
	return Info{
		"type":    "tcp",
		"target":  c.addr,
		"timeout": c.timeout.String(),
	}
}
