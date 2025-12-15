package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"cloud-native-dev-lab/dto"
)

var ErrNotFound = errors.New("resource not found")
var ErrClientTimeout = errors.New("client timeout or canceled")
var ErrBadResponse = errors.New("bad response from provider")

type PokemonAPIClient struct {
	client  *http.Client
	baseURL string
}

func NewPokemonAPIClient(baseURL string) *PokemonAPIClient {
	return &PokemonAPIClient{
		client:  &http.Client{Timeout: time.Second * 5},
		baseURL: baseURL,
	}
}

func (c *PokemonAPIClient) WithClient(httpClient *http.Client) *PokemonAPIClient {
	c.client = httpClient
	return c
}

func (c *PokemonAPIClient) GetByName(ctx context.Context, name string) (*dto.Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request build failed: %w: %v", ErrBadResponse, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("call external canceled or timeout: %w: %v", ErrClientTimeout, err)
		}
		return nil, fmt.Errorf("network error during call: %w: %v", ErrBadResponse, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("unexpected status %d from provider: %w", resp.StatusCode, ErrBadResponse)
	}

	var pokemon dto.Pokemon
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w: %v", ErrBadResponse, err)
	}

	return &pokemon, nil
}
