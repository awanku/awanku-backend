package appctx

import (
	"context"
	"net/http"

	hansip "github.com/asasmoyo/pq-hansip"
	"github.com/awanku/awanku/pkg/core"
)

type Config struct {
	Environment     string
	DB              *hansip.Cluster
	GithubAppConfig *core.GithubAppConfig
}

// Middleware inject stuff into request context
func Middleware(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), KeyEnvironment, config.Environment)
			ctx = context.WithValue(ctx, KeyDatabase, config.DB)
			ctx = context.WithValue(ctx, KeyGithubAppConfig, config.GithubAppConfig)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
