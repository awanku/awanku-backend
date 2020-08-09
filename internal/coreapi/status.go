package coreapi

import (
	"net/http"

	hansip "github.com/asasmoyo/pq-hansip"
	"github.com/awanku/awanku/internal/coreapi/utils/apihelper"
)

type statusResponse struct {
	Database hansip.ClusterHealth
}

// @Id api.status
// @Summary Get API health status
// @Router /status [get]
// @Produce json
// @Success 200 {object} statusResponse
func statusHandler(db *hansip.Cluster) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apihelper.JSON(w, http.StatusOK, db.Health())
	}
}
