package core

import (
	"fmt"
	"net/http"
	"net/url"

	"trakt-go/util"
)

type ApiResponse struct {
	StatusCode int
	Body       map[string]any
}

type Client struct {
	headers map[string]string
	scheme  string
	baseUrl string
}

func NewClient(appName string, appVersion string, apiKey string, apiVersion int) *Client {
	return &Client{
		map[string]string{
			"Content-Type":      "application/json",
			"User-Agent":        fmt.Sprintf("%s/%s", appName, appVersion),
			"trakt-api-key":     apiKey,
			"trakt-api-version": fmt.Sprintf("%d", apiVersion),
		},
		traktScheme,
		traktBaseUrl,
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

func (c *Client) do(method, path string, queryParams, payload map[string]any) (*ApiResponse, error) {
	serializedPayload, err := util.SerializeRequest(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, c.constructUrl(path, queryParams), serializedPayload)
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

	responseBody, err := util.DeserializeResponse(resp)
	if err != nil {
		return nil, err
	}
	return &ApiResponse{resp.StatusCode, responseBody}, nil
}

func (c *Client) Get(path string, queryParams map[string]any) (*ApiResponse, error) {
	return c.do(http.MethodGet, path, queryParams, nil)
}

func (c *Client) Post(path string, queryParams, payload map[string]any) (*ApiResponse, error) {
	return c.do(http.MethodPost, path, queryParams, payload)
}
