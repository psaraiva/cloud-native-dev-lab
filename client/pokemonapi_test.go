package client_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"cloud-native-dev-lab/client"
	"cloud-native-dev-lab/dto"
)

type MockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m *MockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Response, nil
}

func createMockClient(resp *http.Response, err error) *client.PokemonAPIClient {
	return client.NewPokemonAPIClient("http://fake-url").WithClient(&http.Client{
		Transport: &MockRoundTripper{Response: resp, Err: err},
		Timeout:   time.Second * 5,
	})
}

func TestGetByName(t *testing.T) {
	successResponse := `{"id": 143, "name": "snorlax", "abilities": [{"ability": {"name": "thick-fat", "url": ""}}]}`

	testCases := []struct {
		name              string
		inputName         string
		mockResp          *http.Response
		mockErr           error
		ctx               context.Context
		expectedPokemon   *dto.Pokemon
		expectedErrorType error
	}{
		{
			name:      "GetByName_Success_ValidResponse",
			inputName: "snorlax",
			mockResp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(successResponse)),
				Header:     make(http.Header),
			},
			mockErr: nil,
			ctx:     context.Background(),
			expectedPokemon: &dto.Pokemon{
				ID:   143,
				Name: "snorlax",
				Abilities: []dto.PokemonAbility{
					{Ability: dto.AbilityDetails{Name: "thick-fat"}},
				},
			},
			expectedErrorType: nil,
		},
		{
			name:              "GetByName_Failure_Timeout",
			inputName:         "testmon",
			mockResp:          nil,
			mockErr:           context.DeadlineExceeded,
			ctx:               context.Background(),
			expectedPokemon:   nil,
			expectedErrorType: client.ErrClientTimeout,
		},
		{
			name:              "GetByName_Failure_GenericNetworkError",
			inputName:         "testmon",
			mockResp:          nil,
			mockErr:           errors.New("connection reset by peer"),
			ctx:               context.Background(),
			expectedPokemon:   nil,
			expectedErrorType: client.ErrBadResponse,
		},
		{
			name:              "GetByName_Failure_NotFound404",
			inputName:         "testmon",
			mockResp:          &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader(""))},
			mockErr:           nil,
			ctx:               context.Background(),
			expectedPokemon:   nil,
			expectedErrorType: client.ErrNotFound,
		},
		{
			name:              "GetByName_Failure_InternalAPIError500",
			inputName:         "testmon",
			mockResp:          &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))},
			mockErr:           nil,
			ctx:               context.Background(),
			expectedPokemon:   nil,
			expectedErrorType: client.ErrBadResponse,
		},
		{
			name:              "GetByName_Failure_DecodeError",
			inputName:         "testmon",
			mockResp:          &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("isso não é JSON válido"))},
			mockErr:           nil,
			ctx:               context.Background(),
			expectedPokemon:   nil,
			expectedErrorType: client.ErrBadResponse,
		},
		{
			name:              "GetByName_Failure_InvalidInputRequest",
			inputName:         "snorlax\n",
			mockResp:          nil,
			mockErr:           nil,
			ctx:               context.Background(),
			expectedPokemon:   nil,
			expectedErrorType: client.ErrBadResponse,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := createMockClient(tc.mockResp, tc.mockErr)
			pokemon, err := c.GetByName(tc.ctx, tc.inputName)

			if tc.expectedErrorType == nil {
				if err != nil {
					t.Fatalf("Expected success, but received error: %v", err)
				}
				if pokemon.Name != tc.expectedPokemon.Name {
					t.Errorf("Expected name '%s', but got '%s'", tc.expectedPokemon.Name, pokemon.Name)
				}
				return
			}

			if err == nil {
				t.Fatal("Expected an error, but received success")
			}

			if !errors.Is(err, tc.expectedErrorType) {
				t.Errorf("Incorrect typed error. Expected error type '%v', but got: '%v'", tc.expectedErrorType, err)
			}
		})
	}
}
