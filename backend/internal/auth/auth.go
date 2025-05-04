package auth

import (
	"sentimenta/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewOAuthConfig(config config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.GOOGLE_CLIENT_ID,
		ClientSecret: config.GOOGLE_CLIENT_SECRET,
		RedirectURL:  config.GOOGLE_CLIENT_CALLBACK, // должен совпадать с тем, что в Google Console
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}
