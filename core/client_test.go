package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeClient(t *testing.T) {
	client := NewClient("appName", "appVersion", "apiKey", 2)
	assert.Equal(
		t,
		map[string]string{
			"Content-Type":      "application/json",
			"User-Agent":        "appName/appVersion",
			"trakt-api-key":     "apiKey",
			"trakt-api-version": "2",
		},
		client.headers,
	)
}
