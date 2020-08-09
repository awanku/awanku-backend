package core

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

// supported oauth providers
const (
	OauthProviderGithub               = "github"
	OauthProviderGoogle               = "google"
	OauthAuthorizationCodeMaxDuration = 5 * time.Minute
)

// OauthUserData represents user data provided by third party oauth services
type OauthUserData struct {
	Provider   string `json:"provider"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Identifier string `json:"identifier"`
}

// OauthAuthorizationCode represents oauth authorization code
type OauthAuthorizationCode struct {
	Code      string
	UserID    int64
	ExpiresAt time.Time
}

// OauthToken represents oauth token
type OauthToken struct {
	ID                 int64
	UserID             int64
	AccessToken        []byte
	AccessTokenHash    []byte
	RefreshToken       []byte
	RefreshTokenHash   []byte
	ExpiresAt          time.Time
	RequesterIP        string
	RequesterUserAgent string
	DeletedAt          *time.Time
}

// Token returns standar token representation
func (t *OauthToken) Token() *oauth2.Token {
	encodedAccessToken := base64.URLEncoding.EncodeToString(t.AccessToken)
	encodedRefreshToken := base64.URLEncoding.EncodeToString(t.RefreshToken)
	return &oauth2.Token{
		AccessToken:  fmt.Sprintf("%d:%s", t.ID, encodedAccessToken),
		RefreshToken: fmt.Sprintf("%d:%s", t.ID, encodedRefreshToken),
		Expiry:       t.ExpiresAt,
		TokenType:    "bearer",
	}
}

// User represents User
type User struct {
	ID                  int64      `json:"id"`
	Name                string     `json:"name"`
	Email               string     `json:"email"`
	GoogleLoginEmail    *string    `json:"google_login_email"`
	GithubLoginUsername *string    `json:"github_login_username"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"-"`
}

// SetOauth2Identifier sets identifier based on provider
func (u *User) SetOauth2Identifier(provider string, identifier *string) {
	switch provider {
	case OauthProviderGithub:
		u.GithubLoginUsername = identifier
	case OauthProviderGoogle:
		u.GoogleLoginEmail = identifier
	}
}

// Workspace represents workspace
type Workspace struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// RepositoryProvider represents repository provider
type RepositoryProvider string

// repository providers
var (
	RepositoryProviderGithubV1 RepositoryProvider = "github-v1"
)

// RepositoryConnection represents workspace repository connection
// TODO: provider is shown with version in get repository connections endpoint, this should not happen
type RepositoryConnection struct {
	ID          int64              `json:"id"`
	WorkspaceID int64              `json:"-"`
	Identifier  string             `json:"identifier"`
	Provider    RepositoryProvider `json:"provider"`
	Payload     interface{}        `json:"-"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   *time.Time         `json:"updated_at"`
	DeletedAt   *time.Time         `json:"-"`
}

// GithubRepositoryV1Payload represents github repository payload.
// This will be stored as JSON in database, this causes int64 to be treated as float64 when
// it is readed back into interface{} type.
// Therefore we store int64 as string.
type GithubRepositoryV1Payload struct {
	InstallationID string `json:"installation_id"`
}

func (p *GithubRepositoryV1Payload) ParseInstallationID() (int64, error) {
	parsed, err := strconv.ParseInt(p.InstallationID, 10, 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
