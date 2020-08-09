package workspace

import (
	"context"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/pkg/core"
)

func getWorkspaceByID(ctx context.Context, id int64) (*core.Workspace, error) {
	db := appctx.Database(ctx)

	var query = `
        select *
        from workspaces
        where id = ? and deleted_at is null
    `
	var workspace core.Workspace
	err := db.Query(&workspace, query, id)
	if err != nil {
		return nil, err
	}
	if workspace.ID == 0 {
		return nil, nil
	}
	return &workspace, nil
}

func getUserWorkspaces(ctx context.Context, userID int64) ([]*core.Workspace, error) {
	db := appctx.Database(ctx)

	var query = `
        select workspaces.*
        from workspaces
        join workspace_users on workspaces.id = workspace_users.workspace_id
        where
            workspace_users.user_id = ?
            and workspaces.deleted_at is null
            and workspace_users.deleted_at is null
    `
	var workspaces []*core.Workspace
	err := db.Query(&workspaces, query, userID)
	if err != nil {
		return []*core.Workspace{}, err
	}
	return workspaces, nil
}
