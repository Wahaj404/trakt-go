package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListOptions_Nil(t *testing.T) {
	var o *ListOptions
	assert.Nil(t, o.toQuery())
}

func TestListOptions_Empty(t *testing.T) {
	// Documented choice: zero-value struct returns nil (same as nil receiver),
	// so callers can pass the result straight through to Client.GetInto.
	o := &ListOptions{}
	assert.Nil(t, o.toQuery())
}

func TestListOptions_PageOnly(t *testing.T) {
	o := &ListOptions{Page: 2}
	assert.Equal(t, map[string]any{"page": 2}, o.toQuery())
}

func TestListOptions_Full(t *testing.T) {
	o := &ListOptions{
		Page:     3,
		Limit:    25,
		Extended: []ExtendedLevel{ExtendedFull, ExtendedEpisodes},
	}
	assert.Equal(t, map[string]any{
		"page":     3,
		"limit":    25,
		"extended": "full,episodes",
	}, o.toQuery())
}
