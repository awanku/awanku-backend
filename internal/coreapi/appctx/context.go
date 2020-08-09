package appctx

import (
	"context"

	hansip "github.com/asasmoyo/pq-hansip"
	"github.com/awanku/awanku/pkg/core"
)

// Key context key
type Key string

// context keys
const (
	KeyEnvironment       Key = "environment"
	KeyDatabase          Key = "database"
	KeyDatabaseTx        Key = "database_transaction"
	KeyAuthenticatedUser Key = "authenticated_user"
	KeyCurrentWorkspace  Key = "current_workspace"
	KeyGithubAppConfig   Key = "github_app_config"
)

// Environment fetch environment name from context
func Environment(ctx context.Context) string {
	raw := ctx.Value(KeyEnvironment)
	if val, ok := raw.(string); ok {
		return val
	}
	return ""
}

// Database fetch database instance from context
func Database(ctx context.Context) *hansip.Cluster {
	raw := ctx.Value(KeyDatabase)
	if val, ok := raw.(*hansip.Cluster); ok {
		return val
	}
	return nil
}

// CreateDatabaseTx creates new context with database transaction
func CreateDatabaseTx(ctx context.Context) context.Context {
	tx, err := Database(ctx).NewTransaction()
	if err != nil {
		return nil
	}
	return context.WithValue(ctx, KeyDatabaseTx, tx)
}

// DatabaseTx fetch database transaction from context
func DatabaseTx(ctx context.Context) hansip.Transaction {
	raw := ctx.Value(KeyDatabaseTx)
	if val, ok := raw.(hansip.Transaction); ok {
		return val
	}
	return nil
}

// AuthenticatedUser fetch authenticated user from context
func AuthenticatedUser(ctx context.Context) *core.User {
	raw := ctx.Value(KeyAuthenticatedUser)
	if val, ok := raw.(*core.User); ok {
		return val
	}
	return nil
}

// CurrentWorkspace fetch current workspace from context
func CurrentWorkspace(ctx context.Context) *core.Workspace {
	raw := ctx.Value(KeyCurrentWorkspace)
	if val, ok := raw.(*core.Workspace); ok {
		return val
	}
	return nil
}

// GithubAppConfig fetch github app config from context
func GithubAppConfig(ctx context.Context) *core.GithubAppConfig {
	raw := ctx.Value(KeyGithubAppConfig)
	if val, ok := raw.(*core.GithubAppConfig); ok {
		return val
	}
	return nil
}
