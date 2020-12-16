package generator

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
)

type Generator struct {
	Swagger *openapi3.Swagger
}

const baseOpenAPIData = `
{
	"openapi": "3.0.3",
	"info": {
	  "title": "Aidbox client",
	  "version": "1.0"
	},
	"servers": [
	  {
		"url": "https://test.aidbox.app/"
	  }
	],
	"paths": {
	  "/__healthcheck": {
		"get": {
		  "summary": "Checks if the server is running",
		  "security": [],
		  "responses": {
			"200": {
			  "description": "OK"
			}
		  }
		}
	  }
	},
	"components": {
	  "securitySchemes": {
		"BasicAuth": {
		  "type": "http",
		  "scheme": "basic"
		}
	  }
	},
	"security": [
	  {
		"BasicAuth": [
		  "base"
		]
	  }
	]
  }
`

func New() *Generator {
	s, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData([]byte(baseOpenAPIData))
	if err != nil {
		panic(fmt.Errorf("loading base openapi data: %w", err))
	}
	return &Generator{Swagger: s}
}

func (g *Generator) String() string {
	return fmt.Sprintf("%+v", g.Swagger)
}

func (g *Generator) Do(jsonSchema io.Reader, oapiSchema io.Writer) error {
	json.NewEncoder(oapiSchema).Encode(g.Swagger)
	return nil
}
