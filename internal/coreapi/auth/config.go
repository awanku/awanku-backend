package auth

import (
	"github.com/awanku/awanku/pkg/core"
	"github.com/awanku/awanku/pkg/oauth2provider"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
	oauth2google "golang.org/x/oauth2/google"
)

type oauthProvider interface {
	LoginURL(state string) string
	ExchangeCode(code string) (*core.OauthUserData, error)
}

func oauth2Config(environment, provider string) *oauth2.Config {
	data := map[string]map[string]*oauth2.Config{
		"production": {
			core.OauthProviderGithub: &oauth2.Config{
				ClientID:     "857f091db02afc686f98",
				ClientSecret: "7641cdcea107b66a962ac73988b5d77bd2efe13c",
				Scopes:       []string{"read:user", "user:email"},
				Endpoint:     oauth2github.Endpoint,
				RedirectURL:  "https://api.awanku.id/v1/auth/github/callback",
			},
			core.OauthProviderGoogle: &oauth2.Config{
				ClientID:     "757848106543-b069r475lcql7373vmhk3179u5l1anek.apps.googleusercontent.com",
				ClientSecret: "R_JRM20ol-YFqzbVilo81sey",
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
				Endpoint:    oauth2google.Endpoint,
				RedirectURL: "https://api.awanku.id/v1/auth/google/callback",
			},
		},
		"development": {
			core.OauthProviderGithub: &oauth2.Config{
				ClientID:     "6b068bb4d449eb24b8d8",
				ClientSecret: "32118588e79d9132c1c5fa36ec3ad2fcc73bb453",
				Scopes:       []string{"read:user", "user:email"},
				Endpoint:     oauth2github.Endpoint,
				RedirectURL:  "http://api.dev.awanku.xyz/v1/auth/github/callback",
			},
			core.OauthProviderGoogle: &oauth2.Config{
				ClientID:     "757848106543-7joqgt09qgmmvt131is9b5i62bcqd2co.apps.googleusercontent.com",
				ClientSecret: "8mqdZmkeP3O5fkZbEVJdOR05",
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					"https://www.googleapis.com/auth/userinfo.profile",
				},
				Endpoint:    oauth2google.Endpoint,
				RedirectURL: "http://api.dev.awanku.xyz/v1/auth/google/callback",
			},
		},
	}
	return data[environment][provider]
}

func oauth2Provider(provider, environment string) oauthProvider {
	switch provider {
	case core.OauthProviderGoogle:
		return &oauth2provider.GoogleProvider{
			Config: oauth2Config(environment, core.OauthProviderGoogle),
		}
	case core.OauthProviderGithub:
		return &oauth2provider.GithubProvider{
			Config: oauth2Config(environment, core.OauthProviderGithub),
		}
	}
	return nil
}
