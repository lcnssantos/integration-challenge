package main

import (
	"github.com/lcnssantos/integration-challenge/internal/app"
	"github.com/lcnssantos/integration-challenge/internal/infra/concurrency"
	"github.com/lcnssantos/integration-challenge/internal/infra/configuration"
	"github.com/lcnssantos/integration-challenge/internal/infra/httpclient"
	"github.com/lcnssantos/integration-challenge/internal/infra/httpserver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	configuration.LoadEnv()

	if err := configuration.Environment.Validate(); err != nil {
		log.Panic().Err(err).Msg("invalid environment configuration")
	}

	client := httpclient.NewHttpClient()

	pubSub := concurrency.NewPubSub[app.WebhookResponse]()

	serviceA := app.NewServiceAImpl(client, configuration.Environment.ServiceABaseUrl)
	serviceB := app.NewServiceBImpl(client, configuration.Environment.ServiceBBaseUrl)

	serviceC := app.NewServiceCImpl(
		client,
		configuration.Environment.ServiceCBaseUrl,
		configuration.Environment.MyBaseUrl,
		pubSub,
	)

	controller := httpserver.NewController([]app.Strategy{serviceA, serviceB, serviceC}, pubSub)

	server := httpserver.NewServer(8080, controller)
	server.Listen()
}
