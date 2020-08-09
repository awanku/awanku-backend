package repository

import (
	"context"
	"fmt"
	"net/http"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/pkg/core"
	"github.com/bradleyfalzon/ghinstallation"
	githubService "github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func fetchGithubIdentifier(ctx context.Context, installationID int64) (string, error) {
	config := appctx.GithubAppConfig(ctx)

	transport, err := ghinstallation.NewAppsTransport(http.DefaultTransport, config.AppID, config.PrivateKey)
	if err != nil {
		return "", err
	}

	client := githubService.NewClient(&http.Client{Transport: transport})

	installation, _, err := client.Apps.GetInstallation(ctx, installationID)
	if err != nil {
		return "", err
	}

	return installation.GetAccount().GetHTMLURL(), nil
}

func fetchRepositories(ctx context.Context, conn *core.RepositoryConnection) ([]*core.Repository, error) {
	switch conn.Provider {
	case core.RepositoryProviderGithubV1:
		installationID, err := fetchGithubAppInstallationID(ctx, conn.Provider, conn.Payload)
		if err != nil {
			return []*core.Repository{}, err
		}
		return fetchGithubRepositories(ctx, installationID)
	}
	return nil, fmt.Errorf("unknown provider: %s", conn.Provider)
}

func fetchGithubAppInstallationID(ctx context.Context, provider core.RepositoryProvider, payload interface{}) (int64, error) {
	switch provider {
	case core.RepositoryProviderGithubV1:
		parsed := core.GithubRepositoryV1Payload{
			InstallationID: payload.(map[string]interface{})["installation_id"].(string),
		}
		installationID, err := parsed.ParseInstallationID()
		if err != nil {
			return 0, fmt.Errorf("failed to parse payload: %s", err)
		}
		return installationID, nil
	}
	return 0, fmt.Errorf("unknown provider: %s", provider)
}

func fetchGithubRepositories(ctx context.Context, installationID int64) ([]*core.Repository, error) {
	config := appctx.GithubAppConfig(ctx)

	transport, err := ghinstallation.NewAppsTransport(http.DefaultTransport, config.AppID, config.PrivateKey)
	if err != nil {
		return []*core.Repository{}, err
	}

	client := githubService.NewClient(&http.Client{Transport: transport})
	token, _, err := client.Apps.CreateInstallationToken(ctx, installationID, &githubService.InstallationTokenOptions{})
	if err != nil {
		return []*core.Repository{}, err
	}

	client = githubService.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.GetToken()})))
	repos, _, err := client.Apps.ListRepos(ctx, nil)
	if err != nil {
		return []*core.Repository{}, err
	}

	var results []*core.Repository
	for _, repo := range repos {
		results = append(results, &core.Repository{
			Name: repo.GetFullName(),
			URL:  repo.GetHTMLURL(),
		})
	}
	return results, nil
}
