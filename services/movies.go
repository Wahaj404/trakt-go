package services

import (
	"context"
	"fmt"
	"trakt-go/trakt"
)

type MoviesService struct {
	client *trakt.Client
}

func NewMoviesService(c *trakt.Client) *MoviesService {
	return &MoviesService{client: c}
}

func (m *MoviesService) GetMovie(ctx context.Context, id string) (*trakt.Movie, error) {
	var movie trakt.Movie
	path := fmt.Sprintf("/movies/%s?extended=full", id)
	if err := m.client.DoRequest(ctx, "GET", path, nil, &movie); err != nil {
		return nil, err
	}
	return &movie, nil
}
