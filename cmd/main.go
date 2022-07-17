package main

import (
	"github.com/lcnssantos/integration-challenge/internal/infra/configuration"
	"github.com/lcnssantos/integration-challenge/internal/infra/httpserver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	configuration.LoadEnv()

	if err := configuration.Environment.Validate(); err != nil {
		log.Panic().Err(err).Msg("invalid environment configuration")
	}

	server := httpserver.NewServer(8080, httpserver.Controller{})
	server.Listen()
}
