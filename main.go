package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokemonReference struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Pokemon struct {
	Id int           `json:"id"`
	Name string      `json:"name"`
	Height int       `json:"height"`
	Weight int       `json:"weight"`
}

type GetListOfPokemonResult struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []PokemonReference `json:"results"`
}

func (p PokemonReference) String() string {
	return fmt.Sprintf("PokemonReference(Name: %s, Url: %s)", p.Name, p.Url)
}

func (p Pokemon) String() string {
	return fmt.Sprintf("Pokemon(Id: %d, Name: %s, Height:%d, Weight:%d)", p.Id, p.Name, p.Height, p.Weight)
}

func GetPokemon(url string) (Pokemon, error) {
	
	resp, err := http.Get(url)

	if err != nil {
		return Pokemon{}, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	result := Pokemon{}

	parse_error := json.Unmarshal(body, &result)

	if parse_error != nil {
		fmt.Println("failed to parse", parse_error)
	}

	return result, nil
}

func main() {
	// TODO: environment variable?
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon")

	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	result := GetListOfPokemonResult{}

	parse_error := json.Unmarshal(body, &result)

	if parse_error != nil {
		fmt.Println("failed to parse", parse_error)
	}

	fmt.Println("Pokemon count", result.Count)
	fmt.Println("Next link", result.Next)

	for _, pokemon := range result.Results {
		detail, err := GetPokemon(pokemon.Url)
		if err != nil {
			fmt.Println("Error fetching", pokemon.Url)
		}

		fmt.Println(detail)
	}
}
