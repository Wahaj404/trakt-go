package core

type ITraktApiConfig interface {
	Scheme() string
	BaseUrl() string
}

type TraktApiConfig struct{}

func (tac *TraktApiConfig) Scheme() string {
	return "https"
}

func (tac *TraktApiConfig) BaseUrl() string {
	return "api.trakt.tv"
}
