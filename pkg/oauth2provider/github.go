package oauth2provider

import (
	"context"

	"github.com/awanku/awanku/pkg/core"
	githubService "github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

type GithubProvider struct {
	Config *oauth2.Config
}

func (p *GithubProvider) LoginURL(state string) string {
	return p.Config.AuthCodeURL(state)
}

func (p *GithubProvider) ExchangeCode(code string) (*core.OauthUserData, error) {
	token, err := p.Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	tokenSource := p.Config.TokenSource(context.Background(), token)
	client := oauth2.NewClient(context.Background(), tokenSource)
	svc := githubService.NewClient(client)

	githubUser, _, err := svc.Users.Get(context.Background(), "")
	if err != nil {
		return nil, err
	}

	emails, _, err := svc.Users.ListEmails(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	var primaryEmail string
	for _, email := range emails {
		if email.GetEmail() != "" && email.GetVerified() {
			primaryEmail = email.GetEmail()
			break
		}
	}

	// if we don't get primary email, just get first email
	if primaryEmail == "" {
		primaryEmail = emails[0].GetEmail()
	}

	userData := core.OauthUserData{
		Provider:   core.OauthProviderGithub,
		Name:       githubUser.GetName(),
		Email:      primaryEmail,
		Identifier: githubUser.GetLogin(),
	}
	return &userData, nil
}
