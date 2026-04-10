package core

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Pagination struct {
	Page      int // 0 if header absent
	Limit     int // 0 if header absent
	PageCount int
	ItemCount int
}

type RateLimit struct {
	Name      string    `json:"name"`
	Period    int       `json:"period"`
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Until     time.Time `json:"until"`
}

// parsePagination returns a *Pagination if any of the four X-Pagination-*
// headers is present on the response, otherwise nil. Individual malformed
// header values are silently skipped.
func parsePagination(h http.Header) *Pagination {
	pageStr := h.Get("X-Pagination-Page")
	limitStr := h.Get("X-Pagination-Limit")
	pageCountStr := h.Get("X-Pagination-Page-Count")
	itemCountStr := h.Get("X-Pagination-Item-Count")

	if pageStr == "" && limitStr == "" && pageCountStr == "" && itemCountStr == "" {
		return nil
	}

	p := &Pagination{}
	if v, err := strconv.Atoi(pageStr); err == nil {
		p.Page = v
	}
	if v, err := strconv.Atoi(limitStr); err == nil {
		p.Limit = v
	}
	if v, err := strconv.Atoi(pageCountStr); err == nil {
		p.PageCount = v
	}
	if v, err := strconv.Atoi(itemCountStr); err == nil {
		p.ItemCount = v
	}
	return p
}

// parseRateLimit returns a *RateLimit parsed from the X-Ratelimit header,
// which carries a JSON object. Returns nil if the header is absent or
// malformed.
func parseRateLimit(h http.Header) *RateLimit {
	raw := h.Get("X-Ratelimit")
	if raw == "" {
		return nil
	}
	var rl RateLimit
	if err := json.Unmarshal([]byte(raw), &rl); err != nil {
		return nil
	}
	return &rl
}
