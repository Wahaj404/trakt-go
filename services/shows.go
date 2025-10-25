package services

import (
	"context"
	"fmt"
	"trakt-go/trakt"
)

type ShowsService struct {
	client *trakt.Client
}

func NewShowsService(c *trakt.Client) *ShowsService {
	return &ShowsService{client: c}
}

func (m *MoviesService) GetShows(ctx context.Context, id string) (*trakt.Show, error) {
	var movie trakt.Show
	path := fmt.Sprintf("/shows/%s?extended=full", id)
	if err := m.client.DoRequest(ctx, "GET", path, nil, &movie); err != nil {
		return nil, err
	}
	return &movie, nil
}
