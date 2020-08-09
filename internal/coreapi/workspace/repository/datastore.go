package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/pkg/core"
)

var errConnectionAlreadyExists = errors.New("repository connection already exists")

func saveRepositoryConnection(ctx context.Context, conn *core.RepositoryConnection) error {
	db := appctx.Database(ctx)

	var query = `
        insert into workspace_repository_connections (workspace_id, identifier, provider, payload, created_at)
        values (?, ?, ?, ?, now())
    `
	err := db.WriterExec(query, conn.WorkspaceID, conn.Identifier, conn.Provider, conn.Payload)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return errConnectionAlreadyExists
		}
	}
	return err
}

func getConnections(ctx context.Context, workspaceID int64) ([]*core.RepositoryConnection, error) {
	db := appctx.Database(ctx)

	var query = `
        select *
        from workspace_repository_connections
        where workspace_id = ? and deleted_at is null
    `
	var conns []*core.RepositoryConnection
	err := db.Query(&conns, query, workspaceID)
	if err != nil {
		return []*core.RepositoryConnection{}, nil
	}
	return conns, nil
}
