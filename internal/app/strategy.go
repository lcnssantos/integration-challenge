package app

import (
	"context"

	"github.com/lcnssantos/integration-challenge/internal/domain"
)

type Strategy interface {
	Query(ctx context.Context, currency domain.Currency) (*domain.Price, error)
}
