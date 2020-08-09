package auth

import (
	"testing"
	"time"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/pkg/core"
	"github.com/awanku/awanku/pkg/testutil"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	ctx, close := testutil.Context()
	defer close()

	t.Run("user does not exists", func(t *testing.T) {
		user, err := getUserByID(ctx, -1000)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("user exists", func(t *testing.T) {
		user := testutil.UserFactory(ctx, 1)[0]

		retrieved, err := getUserByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user, retrieved)
	})
}

func TestGetOrCreateUserByEmail(t *testing.T) {
	ctx, close := testutil.Context()
	ctx = appctx.CreateDatabaseTx(ctx)
	defer close()

	t.Run("new user", func(t *testing.T) {
		githubUsername := faker.Word() + faker.Username()
		googleEmail := faker.Word() + "_" + faker.Email()
		user := &core.User{
			Name:                faker.Name() + " " + faker.Word(),
			Email:               faker.Word() + "_" + faker.Email(),
			GithubLoginUsername: &githubUsername,
			GoogleLoginEmail:    &googleEmail,
		}
		err := getOrCreateUserByEmail(ctx, user)
		assert.NoError(t, err)
		assert.True(t, user.ID > 0)
		assert.NotNil(t, user.CreatedAt)
		assert.Nil(t, user.UpdatedAt)
		assert.Nil(t, user.DeletedAt)
	})

	t.Run("existing user", func(t *testing.T) {
		user := testutil.UserFactory(ctx, 1)[0]
		err := getOrCreateUserByEmail(ctx, user)
		assert.NoError(t, err)
		assert.True(t, user.ID > 0)
		assert.NotNil(t, user.CreatedAt)
		assert.NotNil(t, user.UpdatedAt)
		assert.Nil(t, user.DeletedAt)
	})
}

func TestSaveOauthAuthorizationCode(t *testing.T) {
	ctx, close := testutil.Context()
	ctx = appctx.CreateDatabaseTx(ctx)
	defer close()

	user := testutil.UserFactory(ctx, 1)[0]
	code := faker.Word()
	err := saveOauthAuthorizationCode(ctx, user.ID, code)
	assert.NoError(t, err)
}

func TestGetOauthAuthorizationCodeByCode(t *testing.T) {
	ctx, close := testutil.Context()
	defer close()

	t.Run("code exists", func(t *testing.T) {
		user := testutil.UserFactory(ctx, 1)[0]
		authCode := testutil.OauthAuthorizationCodeFactory(ctx, user.ID)
		retrieved, err := getOauthAuthorizationCodeBycode(ctx, authCode.Code)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, retrieved.UserID)
		assert.Equal(t, authCode.Code, retrieved.Code)

		// another call with same code should not work
		retrieved, err = getOauthAuthorizationCodeBycode(ctx, authCode.Code)
		assert.NoError(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("code does not exists", func(t *testing.T) {
		retrieved, err := getOauthAuthorizationCodeBycode(ctx, "does not exists")
		assert.NoError(t, err)
		assert.Nil(t, retrieved)
	})
}

func TestGetOauthTokenByID(t *testing.T) {
	ctx, close := testutil.Context()
	defer close()

	t.Run("token does not exists", func(t *testing.T) {
		token, err := getOauthTokenByID(ctx, -123)
		assert.NoError(t, err)
		assert.Nil(t, token)
	})

	t.Run("token exists", func(t *testing.T) {
		user := testutil.UserFactory(ctx, 1)[0]
		token := testutil.OauthTokenFactory(ctx, user.ID, "secret")

		retrieved, err := getOauthTokenByID(ctx, token.ID)
		assert.NoError(t, err)
		assert.True(t, retrieved.ID > 0)
		assert.True(t, retrieved.ExpiresAt.After(time.Now()))
		assert.Nil(t, retrieved.DeletedAt)
	})
}

func TestDeleteOauthToken(t *testing.T) {
	ctx, close := testutil.Context()
	defer close()

	t.Run("token does not exists", func(t *testing.T) {
		err := deleteOauthToken(ctx, -123)
		assert.NoError(t, err)
	})

	t.Run("token exists", func(t *testing.T) {
		user := testutil.UserFactory(ctx, 1)[0]
		token := testutil.OauthTokenFactory(ctx, user.ID, "secret")

		err := deleteOauthToken(ctx, token.ID)
		assert.NoError(t, err)

		var retrieved core.OauthToken
		err = appctx.Database(ctx).Query(&retrieved, "select * from oauth_tokens where id = ?", token.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved.DeletedAt)
	})
}

func TestSaveOauthToken(t *testing.T) {
	ctx, close := testutil.Context()
	defer close()

	user := testutil.UserFactory(ctx, 1)[0]

	token, err := core.BuildOauthToken([]byte("secret"), 10)
	assert.NoError(t, err)

	token.UserID = user.ID
	token.RequesterIP = "127.0.0.1"
	token.RequesterUserAgent = "testing"
	err = saveOauthToken(ctx, token)
	assert.NoError(t, err)

	var retrieved core.OauthToken
	err = appctx.Database(ctx).Query(&retrieved, "select * from oauth_tokens where id = ?", token.ID)
	assert.NoError(t, err)
	assert.Equal(t, retrieved.ID, token.ID)
}
