package core

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"trakt-go/util"

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
	assert.Equal(t, "https", client.scheme)
	assert.Equal(t, "api.trakt.tv", client.baseUrl)
}

func NewTestClientServer(handler http.HandlerFunc) (*Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	traktBaseUrl = strings.TrimPrefix(server.URL, "http://")
	traktScheme = "http"
	client := NewClient("test", "test", "test", 2)
	return client, server
}

func TestGet(t *testing.T) {
	path := "/path"
	queryParams := map[string]any{"q1": "v1"}
	response := map[string]any{"r1": "v1"}

	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, path, r.URL.Path)
		assert.Equal(t, "q1=v1", r.URL.RawQuery)
		w.WriteHeader(http.StatusOK)
		marshaled, _ := json.Marshal(response)
		w.Write([]byte(marshaled))
	}))
	defer server.Close()

	resp, err := client.Get(path, queryParams)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, response, resp.Body)
}

func TestPost(t *testing.T) {
	path := "/path"
	queryParams := map[string]any{"q1": "v1"}
	payload := map[string]any{"b1": "v1"}
	response := map[string]any{"r1": "v1"}

	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, path, r.URL.Path)
		assert.Equal(t, "q1=v1", r.URL.RawQuery)
		deserializedBody, err := util.Deserialize(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, payload, deserializedBody)
		w.WriteHeader(http.StatusOK)
		marshaled, _ := json.Marshal(response)
		w.Write([]byte(marshaled))
	}))
	defer server.Close()

	resp, err := client.Post(path, queryParams, payload)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, response, resp.Body)
}
