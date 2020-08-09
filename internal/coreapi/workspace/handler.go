package workspace

import (
	"net/http"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/internal/coreapi/utils/apihelper"
)

// @Id api.v1.workspace.listAll
// @Summary List all workspaces owned by current authenticated user
// @Tags Workspace
// @Security oauthAccessToken
// @Router /v1/workspaces [get]
// @Produce json
// @Success 200 {array} core.Workspace
// @Failure 400 {object} apihelper.HTTPError
// @Failure 401 {object} apihelper.HTTPError
// @Failure 500 {object} apihelper.InternalServerError
func HandleListAll(w http.ResponseWriter, r *http.Request) {
	currentUser := appctx.AuthenticatedUser(r.Context())

	workspaces, err := getUserWorkspaces(r.Context(), currentUser.ID)
	if err != nil {
		apihelper.InternalServerErrResp(w, err)
		return
	}

	apihelper.JSON(w, http.StatusOK, workspaces)
}
