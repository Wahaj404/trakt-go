package core

import (
	"context"
	"fmt"
)

type MoviesResource struct {
	client *Client
}

// Popular returns the most popular movies.
// GET /movies/popular
func (m *MoviesResource) Popular(ctx context.Context, opts *ListOptions) ([]Movie, *ApiResponse, error) {
	var out []Movie
	resp, err := m.client.GetInto(ctx, "/movies/popular", opts.toQuery(), &out)
	return out, resp, err
}

// Trending returns the most-watched movies right now, each wrapped with a
// watcher count.
// GET /movies/trending
func (m *MoviesResource) Trending(ctx context.Context, opts *ListOptions) ([]TrendingMovie, *ApiResponse, error) {
	var out []TrendingMovie
	resp, err := m.client.GetInto(ctx, "/movies/trending", opts.toQuery(), &out)
	return out, resp, err
}

// Get returns a single movie summary by id (trakt id, slug, or imdb id).
// GET /movies/:id
func (m *MoviesResource) Get(ctx context.Context, id string, extended ...ExtendedLevel) (*Movie, *ApiResponse, error) {
	var out Movie
	opts := &ListOptions{Extended: extended}
	resp, err := m.client.GetInto(ctx, fmt.Sprintf("/movies/%s", id), opts.toQuery(), &out)
	if err != nil {
		return nil, resp, err
	}
	return &out, resp, nil
}

// Aliases returns all title aliases for a movie.
// GET /movies/:id/aliases
func (m *MoviesResource) Aliases(ctx context.Context, id string) ([]Alias, *ApiResponse, error) {
	var out []Alias
	resp, err := m.client.GetInto(ctx, fmt.Sprintf("/movies/%s/aliases", id), nil, &out)
	return out, resp, err
}
