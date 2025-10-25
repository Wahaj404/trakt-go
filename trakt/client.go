package trakt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents the Trakt API client
type Client struct {
	BaseURL    string
	ClientID   string
	HTTPClient *http.Client
}

// NewClient creates a new Trakt API client
func NewClient(baseURL, clientID string) *Client {
	return &Client{
		BaseURL:    baseURL,
		ClientID:   clientID,
		HTTPClient: &http.Client{},
	}
}

// DoRequest performs HTTP requests for Trakt services
func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}, out interface{}) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("trakt-api-key", c.ClientID)
	req.Header.Set("trakt-api-version", "2")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(b))
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return nil
}
