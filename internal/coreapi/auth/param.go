package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

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
		validation.Field(&p.Provider, validation.Required, validation.In(core.OauthProviderGithub, core.OauthProviderGoogle)),
		validation.Field(&p.RedirectTo, validation.Required, is.URL),
	)
}

func (p getProviderConnectParam) encodeState() (string, error) {
	payload := map[string]string{
		"redirect_to": p.RedirectTo,
	}
	marshalled, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(marshalled)
	return encoded, nil
}

type getProviderCallbackParam struct {
	Provider string `json:"provider" swaggerignore:"true"`
	Code     string `json:"code" validate:"required"`
	State    string `json:"state" validate:"required"`
}

func (p getProviderCallbackParam) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Provider, validation.Required, validation.In("github", "google")),
		validation.Field(&p.Code, validation.Required),
		validation.Field(&p.State, validation.Required),
	)
}

func (p getProviderCallbackParam) decodeState() (map[string]string, error) {
	decodedState, err := base64.URLEncoding.DecodeString(p.State)
	if err != nil {
		return nil, err
	}

	var decodedData map[string]string
	err = json.Unmarshal(decodedState, &decodedData)
	if err != nil {
		return nil, err
	}
	return decodedData, nil
}

type postTokenParam struct {
	GrantType    string `json:"grant_type" schema:"grant_type" validate:"required"`
	Code         string `json:"code" schema:"code"`
	RefreshToken string `json:"refresh_token" schema:"refresh_token"`

	retrievedCode       *core.OauthAuthorizationCode
	retrievedOauthToken *core.OauthToken

	oauthTokenSecretKey []byte
}

func (p *postTokenParam) Validate(ctx context.Context) error {
	return validation.ValidateStruct(p,
		validation.Field(&p.GrantType, validation.Required, validation.In("authorization_code", "refresh_token")),
		validation.Field(&p.Code, validation.By(p.validateCode(ctx))),
		validation.Field(&p.RefreshToken, validation.By(p.validateRefreshToken(ctx))),
	)
}

func (p *postTokenParam) validateCode(ctx context.Context) validation.RuleFunc {
	return func(value interface{}) error {
		if p.GrantType != "authorization_code" {
			return nil
		}

		code, _ := value.(string)
		if p.GrantType == "authorization_code" && value == "" {
			return validation.ErrRequired
		}

		var err error
		p.retrievedCode, err = getOauthAuthorizationCodeBycode(ctx, code)
		if err != nil {
			return validation.NewInternalError(err)
		}

		if p.retrievedCode != nil && p.retrievedCode.Code != "" {
			return nil
		}
		return errors.New("invalid")
	}
}

func (p *postTokenParam) validateRefreshToken(ctx context.Context) validation.RuleFunc {
	return func(value interface{}) error {
		if p.GrantType != "refresh_token" {
			return nil
		}

		tokenRaw, _ := value.(string)
		if p.GrantType == "refresh_token" && value == "" {
			return validation.ErrRequired
		}

		tokenParts := strings.Split(tokenRaw, ":")
		if len(tokenParts) != 2 {
			return errors.New("invalid")
		}

		tokenIDStr := tokenParts[0]
		refreshTokenDecoded, err := base64.URLEncoding.DecodeString(tokenParts[1])
		if err != nil {
			return errors.New("invalid")
		}

		tokenID, _ := strconv.ParseInt(tokenIDStr, 10, 64)
		p.retrievedOauthToken, err = getOauthTokenByID(ctx, tokenID)
		if err != nil {
			return validation.NewInternalError(err)
		}

		valid, err := core.ValidateHMAC(p.oauthTokenSecretKey, refreshTokenDecoded, p.retrievedOauthToken.RefreshTokenHash)
		if err != nil {
			return validation.NewInternalError(err)
		}
		if valid {
			return nil
		}
		return errors.New("invalid")
	}
}
