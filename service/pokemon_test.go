package service_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"cloud-native-dev-lab/client"
	"cloud-native-dev-lab/dto"
	"cloud-native-dev-lab/mocks"
	"cloud-native-dev-lab/service"

	"github.com/stretchr/testify/mock"
)

type testCase struct {
	name              string
	inputName         string
	mockReturnPokemon *dto.Pokemon
	mockReturnError   error
	expectedStatus    string
	expectedContains  string
}

func mockSetup(m *mocks.PokemonProvider, inputName string, retPokemon *dto.Pokemon, retError error) {
	m.On("GetByName", mock.Anything, inputName).
		Return(retPokemon, retError).
		Once()
}

func TestIdentifyPokemon(t *testing.T) {
	tests := []testCase{
		{
			name:      "IdentifyPokemon_ReturnsOkStatus",
			inputName: "snorlax",
			mockReturnPokemon: &dto.Pokemon{
				ID:   143,
				Name: "snorlax",
				Abilities: []dto.PokemonAbility{
					{
						Ability: dto.AbilityDetails{Name: "thick-fat", URL: "placeholder"},
					},
				},
			},
			mockReturnError:  nil,
			expectedStatus:   "ok",
			expectedContains: "thick-fat",
		},
		{
			name:              "IdentifyPokemon_TimeoutReturnsFailStatus",
			inputName:         "timeout-mon",
			mockReturnPokemon: nil,
			mockReturnError:   errors.New("context deadline exceeded"),
			expectedStatus:    "fail",
			expectedContains:  "Failed to retrieve Pokémon data from provider.",
		},
		{
			name:              "IdentifyPokemon_ValidNameNilAbilities",
			inputName:         "magikarp",
			mockReturnPokemon: &dto.Pokemon{ID: 129, Name: "magikarp", Abilities: nil},
			mockReturnError:   nil,
			expectedStatus:    "ok",
			expectedContains:  "None",
		},
		{
			name:              "IdentifyPokemon_NotFoundReturnsNotFoundStatus",
			inputName:         "missing-mon",
			mockReturnPokemon: nil,
			mockReturnError:   client.ErrNotFound,
			expectedStatus:    "not_found",
			expectedContains:  "The Pokémon not found.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockFetcher := mocks.NewPokemonProvider(t)
			mockSetup(mockFetcher, tc.inputName, tc.mockReturnPokemon, tc.mockReturnError)

			ctx := context.Background()
			if strings.Contains(tc.name, "Timeout") {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*100)
				defer cancel()
			}

			svc := service.NewPokemonService(mockFetcher)
			result := svc.IdentifyPokemon(ctx, tc.inputName)

			if result["status"] != tc.expectedStatus {
				t.Fatalf("Expected status %s, but got %s. Result Details: %v", tc.expectedStatus, result["status"], result)
			}

			if tc.expectedStatus == "ok" {
				answer := result["pokemon_answer"]
				if !strings.Contains(answer, tc.expectedContains) {
					t.Errorf("Expected message '%s', but got %s", tc.expectedContains, answer)
				}
			} else if tc.expectedStatus == "fail" {
				errMsg := result["error"]
				if !strings.Contains(errMsg, tc.expectedContains) {
					t.Errorf("Expected message '%s', but got %s", tc.expectedContains, errMsg)
				}
			}

			mockFetcher.AssertExpectations(t)
		})
	}
}

func TestWelcomeMessage_Table(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "Should return the default welcome message",
			expected: "Welcome to this laboratory!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewPokemonService(nil)
			actual := svc.WelcomeMessage()

			if actual != tt.expected {
				t.Errorf("Expected message '%s', but got '%s'", tt.expected, actual)
			}
		})
	}
}
