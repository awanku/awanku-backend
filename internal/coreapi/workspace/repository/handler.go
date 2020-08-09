package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/internal/coreapi/utils/apihelper"
	"github.com/awanku/awanku/pkg/core"
)

// @Id api.v1.workspace.repository.listAll
// @Summary List all repositories in a workspace
// @Tags Workspace
// @Security oauthAccessToken
// @Param workspace_id path integer true "Workspace id"
// @Router /v1/workspaces/{workspace_id}/repositories [get]
// @Produce json
// @Success 200 {array} core.Repository
// @Failure 400 {object} apihelper.HTTPError
// @Failure 401 {object} apihelper.HTTPError
// @Failure 404
// @Failure 500 {object} apihelper.InternalServerError
func HandleListAllRepositories(w http.ResponseWriter, r *http.Request) {
	currentWorkspace := appctx.CurrentWorkspace(r.Context())

	conns, err := getConnections(r.Context(), currentWorkspace.ID)
	if err != nil {
		apihelper.InternalServerErrResp(w, err)
		return
	}

	var wg sync.WaitGroup
	var mut sync.Mutex
	var repositories []*core.Repository
	var errors []error
	for _, conn := range conns {
		wg.Add(1)
		go func(c *core.RepositoryConnection) {
			reps, err := fetchRepositories(r.Context(), c)
			mut.Lock()
			if err != nil {
				errors = append(errors, err)
			}
			repositories = append(repositories, reps...)
			mut.Unlock()
			wg.Done()
		}(conn)
	}
	wg.Wait()

	if len(errors) > 0 {
		apihelper.InternalServerErrResp(w, fmt.Errorf("failed fetching workspace repositories: %v", errors))
		return
	}

	if repositories == nil {
		repositories = []*core.Repository{}
	}
	apihelper.JSON(w, http.StatusOK, repositories)
}

// @Id api.v1.workspace.connections.listAll
// @Summary List all repository connections in a workspace
// @Tags Workspace
// @Security oauthAccessToken
// @Param workspace_id path integer true "Workspace id"
// @Router /v1/workspaces/{workspace_id}/connections [get]
// @Produce json
// @Success 200 {array} core.RepositoryConnection
// @Failure 400 {object} apihelper.HTTPError
// @Failure 401 {object} apihelper.HTTPError
// @Failure 404
// @Failure 500 {object} apihelper.InternalServerError
func HandleListAllConnections(w http.ResponseWriter, r *http.Request) {
	currentWorkspace := appctx.CurrentWorkspace(r.Context())

	conns, err := getConnections(r.Context(), currentWorkspace.ID)
	if err != nil {
		apihelper.InternalServerErrResp(w, err)
		return
	}

	if conns == nil {
		conns = []*core.RepositoryConnection{}
	}
	apihelper.JSON(w, http.StatusOK, conns)
}

// @Id api.v1.workspace.repository.provider.github.connect
// @Summary Start connecting Github repository
// @Tags Workspace
// @Security oauthAccessToken
// @Param workspace_id path integer true "Workspace id"
// @Router /v1/workspaces/{workspace_id}/providers/github [get]
// @Produce json
// @Success 301
// @Header 301 {string} location "setup github connection url"
// @Failure 400 {object} apihelper.HTTPError
// @Failure 401 {object} apihelper.HTTPError
// @Failure 404
// @Failure 500 {object} apihelper.InternalServerError
func HandleConnectGithub(w http.ResponseWriter, r *http.Request) {
	url := appctx.GithubAppConfig(r.Context()).InstallURL
	apihelper.RedirectResp(w, url)
}

// @Id api.v1.workspace.repository.provider.github.save
// @Summary Save Github repository connection
// @Tags Workspace
// @Security oauthAccessToken
// @Param workspace_id path integer true "Workspace id"
// @Router /v1/workspaces/{workspace_id}/providers/github [post]
// @Produce json
// @Success 201
// @Failure 400 {object} apihelper.HTTPError
// @Failure 401 {object} apihelper.HTTPError
// @Failure 404
// @Failure 500 {object} apihelper.InternalServerError
func HandleSaveGithubConnection(w http.ResponseWriter, r *http.Request) {
	currentWorkspace := appctx.CurrentWorkspace(r.Context())

	param := saveGithubConnection{
		Provider: core.RepositoryProviderGithubV1,
	}
	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		apihelper.BadRequestErrResp(w, "invalid_request", map[string]string{
			"request_body": "malformed format",
		})
		return
	}
	if err := param.Validate(r.Context()); err != nil {
		apihelper.ValidationErrResp(w, err)
		return
	}

	payload := &core.GithubRepositoryV1Payload{InstallationID: param.ParseInstallationID()}

	identifier, err := fetchGithubIdentifier(r.Context(), param.InstallationID)
	if err != nil {
		apihelper.ValidationErrResp(w, map[string]string{
			"installation_id": "invalid",
		})
		return
	}

	conn := core.RepositoryConnection{
		WorkspaceID: currentWorkspace.ID,
		Identifier:  identifier,
		Provider:    core.RepositoryProviderGithubV1,
		Payload:     payload,
	}
	err = saveRepositoryConnection(r.Context(), &conn)
	if err == errConnectionAlreadyExists {
		apihelper.ValidationErrResp(w, map[string]string{
			"installation_id": "github connection with same organization/user already exists",
		})
		return
	}
	if err != nil {
		apihelper.InternalServerErrResp(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
