package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/lcnssantos/integration-challenge/internal/app"
	"github.com/lcnssantos/integration-challenge/internal/domain"
	"github.com/lcnssantos/integration-challenge/internal/infra/concurrency"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	serviceA app.Strategy
	serviceB app.Strategy
}

func NewController(serviceA app.Strategy, serviceB app.Strategy) Controller {
	return Controller{
		serviceA: serviceA,
		serviceB: serviceB,
	}
}

func (c Controller) Query(ctx *gin.Context) {
	currency := domain.Currency(ctx.Param("currency"))

	prices := []domain.Price{}

	tasks := concurrency.ExecuteConcurrentTasks(concurrency.TaskInput{
		Task: func() (interface{}, error) {
			return c.serviceA.Query(ctx, currency)
		},
		Tag: "query-service-a",
	}, concurrency.TaskInput{
		Task: func() (interface{}, error) {
			return c.serviceB.Query(ctx, currency)
		},
		Tag: "query-service-b",
	})

	priceATask := tasks[0]
	priceBTask := tasks[1]

	if priceATask.Err == nil {
		prices = append(prices, *priceATask.Result.(*domain.Price))
	} else {
		log.Error().Err(priceATask.Err).Msg("error querying service A")
	}

	if priceBTask.Err == nil {
		prices = append(prices, *priceBTask.Result.(*domain.Price))
	} else {
		log.Error().Err(priceBTask.Err).Msg("error querying service B")
	}

	if len(prices) == 0 {
		ctx.JSON(404, gin.H{
			"error": "not found",
		})
		return
	}

	bestPrice := app.GetBestOffer(prices)

	ctx.JSON(200, bestPrice)
}
