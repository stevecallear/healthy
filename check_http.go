package healthy

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type HTTPCheck struct {
	client     *http.Client
	url        string
	statusCode int
}

func HTTP(url string) *HTTPCheck {
	return &HTTPCheck{
		client:     &http.Client{Timeout: time.Second},
		url:        url,
		statusCode: http.StatusOK,
	}
}

func (c *HTTPCheck) Timeout(t time.Duration) *HTTPCheck {
	c.client.Timeout = t
	return c
}

func (c *HTTPCheck) Expect(statusCode int) *HTTPCheck {
	c.statusCode = statusCode
	return c
}

func (c *HTTPCheck) Healthy(ctx context.Context) error {
	req, err := http.NewRequest(http.MethodGet, c.url, nil)
	if err != nil {
		return Fatal(err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != c.statusCode {
		return fmt.Errorf("incorrect status code: %d", res.StatusCode)
	}

	return nil
}

func (c *HTTPCheck) Metadata() Metadata {
	return Metadata{
		"type":    "http",
		"target":  c.url,
		"timeout": c.client.Timeout.String(),
	}
}
