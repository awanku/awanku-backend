package coreapi

import (
	"net/http"

	"github.com/awanku/awanku/internal/coreapi/appctx"
	"github.com/awanku/awanku/internal/coreapi/auth"
	"github.com/awanku/awanku/internal/coreapi/user"
	"github.com/awanku/awanku/internal/coreapi/workspace"
	workspaceProject "github.com/awanku/awanku/internal/coreapi/workspace/project"
	workspaceProjectResource "github.com/awanku/awanku/internal/coreapi/workspace/project/resource"
	workspaceRepository "github.com/awanku/awanku/internal/coreapi/workspace/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func (s *Server) initRoutes() {
	s.router.Use(baseMiddleware)

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("See https://api.awanku.id/docs/ for API documentation"))
	})

	s.router.Get("/status", statusHandler(s.db))

	s.router.Route("/v1", func(r chi.Router) {
		r.Use(appctx.Middleware(appctx.Config{
			Environment:     s.Config.Environment,
			DB:              s.db,
			GithubAppConfig: s.githubAppConfig,
		}))

		r.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPatch, http.MethodDelete},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           5 * 60,
		}).Handler)

		r.Route("/auth", func(r chi.Router) {
			r.Get("/{provider:[a-z]+}/connect", auth.HandleOauthProviderConnect)
			r.Get("/{provider:[a-z]+}/callback", auth.HandleOauthProviderCallback)
			r.Post("/token", auth.HandleExchangeOauthToken(s.oauthTokenSecretKey))
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(auth.OauthTokenValidatorMiddleware(s.oauthTokenSecretKey))

			r.Get("/me", user.HandleGetMe)
		})

		r.Route("/workspaces", func(r chi.Router) {
			r.Use(auth.OauthTokenValidatorMiddleware(s.oauthTokenSecretKey))

			r.Get("/", workspace.HandleListAll)

			r.Route("/{workspace_id:[0-9]+}", func(r chi.Router) {
				r.Use(workspace.CurrentWorkspaceMiddleware)

				r.Route("/repositories", func(r chi.Router) {
					r.Get("/", workspaceRepository.HandleListAllRepositories)
					r.Get("/connections", workspaceRepository.HandleListAllConnections)

					r.Route("/providers", func(r chi.Router) {
						r.Get("/github", workspaceRepository.HandleConnectGithub)
						r.Post("/github", workspaceRepository.HandleSaveGithubConnection)
					})
				})

				r.Route("/projects", func(r chi.Router) {
					r.Get("/", workspaceProject.HandleListAll)

					r.Route("/{project_id:[0-9]+}", func(r chi.Router) {
						r.Route("/resources", func(r chi.Router) {
							r.Get("/", workspaceProjectResource.HandleListAll)
							r.Post("/", workspaceProjectResource.HandleCreate)

							r.Route("/{resource_id:[0-9]+}", func(r chi.Router) {
								r.Get("/", workspaceProjectResource.HandleGet)
								r.Patch("/", workspaceProjectResource.HandleUpdate)
								r.Delete("/", workspaceProjectResource.HandleDelete)
							})
						})
					})
				})
			})
		})
	})
}

func baseMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0, private, must-revalidate")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}
