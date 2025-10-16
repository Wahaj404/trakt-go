package core

import "fmt"

type Client struct {
	headers map[string]string
}

func MakeClient(appName string, appVersion string, apiKey string, apiVersion int) *Client {
	return &Client{
		map[string]string{
			"Content-Type":      "application/json",
			"User-Agent":        fmt.Sprintf("%s/%s", appName, appVersion),
			"trakt-api-key":     apiKey,
			"trakt-api-version": fmt.Sprintf("%d", apiVersion),
		},
	}
}
