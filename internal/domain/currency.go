package domain

type Currency string

type Price struct {
	Value    float64  `json:"value"`
	Currency Currency `json:"currency"`
}
