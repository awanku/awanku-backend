package coreapi

import (
	"io/ioutil"

	"github.com/awanku/awanku/pkg/core"
	"github.com/caarlos0/env"
)

type Config struct {
	Environment             string `env:"ENVIRONMENT"`
	DatabaseURL             string `env:"DATABASE_URL"`
	OAuthSecretKey          string `env:"OAUTH_SECRET_KEY"`
	GithubAppID             int64  `env:"GITHUB_APP_ID"`
	GithubAppPrivateKeyPath string `env:"GITHUB_APP_PRIVATE_KEY_PATH"`
	GithubAppInstallURL     string `env:"GITHUB_APP_INSTALL_URL"`
}

func (c *Config) Load() error {
	return env.Parse(c)
}

func (c *Config) GithubAppConfig() (*core.GithubAppConfig, error) {
	privateKey, err := ioutil.ReadFile(c.GithubAppPrivateKeyPath)
	if err != nil {
		return nil, err
	}

	config := core.GithubAppConfig{
		AppID:      c.GithubAppID,
		PrivateKey: privateKey,
		InstallURL: c.GithubAppInstallURL,
	}
	return &config, nil
}
