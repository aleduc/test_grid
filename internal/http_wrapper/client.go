package http_wrapper

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Client implements http request logic.
type Client struct {
	client *http.Client
}

func NewClient(requestTimeout time.Duration) *Client {
	client := &http.Client{
		Timeout: requestTimeout,
	}
	return &Client{client: client}
}

// MakeGetRequest sends http GET request with body.
func (c *Client) MakeGetRequest(ctx context.Context, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}
