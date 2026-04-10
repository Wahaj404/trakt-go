package core

import "fmt"

type APIError struct {
	StatusCode      int
	Message         string // from response body "error" field, or empty
	Description     string // from response body "error_description" field (OAuth), or empty
	WWWAuthenticate string // from WWW-Authenticate header on 401, or empty
	RetryAfter      string // from Retry-After header on 429, or empty
	RawBody         []byte // preserved raw body for debugging
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("trakt: %d %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("trakt: %d", e.StatusCode)
}
