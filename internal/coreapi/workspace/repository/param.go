package repository

import (
	"context"
	"fmt"

	"github.com/awanku/awanku/pkg/core"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type getProviderConnectParam struct {
	Provider   string `json:"provider" validate:"required" swaggerignore:"true"`
	RedirectTo string `json:"redirect_to" validate:"required"`
}

func (p getProviderConnectParam) Validate(ctx context.Context) error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Provider, validation.Required, validation.In(core.OauthProviderGithub)),
		validation.Field(&p.RedirectTo, validation.Required, is.URL),
	)
}

type saveRepositoryConnectionParam struct {
	Provider string `json:"provider" validate:"required" swaggerignore:"true"`
	Code     string `json:"code" validate:"required"`
}

func (p saveRepositoryConnectionParam) Validate(ctx context.Context) error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Provider, validation.Required, validation.In(core.OauthProviderGithub)),
		validation.Field(&p.Code, validation.Required),
	)
}

type saveGithubConnection struct {
	Provider       core.RepositoryProvider `json:"provider" validate:"required" swaggerignore:"true"`
	InstallationID int64                   `json:"installation_id" validate:"required"`
	Action         string                  `json:"action" validate:"required"`
}

func (c saveGithubConnection) Validate(ctx context.Context) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Provider, validation.Required),
		validation.Field(&c.InstallationID, validation.Required, validation.Min(0).Exclusive()),
		validation.Field(&c.Action, validation.Required, validation.In("install", "update")),
	)
}

func (c saveGithubConnection) ParseInstallationID() string {
	return fmt.Sprintf("%d", c.InstallationID)
}
