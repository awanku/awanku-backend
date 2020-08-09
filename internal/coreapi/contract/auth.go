package contract

import (
	"github.com/awanku/awanku/pkg/core"
)

type AuthProvider interface {
	LoginURL(state string) string
	ExchangeCode(code string) (*core.OauthUserData, error)
}

type AuthStore interface {
	CreateOauthAuthorizationCode(userID int64, code string) (*core.OauthAuthorizationCode, error)
	GetOauthAuthorizationCodeByCode(code string) (*core.OauthAuthorizationCode, error)
	CreateOauthToken(token *core.OauthToken) error
	GetOauthTokenByID(id int64) (*core.OauthToken, error)
	DeleteOauthToken(id int64) error
}
