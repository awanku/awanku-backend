package workspace

import (
	"context"
	"net/http"
	"strconv"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/internal/coreapi/utils/apihelper"
	"github.com/go-chi/chi"
)

func CurrentWorkspaceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workspaceID := chi.URLParam(r, "workspace_id")
		parsedWorkspaceID, _ := strconv.ParseInt(workspaceID, 10, 64)
		if parsedWorkspaceID <= 0 {
			apihelper.BadRequestErrResp(w, "bad_request", map[string]string{
				"workspace_id": "invalid",
			})
			return
		}

		workspace, err := getWorkspaceByID(r.Context(), parsedWorkspaceID)
		if err != nil {
			apihelper.InternalServerErrResp(w, err)
			return
		}
		if workspace == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), appctx.KeyCurrentWorkspace, workspace)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
