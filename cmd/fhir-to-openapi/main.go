package main

import (
	"log"
	"os"

	"github.com/gotidy/json-schema-to-openapi/pkg/generator"
)

func main() {
	config := ParseFlags()

	input := os.Stdin
	if config.Input != "" {
		var err error
		input, err = os.Open(config.Input)
		if err != nil {
			log.Fatalf("Opening file «%s»: %s", config.Input, err)
		}
		defer input.Close()
	}
	output := os.Stdout
	if config.Output != "" {
		var err error
		output, err = os.Create(config.Output)
		if err != nil {
			log.Fatalf("Opening file «%s»: %s", config.Output, err)
		}
		defer output.Close()
	}

	if err := generator.New().Do(input, output); err != nil {
		log.Fatalf("Generation OpenAPI: %s", err)
	}
}
