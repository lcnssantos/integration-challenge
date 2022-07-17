package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type HttpClient struct {
	httpClient *http.Client
}

func NewHttpClient() HttpClient {
	return HttpClient{
		httpClient: &http.Client{Timeout: time.Duration(20) * time.Second},
	}
}

func (h *HttpClient) Post(ctx context.Context, url string, body interface{}, output interface{}) error {
	bodyEncoded, err := json.Marshal(&body)

	if err != nil {
		log.Error().Err(err).Msg("failed to encode body")
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyEncoded))

	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return err
	}

	startTime := time.Now()

	res, err := h.httpClient.Do(req)

	if err != nil {
		log.Error().Err(err).Msg("failed to do request")
		return err
	}

	difference := time.Since(startTime)

	log.Debug().
		Str("url", url).
		Str("method", http.MethodPost).
		Dur("duration", difference).
		Int("status", res.StatusCode).
		Msgf("HTTP Request | %s", url)

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("failed to post to %s, status code %d", url, res.StatusCode))
	}

	err = json.NewDecoder(res.Body).Decode(&output)

	if err != nil {
		log.Error().Err(err).Msg("failed to decode response")
		return err
	}

	return nil
}

func (h *HttpClient) Get(ctx context.Context, url string, output interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return err
	}

	startTime := time.Now()

	res, err := h.httpClient.Do(req)

	if err != nil {
		log.Error().Err(err).Msg("failed to do request")
		return err
	}

	difference := time.Since(startTime)

	log.Debug().
		Str("url", url).
		Str("method", http.MethodGet).
		Dur("duration", difference).
		Int("status", res.StatusCode).
		Msgf("HTTP Request | %s", url)

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintf("http status code: %d | URL: %s", res.StatusCode, url))
	}

	err = json.NewDecoder(res.Body).Decode(&output)

	if err != nil {
		log.Error().Err(err).Msg("failed to decode response")
		return err
	}

	return nil
}
