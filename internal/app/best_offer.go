package app

import "github.com/lcnssantos/integration-challenge/internal/domain"

func GetBestOffer(prices []domain.Price) domain.Price {
	lower := prices[0]

	for _, price := range prices {
		if price.Value < lower.Value {
			lower = price
		}
	}

	return lower
}
