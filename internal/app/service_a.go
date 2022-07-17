package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lcnssantos/integration-challenge/internal/domain"
	"github.com/lcnssantos/integration-challenge/internal/infra/httpclient"
	"github.com/rs/zerolog/log"
)

type serviceACache struct {
	sync.Mutex
	data       *ServiceAResponse
	expiration time.Time
}

type ServiceAImpl struct {
	httpClient httpclient.HttpClient
	baseUrl    string
	cache      serviceACache
}

type ServiceAResponse struct {
	Cotacao float64         `json:"cotacao"`
	Moeda   domain.Currency `json:"moeda"`
	Symbol  string          `json:"symbol"`
}

func NewServiceAImpl(httpClient httpclient.HttpClient, baseUrl string) Strategy {
	return &ServiceAImpl{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

func (s *ServiceAImpl) GetTag() string {
	return "service-a"
}

func (s *ServiceAImpl) Query(ctx context.Context, currency domain.Currency) (*domain.Price, error) {
	url := fmt.Sprintf("%s/cotacao?moeda=%s", s.baseUrl, currency)

	var response ServiceAResponse

	if s.cache.data != nil && s.cache.expiration.After(time.Now()) {
		response = *s.cache.data
		log.Debug().Str("service", "service-a").Msg("using cached response")
	} else {
		s.cache.Mutex.Lock()
		defer s.cache.Mutex.Unlock()
		log.Debug().Str("service", "service-a").Msg("querying service")
		err := s.httpClient.Get(ctx, url, &response)

		if err != nil {
			log.Error().Err(err).Str("service", "service-a").Msg("error querying service")
			return nil, err
		}

		s.cache.data = &response
		s.cache.expiration = time.Now().Add(CACHE_TIME)
	}

	return &domain.Price{
		Value:    response.Cotacao,
		Currency: response.Moeda,
	}, nil
}
