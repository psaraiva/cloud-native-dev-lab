package config

import (
	"fmt"
	"os"
)

type AppConfig struct {
	Port           string
	PokeAPIBaseURL string
	DebugMode      bool
}

func LoadConfig() AppConfig {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Fprintf(os.Stderr, "error critical: Env variable PORT undefined.\n")
		os.Exit(1)
	}

	pokeAPIBaseURL := os.Getenv("POKEMONAPI_BASE_URL")
	if pokeAPIBaseURL == "" {
		fmt.Fprintf(os.Stderr, "error critical: Env variable POKEMONAPI_BASE_URL undefined.\n")
		os.Exit(1)
	}

	debugMode := os.Getenv("ECHO_DEBUG") == "true"
	return AppConfig{
		Port:           port,
		PokeAPIBaseURL: pokeAPIBaseURL,
		DebugMode:      debugMode,
	}
}
