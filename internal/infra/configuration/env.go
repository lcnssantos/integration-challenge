package configuration

import (
	"os"

	"github.com/go-playground/validator/v10"
)

var Environment Env

type Env struct {
	ServiceABaseUrl string `validate:"required"`
	ServiceBBaseUrl string `validate:"required"`
	ServiceCBaseUrl string `validate:"required"`
}

func (e *Env) Validate() error {
	err := validator.New().Struct(*e)
	return err
}

func LoadEnv() {
	Environment = Env{
		ServiceABaseUrl: os.Getenv("SERVICE_A_BASE_URL"),
		ServiceBBaseUrl: os.Getenv("SERVICE_B_BASE_URL"),
		ServiceCBaseUrl: os.Getenv("SERVICE_C_BASE_URL"),
	}
}
