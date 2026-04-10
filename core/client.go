package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Wahaj404/trakt-go/util"
)

type ApiResponse struct {
	StatusCode int
	Body       map[string]any
	Pagination *Pagination // nil if endpoint is not paginated
	RateLimit  *RateLimit  // nil if header absent
}

type Client struct {
	headers map[string]string
	scheme  string
	baseUrl string
}

func NewClient(appName, appVersion, apiKey string, apiVersion int) *Client {
	return newClientWithConfig(appName, appVersion, apiKey, apiVersion, &TraktApiConfig{})
}

func newClientWithConfig(appName, appVersion, apiKey string, apiVersion int, config ITraktApiConfig) *Client {
	return &Client{
		map[string]string{
			"Content-Type":      "application/json",
			"User-Agent":        fmt.Sprintf("%s/%s", appName, appVersion),
			"trakt-api-key":     apiKey,
			"trakt-api-version": fmt.Sprintf("%d", apiVersion),
		},
		config.Scheme(),
		config.BaseUrl(),
	}
}

func (c *Client) constructUrl(path string, queryParams map[string]any) string {
	u := url.URL{
		Scheme: c.scheme,
		Host:   c.baseUrl,
		Path:   path,
	}
	qParams := url.Values{}
	for k, v := range queryParams {
		qParams.Set(k, fmt.Sprintf("%s", v))
	}
	u.RawQuery = qParams.Encode()
	return u.String()
}

func (c *Client) do(ctx context.Context, method, path string, queryParams, payload map[string]any) (*ApiResponse, error) {
	serializedPayload, err := util.Serialize(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, c.constructUrl(path, queryParams), serializedPayload)
	if err != nil {
		return nil, err
	}
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{
			StatusCode:      resp.StatusCode,
			WWWAuthenticate: resp.Header.Get("WWW-Authenticate"),
			RetryAfter:      resp.Header.Get("Retry-After"),
			RawBody:         rawBody,
		}
		var errBody struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		_ = json.Unmarshal(rawBody, &errBody)
		apiErr.Message = errBody.Error
		apiErr.Description = errBody.ErrorDescription
		return nil, apiErr
	}

	responseBody := make(map[string]any)
	if len(rawBody) > 0 {
		_ = json.Unmarshal(rawBody, &responseBody)
	}
	return &ApiResponse{
		StatusCode: resp.StatusCode,
		Body:       responseBody,
		Pagination: parsePagination(resp.Header),
		RateLimit:  parseRateLimit(resp.Header),
	}, nil
}

func (c *Client) Get(ctx context.Context, path string, queryParams map[string]any) (*ApiResponse, error) {
	return c.do(ctx, http.MethodGet, path, queryParams, nil)
}

func (c *Client) Post(ctx context.Context, path string, queryParams, payload map[string]any) (*ApiResponse, error) {
	return c.do(ctx, http.MethodPost, path, queryParams, payload)
}
