package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud-native-dev-lab/client"
	"cloud-native-dev-lab/dto"
)

type Service interface {
	WelcomeMessage() string
	IdentifyPokemon(ctx context.Context, name string) map[string]string
}

type PokemonProvider interface {
	GetByName(ctx context.Context, name string) (*dto.Pokemon, error)
}

type PokemonService struct {
	PokemonProvider PokemonProvider
}

func NewPokemonService(fetcher PokemonProvider) *PokemonService {
	return &PokemonService{
		PokemonProvider: fetcher,
	}
}

func (s *PokemonService) WelcomeMessage() string {
	return "Welcome to this laboratory!"
}

func (s *PokemonService) IdentifyPokemon(ctx context.Context, name string) map[string]string {
	pokemon, err := s.PokemonProvider.GetByName(ctx, name)

	if err != nil {
		if errors.Is(err, client.ErrNotFound) {
			return map[string]string{
				"error":  fmt.Sprintf("Pokémon '%s' not found.", name),
				"status": "not_found",
			}
		}

		return map[string]string{
			"error":  "Failed to retrieve Pokémon data from provider.",
			"status": "fail",
		}
	}

	now := time.Now().Format(time.RFC3339)
	ability := "None"
	if len(pokemon.Abilities) > 0 {
		ability = pokemon.Abilities[0].Ability.Name
	}

	answer := fmt.Sprintf(
		"Pokémon %s (ID: %d) has the main ability '%s'. Generated in: %s",
		pokemon.Name,
		pokemon.ID,
		ability,
		now,
	)

	return map[string]string{
		"pokemon_id":     fmt.Sprintf("%d", pokemon.ID),
		"pokemon_answer": answer,
		"status":         "ok",
	}
}
