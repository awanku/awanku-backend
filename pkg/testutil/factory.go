package testutil

import (
	"context"
	"time"

	"github.com/awanku/awanku/pkg/core"
	"github.com/bxcodec/faker/v3"
)

func UserFactory(ctx context.Context, n int) []*core.User {
	users := []*core.User{}
	for i := 0; i < n; i++ {
		githubUsername := faker.Word() + faker.Word()
		googleEmail := faker.Word() + "_" + faker.Email()
		user := &core.User{
			Name:                faker.Name() + " " + faker.Word(),
			Email:               faker.Word() + "_" + faker.Email(),
			GithubLoginUsername: &githubUsername,
			GoogleLoginEmail:    &googleEmail,
		}
		if err := orm(ctx).Insert(user); err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	return users
}

func OauthAuthorizationCodeFactory(ctx context.Context, userID int64) *core.OauthAuthorizationCode {
	oauthCode := &core.OauthAuthorizationCode{
		UserID:    userID,
		Code:      faker.Word() + faker.Word(),
		ExpiresAt: time.Now().Add(core.OauthAuthorizationCodeMaxDuration),
	}
	if err := orm(ctx).Insert(oauthCode); err != nil {
		panic(err)
	}
	return oauthCode
}

func OauthTokenFactory(ctx context.Context, userID int64, secretKey string) *core.OauthToken {
	token, err := core.BuildOauthToken([]byte(secretKey), 20)
	if err != nil {
		panic(err)
	}
	token.UserID = userID
	token.RequesterIP = "127.0.0.1"
	token.RequesterUserAgent = "testing"
	_, err = orm(ctx).Query(token, `
        insert into oauth_tokens (user_id, access_token_hash, refresh_token_hash, expires_at, requester_ip, requester_user_agent)
        values (?, ?, ?, ?, ?, ?)
        returning id
    `, token.UserID, token.AccessTokenHash, token.RefreshTokenHash, token.ExpiresAt, token.RequesterIP, token.RequesterUserAgent)
	if err != nil {
		panic(err)
	}
	return token
}
