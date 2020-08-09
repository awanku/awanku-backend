package oauth2provider

import (
	"context"

	"github.com/awanku/awanku/pkg/core"
	"golang.org/x/oauth2"
	googleService "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleProvider struct {
	Config *oauth2.Config
}

func (p *GoogleProvider) LoginURL(state string) string {
	return p.Config.AuthCodeURL(state)
}

func (p *GoogleProvider) ExchangeCode(code string) (*core.OauthUserData, error) {
	token, err := p.Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	tokenSource := p.Config.TokenSource(context.Background(), token)
	svc, err := googleService.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, err
	}

	remoteUserData, err := svc.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return nil, err
	}

	userData := core.OauthUserData{
		Provider:   core.OauthProviderGoogle,
		Name:       remoteUserData.Name,
		Email:      remoteUserData.Email,
		Identifier: remoteUserData.Email,
	}
	return &userData, nil
}
