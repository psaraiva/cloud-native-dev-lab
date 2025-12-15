// @title Cloud Native Go 1.19 Lab
// @version 1.0
// @description This lab aims to create a secure development environment for legacy Go applications, in this case version 1.19, by enforcing architectural standards (Ports & Adapters) and consistency (Dev Containers).
// @termsOfService http://swagger.io/terms/

// @contact.name Pedro Saraiva
// @contact.url https://github.com/psaraiva/cloud-native-dev-lab/issues

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /

package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"cloud-native-dev-lab/client"
	"cloud-native-dev-lab/config"
	"cloud-native-dev-lab/docs"
	"cloud-native-dev-lab/service"

	echoSwagger "github.com/swaggo/echo-swagger"
)

var svc service.Service

type AnswerModel struct {
	PokemonAnswer string `json:"pokemon_answer"`
	PokemonID     string `json:"pokemon_id"`
}

// @Summary Welcome Message
// @Tags Status
// @Produce plain
// @Success 200 {string} string "greetings"
// @Router / [get]
func handlerWelcome(c echo.Context) error {
	message := svc.WelcomeMessage()
	return c.String(http.StatusOK, message)
}

// @Summary Pokemon Details
// @Description Fetches a Pokemon by name via the external API and returns processed information through the service.
// @Tags Data
// @Param name path string true "Pokemon Name (e.g., ditto)"
// @Produce json
// @Success 200 {object} main.AnswerModel
// @Router /pokemon/{name} [get]
func handlerAnswer(c echo.Context) error {
	ctx := c.Request().Context()
	pokemonName := c.Param("name")
	data := svc.IdentifyPokemon(ctx, pokemonName)
	return c.JSON(http.StatusOK, data)
}

func main() {
	cfg := config.LoadConfig()

	// The swag version v1.8.12 not have suport -property HostPort, wellcome to old picture.
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Port)

	pokeFetcher := client.NewPokemonAPIClient(cfg.PokeAPIBaseURL)
	svc = service.NewPokemonService(pokeFetcher)

	e := echo.New()
	e.Use(middleware.Logger())
	if cfg.DebugMode {
		e.Debug = true
		fmt.Println("Debug Mode: Actived")
	}

	e.GET("/", handlerWelcome)
	e.GET("/pokemon/:name", handlerAnswer)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
