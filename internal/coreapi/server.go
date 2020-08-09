package coreapi

import (
	"net/http"
	"time"

	hansip "github.com/asasmoyo/pq-hansip"
	"github.com/awanku/awanku/pkg/core"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v9"
)

type Server struct {
	router              chi.Router
	db                  *hansip.Cluster
	oauthTokenSecretKey []byte
	githubAppConfig     *core.GithubAppConfig

	Config *Config
}

func (s *Server) Init() error {
	s.router = chi.NewRouter()
	s.router.Use(middleware.Logger)

	s.oauthTokenSecretKey = []byte(s.Config.OAuthSecretKey)

	db, err := initDB(s.Config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	s.db = db

	githubAppConfig, err := s.Config.GithubAppConfig()
	if err != nil {
		panic(err)
	}
	s.githubAppConfig = githubAppConfig

	s.initRoutes()
	return nil
}

func (s *Server) Start() error {
	return http.ListenAndServe("0.0.0.0:3000", s.router)
}

func initDB(dbURL string) (*hansip.Cluster, error) {
	opt, err := pg.ParseURL(dbURL)
	if err != nil {
		return nil, err
	}
	db := hansip.NewCluster(&hansip.Config{
		Primary:        hansip.WrapGoPG(pg.Connect(opt)),
		Replicas:       []hansip.SQL{hansip.WrapGoPG(pg.Connect(opt))},
		PingTimeout:    1 * time.Second,
		ConnCheckDelay: 5 * time.Second,
	})
	return db, nil
}
