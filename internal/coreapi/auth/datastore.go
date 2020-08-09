package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/pkg/core"
)

var errOauthTokenExpired = errors.New("oauth token expired")

func getUserByID(ctx context.Context, id int64) (*core.User, error) {
	db := appctx.Database(ctx)

	var query = `
        select *
        from users
        where id = ? and deleted_at is null
    `
	var returned core.User
	err := db.Query(&returned, query, id)
	if err != nil {
		return nil, err
	}
	if returned.ID == 0 {
		return nil, nil
	}
	return &returned, nil
}

func getOauthAuthorizationCodeBycode(ctx context.Context, code string) (*core.OauthAuthorizationCode, error) {
	db := appctx.Database(ctx)

	var query = `
        delete from oauth_authorization_codes
        where code = ? and expires_at > now()
        returning *
    `
	var returned core.OauthAuthorizationCode
	err := db.Query(&returned, query, code)
	if err != nil {
		return nil, err
	}
	if returned.UserID == 0 || returned.Code == "" {
		return nil, nil
	}
	return &returned, nil
}

func getOauthTokenByID(ctx context.Context, id int64) (*core.OauthToken, error) {
	db := appctx.Database(ctx)

	var query = `
        select *
        from oauth_tokens
        where id = ? and deleted_at is null
    `
	var returned core.OauthToken
	err := db.Query(&returned, query, id)
	if err != nil {
		return nil, err
	}
	if returned.ID == 0 {
		return nil, nil
	}
	if returned.ExpiresAt.Before(time.Now()) {
		return nil, errOauthTokenExpired
	}
	return &returned, nil
}

func deleteOauthToken(ctx context.Context, id int64) error {
	db := appctx.Database(ctx)

	var query = `
        update oauth_tokens
        set deleted_at = now()
        where id = ?
    `
	return db.WriterExec(query, id)
}

func saveOauthToken(ctx context.Context, token *core.OauthToken) error {
	db := appctx.Database(ctx)

	var query = `
        insert into oauth_tokens (user_id, access_token_hash, refresh_token_hash, expires_at, requester_ip, requester_user_agent)
        values (?, ?, ?, ?, ?, ?)
        returning id
    `
	var returned struct {
		ID int64
	}
	err := db.WriterQuery(&returned, query, token.UserID, token.AccessTokenHash, token.RefreshTokenHash, token.ExpiresAt, token.RequesterIP, token.RequesterUserAgent)
	if err != nil {
		return err
	}
	token.ID = returned.ID
	return nil
}

func registerUser(ctx context.Context, user *core.User, authorizationCode string) error {
	ctx = appctx.CreateDatabaseTx(ctx)
	defer func() {
		err := appctx.DatabaseTx(ctx).Commit()
		if err != nil {
			fmt.Println("transaction commit failed:", err)
		}
	}()

	if err := getOrCreateUserByEmail(ctx, user); err != nil {
		return err
	}

	// create workspace on new user
	// updated_at value is nil on new user
	if user.UpdatedAt == nil {
		if err := createWorkspace(ctx, user); err != nil {
			return err
		}
	}

	return saveOauthAuthorizationCode(ctx, user.ID, authorizationCode)
}

func getOrCreateUserByEmail(ctx context.Context, user *core.User) (err error) {
	tx := appctx.DatabaseTx(ctx)
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var query = `
        insert into users (name, email, google_login_email, github_login_username)
        values (?, ?, ?, ?)
        on conflict (email) do update set updated_at = now()
        returning id, created_at, updated_at
    `
	var returned struct {
		ID        int64
		CreatedAt time.Time
		UpdatedAt *time.Time
	}
	err = tx.Query(&returned, query, user.Name, user.Email, user.GoogleLoginEmail, user.GithubLoginUsername)
	if err != nil {
		return err
	}

	user.ID = returned.ID
	user.CreatedAt = returned.CreatedAt
	user.UpdatedAt = returned.UpdatedAt
	return nil
}

func createWorkspace(ctx context.Context, user *core.User) (err error) {
	tx := appctx.DatabaseTx(ctx)
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var name = fmt.Sprintf("%s's workspace", user.Name)
	var queryWorkspace = `
        insert into workspaces (name, created_at)
        values (?, now())
        returning id
    `
	var returnedWorkspace struct{ ID int64 }
	err = tx.Query(&returnedWorkspace, queryWorkspace, name)
	if err != nil {
		return err
	}

	var queryWorkspaceUser = `
        insert into workspace_users (workspace_id, user_id, access_level, created_at)
        values (?, ?, 'owner', now())
    `
	err = tx.Exec(queryWorkspaceUser, returnedWorkspace.ID, user.ID)
	return
}

func saveOauthAuthorizationCode(ctx context.Context, userID int64, code string) (err error) {
	tx := appctx.DatabaseTx(ctx)
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var query = `
        insert into oauth_authorization_codes (user_id, code, expires_at)
        values (?, ?, now() + interval '5 minutes')
    `
	err = tx.Exec(query, userID, code)
	return
}
