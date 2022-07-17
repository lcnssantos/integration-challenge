package app

import (
	"context"
	"fmt"

	"github.com/lcnssantos/integration-challenge/internal/domain"
	"github.com/lcnssantos/integration-challenge/internal/infra/httpclient"
)

type ServiceAImpl struct {
	httpClient httpclient.HttpClient
	baseUrl    string
}

type ServiceAResponse struct {
	Cotacao float64         `json:"cotacao"`
	Moeda   domain.Currency `json:"moeda"`
	Symbol  string          `json:"symbol"`
}

func NewServiceAImpl(httpClient httpclient.HttpClient, baseUrl string) ServiceAImpl {
	return ServiceAImpl{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

func (s ServiceAImpl) Query(ctx context.Context, currency domain.Currency) (*domain.Price, error) {
	url := fmt.Sprintf("%s/cotacao?moeda=%s", s.baseUrl, currency)

	var response ServiceAResponse

	err := s.httpClient.Get(ctx, url, &response)

	if err != nil {
		return nil, err
	}

	return &domain.Price{
		Value:    response.Cotacao,
		Currency: response.Moeda,
	}, nil
}
