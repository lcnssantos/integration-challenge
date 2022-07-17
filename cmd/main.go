package main

import (
	"github.com/lcnssantos/integration-challenge/internal/app"
	"github.com/lcnssantos/integration-challenge/internal/infra/configuration"
	"github.com/lcnssantos/integration-challenge/internal/infra/httpclient"
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

	client := httpclient.NewHttpClient()

	serviceA := app.NewServiceAImpl(client, configuration.Environment.ServiceABaseUrl)
	serviceB := app.NewServiceBImpl(client, configuration.Environment.ServiceBBaseUrl)

	controller := httpserver.NewController(serviceA, serviceB)

	server := httpserver.NewServer(8080, controller)
	server.Listen()
}
