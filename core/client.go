package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

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
		"https",
		"api.trakt.tv",
	}
}

func (c *Client) constructUrl(path string, queryParams map[string]any) string {
	u := url.URL{
		Scheme: c.scheme,
		Host: c.baseUrl,
		Path: path,
	}
	qParams := url.Values{}
	for k, v := range queryParams {
		qParams.Set(k, fmt.Sprintf("%s", v))
	}
	u.RawQuery = qParams.Encode()
	return u.String()
}

func (c *Client) do(method, path string, queryParams, payload map[string]any) (map[string]any, error) {
	serializedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	payloadBuffer := bytes.NewBuffer(serializedPayload)

	req, err := http.NewRequest(method, c.constructUrl(path, queryParams), payloadBuffer)
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
	responseBuffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	responseBody := make(map[string]any)
	json.Unmarshal(responseBuffer, &responseBody)

	return responseBody, nil
}

func (c *Client) Get(path string, queryParams map[string]any) (map[string]any, error) {
	return c.do(http.MethodGet, path, queryParams, nil)
}

func (c *Client) Post(path string, queryParams, payload map[string]any) (map[string]any, error) {
	return c.do(http.MethodPost, path, queryParams, payload)
}
