package dto

type AbilityDetails struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonAbility struct {
	Ability AbilityDetails `json:"ability"`
}

type Pokemon struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	Abilities []PokemonAbility `json:"abilities"`
}
