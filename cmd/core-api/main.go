package main

import (
	"log"

	"github.com/awanku/awanku/internal/coreapi"
)

// @title Awanku API
// @version 0.1

// @contact.name Awanku Support
// @contact.email hello@awanku.id

// @host api.awanku.id
// @schemes	https

// @securityDefinitions.apikey oauthAccessToken
// @in header
// @name Authorization
// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://api.awanku.id/v1/auth/token
// @authorizationUrl https://api.awanku.id/v1/auth/{provider}/connect

func main() {
	log.Println("Starting server...")

	conf := &coreapi.Config{}
	if err := conf.Load(); err != nil {
		log.Panicln("failed to parse config from environment variable:", err)
	}

	srv := coreapi.Server{
		Config: conf,
	}
	srv.Init()
	srv.Start()
}
