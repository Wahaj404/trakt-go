package core

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovies_Popular(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/movies/popular", r.URL.Path)
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"title":"The Dark Knight","year":2008,"ids":{"trakt":4,"slug":"the-dark-knight-2008","imdb":"tt0468569","tmdb":155}},
			{"title":"Inception","year":2010,"ids":{"trakt":417,"slug":"inception-2010","imdb":"tt1375666","tmdb":27205}}
		]`))
	}))
	defer server.Close()

	movies, resp, err := client.Movies.Popular(context.Background(), &ListOptions{Page: 1, Limit: 10})
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Len(t, movies, 2)
	assert.Equal(t, "The Dark Knight", movies[0].Title)
	assert.Equal(t, 2008, movies[0].Year)
	assert.Equal(t, 4, movies[0].IDs.Trakt)
	assert.Equal(t, "the-dark-knight-2008", movies[0].IDs.Slug)
	assert.Equal(t, "tt0468569", movies[0].IDs.IMDB)
	assert.Equal(t, "Inception", movies[1].Title)
}

func TestMovies_Trending(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/movies/trending", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"watchers":123,"movie":{"title":"Dune","year":2021,"ids":{"trakt":316842,"slug":"dune-2021","imdb":"tt1160419","tmdb":438631}}},
			{"watchers":45,"movie":{"title":"Tenet","year":2020,"ids":{"trakt":316841,"slug":"tenet-2020","imdb":"tt6723592","tmdb":577922}}}
		]`))
	}))
	defer server.Close()

	trending, resp, err := client.Movies.Trending(context.Background(), nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Len(t, trending, 2)
	assert.Equal(t, 123, trending[0].Watchers)
	assert.Equal(t, "Dune", trending[0].Movie.Title)
	assert.Equal(t, 2021, trending[0].Movie.Year)
}

func TestMovies_Get(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/movies/tron-legacy-2010", r.URL.Path)
		assert.Equal(t, "full", r.URL.Query().Get("extended"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"title":"TRON: Legacy",
			"year":2010,
			"ids":{"trakt":12601,"slug":"tron-legacy-2010","imdb":"tt1104001","tmdb":20526},
			"tagline":"The Grid. A digital frontier.",
			"overview":"Sam Flynn, a rebellious 27-year-old, is haunted by the mysterious disappearance of his father.",
			"released":"2010-12-17",
			"runtime":125,
			"country":"us",
			"rating":6.8,
			"votes":5000,
			"genres":["action","adventure","science fiction"],
			"certification":"PG"
		}`))
	}))
	defer server.Close()

	movie, resp, err := client.Movies.Get(context.Background(), "tron-legacy-2010", ExtendedFull)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.NotNil(t, movie)
	assert.Equal(t, "TRON: Legacy", movie.Title)
	assert.Equal(t, 2010, movie.Year)
	assert.Equal(t, "tt1104001", movie.IDs.IMDB)
	// Extended=full fields:
	assert.Equal(t, "Sam Flynn, a rebellious 27-year-old, is haunted by the mysterious disappearance of his father.", movie.Overview)
	assert.Equal(t, 125, movie.Runtime)
	assert.Equal(t, 6.8, movie.Rating)
	assert.Equal(t, []string{"action", "adventure", "science fiction"}, movie.Genres)
}

func TestMovies_Get_NotFound(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	movie, _, err := client.Movies.Get(context.Background(), "does-not-exist")
	assert.Nil(t, movie)
	var apiErr *APIError
	assert.True(t, errors.As(err, &apiErr))
	assert.Equal(t, 404, apiErr.StatusCode)
}

func TestMovies_Aliases(t *testing.T) {
	client, server := NewTestClientServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/movies/tron-legacy-2010/aliases", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"title":"TRON: Legacy","country":"us"},
			{"title":"Tron: El Legado","country":"es"}
		]`))
	}))
	defer server.Close()

	aliases, resp, err := client.Movies.Aliases(context.Background(), "tron-legacy-2010")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Len(t, aliases, 2)
	assert.Equal(t, "TRON: Legacy", aliases[0].Title)
	assert.Equal(t, "us", aliases[0].Country)
	assert.Equal(t, "es", aliases[1].Country)
}
