package httpserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lcnssantos/integration-challenge/internal/app"
	"github.com/lcnssantos/integration-challenge/internal/domain"
	"github.com/lcnssantos/integration-challenge/internal/infra/cache"
	"github.com/lcnssantos/integration-challenge/internal/infra/concurrency"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	strategies []app.Strategy
	pubSub     *concurrency.PubSub[app.WebhookResponse]
	cache      *cache.CacheProxy
}

func NewController(strategies []app.Strategy, pubSub *concurrency.PubSub[app.WebhookResponse], cache *cache.CacheProxy) Controller {
	return Controller{
		strategies: strategies,
		pubSub:     pubSub,
		cache:      cache,
	}
}

func (c *Controller) Subscribe(ctx *gin.Context) {
	var webhookResponse app.WebhookResponse
	if err := ctx.BindJSON(&webhookResponse); err != nil {
		log.Error().Err(err).Msg("error binding json")
		ctx.JSON(400, gin.H{
			"error": "error binding json",
		})
		return
	}

	c.pubSub.Publish(webhookResponse.CorrelationID, webhookResponse)
	ctx.JSON(200, gin.H{
		"message": "ok",
	})
}

func (c *Controller) Query(ctx *gin.Context) {
	currency := domain.Currency(ctx.Param("currency"))

	prices := []domain.Price{}

	tasksInputs := make([]concurrency.TaskInput, len(c.strategies))

	for i := 0; i < len(c.strategies); i++ {
		strategy := c.strategies[i]

		tasksInputs[i] = concurrency.TaskInput{
			Task: func() (interface{}, error) {
				return c.cache.Proxy(func() (interface{}, error) {
					return strategy.Query(ctx, currency)
				}, fmt.Sprintf("%s-query-%s", strategy.GetTag(), currency))
			},
			Tag: strategy.GetTag(),
		}

	}

	tasks := concurrency.ExecuteConcurrentTasks(tasksInputs)

	for _, task := range tasks {
		if task.Err != nil {
			log.Error().Err(task.Err).Msg("error executing task")
			continue
		}

		prices = append(prices, *task.Result.(*domain.Price))
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
