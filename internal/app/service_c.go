package app

import (
	"context"
	"fmt"

	"github.com/lcnssantos/integration-challenge/internal/domain"
	"github.com/lcnssantos/integration-challenge/internal/infra/concurrency"
	"github.com/lcnssantos/integration-challenge/internal/infra/httpclient"
	"github.com/rs/zerolog/log"
)

type ServiceCImpl struct {
	httpClient httpclient.HttpClient
	baseUrl    string
	myBaseUrl  string
	pubSub     *concurrency.PubSub[WebhookResponse]
}

type WebhookResponse struct {
	CorrelationID string          `json:"cid"`
	F             float64         `json:"f"`
	T             domain.Currency `json:"t"`
	V             float64         `json:"v"`
}

type ServiceCResponse struct {
	CorrelationID string `json:"cid"`
}

type ServiceCRequest struct {
	Type     domain.Currency `json:"tipo"`
	Callback string          `json:"callback"`
}

func NewServiceCImpl(httpClient httpclient.HttpClient, baseUrl string, myBaseUrl string, pubSub *concurrency.PubSub[WebhookResponse]) Strategy {
	return &ServiceCImpl{
		httpClient: httpClient,
		baseUrl:    baseUrl,
		myBaseUrl:  myBaseUrl,
		pubSub:     pubSub,
	}
}

func (s *ServiceCImpl) GetTag() string {
	return "service-c"
}

func (s *ServiceCImpl) Query(ctx context.Context, currency domain.Currency) (*domain.Price, error) {
	url := fmt.Sprintf("%s/cotacao", s.baseUrl)

	var msg WebhookResponse

	log.Debug().Str("service", "service-c").Msg("querying service")

	var response ServiceCResponse

	err := s.httpClient.Post(ctx, url, ServiceCRequest{
		Type:     currency,
		Callback: s.myBaseUrl + "/service-c/callback",
	}, &response)

	if err != nil {
		log.Error().Err(err).Str("service", "service-c").Msg("error querying service")
		return nil, err
	}

	ch := s.pubSub.Subscribe(response.CorrelationID)

	msg = <-ch

	s.pubSub.Unsubscribe(response.CorrelationID)

	log.Debug().Interface("msg", msg).Str("service", "service-c").Msg("received message")

	return &domain.Price{
		Value:    msg.V / msg.F,
		Currency: domain.Currency(msg.T),
	}, nil

}
