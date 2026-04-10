package core

import "time"

// IDs is the standard ID set Trakt attaches to every media object.
type IDs struct {
	Trakt int    `json:"trakt"`
	Slug  string `json:"slug"`
	IMDB  string `json:"imdb,omitempty"`
	TMDB  int    `json:"tmdb,omitempty"`
	TVDB  int    `json:"tvdb,omitempty"`
}

// Movie is the full Trakt movie object. Fields marked "extended=full" are
// only populated when the request uses ?extended=full.
type Movie struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
	IDs   IDs    `json:"ids"`

	// Extended=full fields (zero if not requested):
	Tagline               string    `json:"tagline,omitempty"`
	Overview              string    `json:"overview,omitempty"`
	Released              string    `json:"released,omitempty"` // YYYY-MM-DD
	Runtime               int       `json:"runtime,omitempty"`
	Country               string    `json:"country,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
	Trailer               string    `json:"trailer,omitempty"`
	Homepage              string    `json:"homepage,omitempty"`
	Status                string    `json:"status,omitempty"`
	Rating                float64   `json:"rating,omitempty"`
	Votes                 int       `json:"votes,omitempty"`
	CommentCount          int       `json:"comment_count,omitempty"`
	Language              string    `json:"language,omitempty"`
	Languages             []string  `json:"languages,omitempty"`
	AvailableTranslations []string  `json:"available_translations,omitempty"`
	Genres                []string  `json:"genres,omitempty"`
	Certification         string    `json:"certification,omitempty"`
}

// TrendingMovie is the wrapper shape returned by /movies/trending.
type TrendingMovie struct {
	Watchers int   `json:"watchers"`
	Movie    Movie `json:"movie"`
}

// Alias is the shape returned by /movies/:id/aliases.
type Alias struct {
	Title   string `json:"title"`
	Country string `json:"country"`
}
