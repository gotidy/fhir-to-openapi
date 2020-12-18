package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gotidy/fhir-to-openapi/pkg/generator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	config := ParseFlags()

	input := os.Stdin
	if config.Input != "" {
		var err error
		input, err = os.Open(config.Input)
		if err != nil {
			log.Fatal().Msgf("Opening file «%s»: %s", config.Input, err)
		}
		defer input.Close()
	}
	output := os.Stdout
	if config.Output != "" {
		var err error
		output, err = os.Create(config.Output)
		if err != nil {
			log.Fatal().Msgf("Opening file «%s»: %s", config.Output, err)
		}
		defer output.Close()
	}

	format := generator.JSON
	if ext := strings.ToLower(filepath.Ext(config.Output)); ext == ".yaml" || ext == ".yml" {
		format = generator.YAML
	}

	if err := generator.New().Do(input, output, format); err != nil {
		log.Fatal().Msgf("Generation OpenAPI: %s", err)
	}

	log.Info().Msg("Successfully generated.")
}
