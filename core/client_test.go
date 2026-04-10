package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wahaj404/trakt-go/util"

	"github.com/stretchr/testify/assert"
)

type MockConfig struct {
	serverUrl string
}

func (mc *MockConfig) Scheme() string {
	return "http"
}

func (mc *MockConfig) BaseUrl() string {
	return strings.TrimPrefix(mc.serverUrl, fmt.Sprintf("%s://", mc.Scheme()))
}

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
	client := newClientWithConfig("test", "test", "test", 2, &MockConfig{server.URL})
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

	resp, err := client.Get(context.Background(), path, queryParams)
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

	resp, err := client.Post(context.Background(), path, queryParams, payload)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, response, resp.Body)
}

func TestAPIError_EmptyBody(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	resp, err := client.Get(context.Background(), "/missing", nil)
	assert.Nil(t, resp)
	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 404, apiErr.StatusCode)
	assert.Equal(t, "", apiErr.Message)
	assert.Equal(t, "", apiErr.Description)
}

func TestAPIError_MinimalErrorBody(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	_, err := client.Get(context.Background(), "/missing", nil)
	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 404, apiErr.StatusCode)
	assert.Equal(t, "not found", apiErr.Message)
}

func TestAPIError_WWWAuthenticateHeader(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Bearer realm="trakt", error="invalid_token"`)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err := client.Get(context.Background(), "/protected", nil)
	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 401, apiErr.StatusCode)
	assert.Equal(t, `Bearer realm="trakt", error="invalid_token"`, apiErr.WWWAuthenticate)
}

func TestAPIError_RetryAfterHeader(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "42")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	_, err := client.Get(context.Background(), "/busy", nil)
	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 429, apiErr.StatusCode)
	assert.Equal(t, "42", apiErr.RetryAfter)
}

func TestPagination_AllHeaders(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Pagination-Page", "2")
		w.Header().Set("X-Pagination-Limit", "10")
		w.Header().Set("X-Pagination-Page-Count", "5")
		w.Header().Set("X-Pagination-Item-Count", "47")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	resp, err := client.Get(context.Background(), "/list", nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp.Pagination)
	assert.Equal(t, 2, resp.Pagination.Page)
	assert.Equal(t, 10, resp.Pagination.Limit)
	assert.Equal(t, 5, resp.Pagination.PageCount)
	assert.Equal(t, 47, resp.Pagination.ItemCount)
}

func TestPagination_PartialHeaders(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Pagination-Page-Count", "3")
		w.Header().Set("X-Pagination-Item-Count", "12")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	resp, err := client.Get(context.Background(), "/list", nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp.Pagination)
	assert.Equal(t, 0, resp.Pagination.Page)
	assert.Equal(t, 0, resp.Pagination.Limit)
	assert.Equal(t, 3, resp.Pagination.PageCount)
	assert.Equal(t, 12, resp.Pagination.ItemCount)
}

func TestPagination_Absent(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	resp, err := client.Get(context.Background(), "/single", nil)
	assert.Nil(t, err)
	assert.Nil(t, resp.Pagination)
}

func TestRateLimit_Present(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Ratelimit", `{"name":"UNAUTHED_API_GET_LIMIT","period":300,"limit":1000,"remaining":999,"until":"2020-10-27T20:23:16Z"}`)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	resp, err := client.Get(context.Background(), "/thing", nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp.RateLimit)
	assert.Equal(t, "UNAUTHED_API_GET_LIMIT", resp.RateLimit.Name)
	assert.Equal(t, 300, resp.RateLimit.Period)
	assert.Equal(t, 1000, resp.RateLimit.Limit)
	assert.Equal(t, 999, resp.RateLimit.Remaining)
	expectedUntil, _ := time.Parse(time.RFC3339, "2020-10-27T20:23:16Z")
	assert.Equal(t, expectedUntil, resp.RateLimit.Until)
}

func TestRateLimit_Absent(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	resp, err := client.Get(context.Background(), "/thing", nil)
	assert.Nil(t, err)
	assert.Nil(t, resp.RateLimit)
}

func TestAPIError_OAuthErrorShape(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid_grant","error_description":"The provided authorization grant is invalid."}`))
	}))
	defer server.Close()

	_, err := client.Post(context.Background(), "/oauth/token", nil, map[string]any{"grant_type": "authorization_code"})
	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 400, apiErr.StatusCode)
	assert.Equal(t, "invalid_grant", apiErr.Message)
	assert.Equal(t, "The provided authorization grant is invalid.", apiErr.Description)
}

type testThing struct {
	Name string `json:"name"`
}

func TestGetInto(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/thing", r.URL.Path)
		w.Header().Set("X-Pagination-Page", "1")
		w.Header().Set("X-Pagination-Limit", "10")
		w.Header().Set("X-Pagination-Page-Count", "1")
		w.Header().Set("X-Pagination-Item-Count", "1")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"widget"}`))
	}))
	defer server.Close()

	var thing testThing
	resp, err := client.GetInto(context.Background(), "/thing", nil, &thing)
	assert.Nil(t, err)
	assert.Equal(t, "widget", thing.Name)
	assert.Nil(t, resp.Body) // Body is nil when typed unmarshal is used.
	assert.NotNil(t, resp.Pagination)
	assert.Equal(t, 1, resp.Pagination.Page)
}

func TestGetInto_Error(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer server.Close()

	thing := testThing{Name: "unchanged"}
	resp, err := client.GetInto(context.Background(), "/missing", nil, &thing)
	assert.Nil(t, resp)
	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 404, apiErr.StatusCode)
	assert.Equal(t, "not found", apiErr.Message)
	// out should not be mutated on error.
	assert.Equal(t, "unchanged", thing.Name)
}

func TestGetInto_NilOut(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"widget"}`))
	}))
	defer server.Close()

	resp, err := client.GetInto(context.Background(), "/thing", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPostInto(t *testing.T) {
	payload := map[string]any{"b1": "v1"}

	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/thing", r.URL.Path)
		deserializedBody, err := util.Deserialize(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, payload, deserializedBody)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"widget"}`))
	}))
	defer server.Close()

	var thing testThing
	resp, err := client.PostInto(context.Background(), "/thing", nil, payload, &thing)
	assert.Nil(t, err)
	assert.Equal(t, "widget", thing.Name)
	assert.Nil(t, resp.Body) // Body is nil when typed unmarshal is used.
}
