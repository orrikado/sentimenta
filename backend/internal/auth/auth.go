package auth

import (
	"sentimenta/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type OAuth struct {
	GoogleConfig *oauth2.Config
	GithubConfig *oauth2.Config
}

func newGoogleOAuthConfig(config *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GOOGLE_CLIENT_ID,
		ClientSecret: config.GOOGLE_CLIENT_SECRET,
		RedirectURL:  config.GOOGLE_CLIENT_CALLBACK,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func newGithubOAuthConfig(config *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GITHUB_CLIENT_ID,
		ClientSecret: config.GITHUB_CLIENT_SECRET,
		RedirectURL:  config.GITHUB_CLIENT_CALLBACK,
		Scopes:       []string{"email", "profile"},
		Endpoint:     github.Endpoint,
	}
}

func NewOAuth(config *config.Config) *OAuth {
	return &OAuth{
		GoogleConfig: newGoogleOAuthConfig(config),
		GithubConfig: newGithubOAuthConfig(config),
	}
}
