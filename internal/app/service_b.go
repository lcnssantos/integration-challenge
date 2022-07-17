package app

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lcnssantos/integration-challenge/internal/domain"
	"github.com/lcnssantos/integration-challenge/internal/infra/httpclient"
	"github.com/rs/zerolog/log"
)

type ServiceBImpl struct {
	httpClient httpclient.HttpClient
	baseUrl    string
}

type ServiceBResponse struct {
	Cotacao struct {
		Fator    int             `json:"fator"`
		Currency domain.Currency `json:"currency"`
		Valor    string          `json:"valor"`
	} `json:"cotacao"`
}

func NewServiceBImpl(httpClient httpclient.HttpClient, baseUrl string) ServiceBImpl {
	return ServiceBImpl{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

func (s ServiceBImpl) Query(ctx context.Context, currency domain.Currency) (*domain.Price, error) {
	url := fmt.Sprintf("%s/cotacao?curr=%s", s.baseUrl, currency)

	var response ServiceBResponse

	err := s.httpClient.Get(ctx, url, &response)

	if err != nil {
		return nil, err
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
