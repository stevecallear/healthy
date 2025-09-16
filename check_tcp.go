package healthy

import (
	"context"
	"net"
	"time"
)

// TCPCheck represents a TCP health check
type TCPCheck struct {
	addr    string
	timeout time.Duration
}

// TCP returns a TCP health check
func TCP(addr string) *TCPCheck {
	return &TCPCheck{
		addr:    addr,
		timeout: time.Second,
	}
}

// Timeout species the TCP dial timeout
func (c *TCPCheck) Timeout(t time.Duration) *TCPCheck {
	c.timeout = t
	return c
}

// Healtyh returns nil if a TCP connection can be established with
// the target address
func (c *TCPCheck) Healthy(ctx context.Context) error {
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return err
	}

	defer conn.Close()
	return nil
}

// Metadata retuns the check metadata
func (c *TCPCheck) Metadata() Metadata {
	return Metadata{
		"type":    "tcp",
		"target":  c.addr,
		"timeout": c.timeout.String(),
	}
}
