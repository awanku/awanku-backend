package user

import (
	"net/http"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/internal/coreapi/utils/apihelper"
)

// @Id api.v1.users.getMe
// @Summary Get current user data
// @Tags Users
// @Security oauthAccessToken
// @Router /v1/users/me [get]
// @Produce json
// @Success 200 {object} core.User
// @Failure 400 {object} apihelper.HTTPError
// @Failure 401 {object} apihelper.HTTPError
// @Failure 500 {object} apihelper.InternalServerError
func HandleGetMe(w http.ResponseWriter, r *http.Request) {
	user := appctx.AuthenticatedUser(r.Context())
	apihelper.JSON(w, http.StatusOK, user)
}
