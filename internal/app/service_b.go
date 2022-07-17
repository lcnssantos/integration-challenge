package app

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/lcnssantos/integration-challenge/internal/domain"
	"github.com/lcnssantos/integration-challenge/internal/infra/httpclient"
	"github.com/rs/zerolog/log"
)

type serviceBCache struct {
	sync.Mutex
	data       *ServiceBResponse
	expiration time.Time
}
type ServiceBImpl struct {
	httpClient httpclient.HttpClient
	baseUrl    string
	cache      serviceBCache
}

type ServiceBResponse struct {
	Cotacao struct {
		Fator    int             `json:"fator"`
		Currency domain.Currency `json:"currency"`
		Valor    string          `json:"valor"`
	} `json:"cotacao"`
}

func NewServiceBImpl(httpClient httpclient.HttpClient, baseUrl string) *ServiceBImpl {
	return &ServiceBImpl{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

func (s *ServiceBImpl) GetTag() string {
	return "service-b"
}

func (s *ServiceBImpl) Query(ctx context.Context, currency domain.Currency) (*domain.Price, error) {
	url := fmt.Sprintf("%s/cotacao?curr=%s", s.baseUrl, currency)

	var response ServiceBResponse

	if s.cache.data != nil && s.cache.expiration.After(time.Now()) {
		response = *s.cache.data
		log.Debug().Str("service", "service-b").Msg("using cached response")
	} else {
		s.cache.Mutex.Lock()
		defer s.cache.Mutex.Unlock()

		log.Debug().Str("service", "service-b").Msg("querying service")
		err := s.httpClient.Get(ctx, url, &response)

		if err != nil {
			log.Error().Err(err).Str("service", "service-b").Msg("error querying service")
			return nil, err
		}

		s.cache.data = &response
		s.cache.expiration = time.Now().Add(CACHE_TIME)
	}

	value, err := strconv.Atoi(response.Cotacao.Valor)

	if err != nil {
		log.Error().Err(err).Msg("failed to convert value")
		return nil, err
	}

	return &domain.Price{
		Value:    float64(value) / float64(response.Cotacao.Fator),
		Currency: response.Cotacao.Currency,
	}, nil
}
