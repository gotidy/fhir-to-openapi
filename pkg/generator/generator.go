package generator

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/gotidy/ptr"
)

type Format int

const (
	JSON Format = iota
	YAML
)

type Generator struct {
	SkipUnderscore bool
	Swagger        *openapi3.Swagger
	Schema         *Schema
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
	return &Generator{Swagger: s, Schema: &Schema{}, SkipUnderscore: true}
}

func (g *Generator) String() string {
	return fmt.Sprintf("%+v", g.Swagger)
}

func (g *Generator) Do(schema io.Reader, output io.Writer, format Format) error {
	if err := g.encodeSchema(schema); err != nil {
		return err
	}

	g.Swagger.Components.Schemas = g.convertNamedSchemas(g.Schema.Definitions, true)

	b, err := json.MarshalIndent(g.Swagger, "", "    ")
	if err != nil {
		return err
	}
	if format == YAML {
		b, err = yaml.JSONToYAML(b)
		if err != nil {
			return err
		}
	}
	_, err = output.Write(b)
	return err
}

func (g *Generator) initSwagger() {
	// Responses
	g.Swagger.Components.Responses = openapi3.Responses{
		"Error": &openapi3.ResponseRef{Value: &openapi3.Response{
			Description: ptr.String("Error"),
			Content:     openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/OperationOutcome", nil)),
		}},
		"Bundle": &openapi3.ResponseRef{Value: &openapi3.Response{
			Description: ptr.String("OK"),
			Content:     openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/Bundle", nil)),
		}},
	}

	// Parameters
	g.Swagger.Components.Parameters = openapi3.ParametersMap{
		"search": &openapi3.ParameterRef{Value: &openapi3.Parameter{
			Name:     "search",
			In:       "query",
			Required: true,
			Schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:                        "object",
					AdditionalPropertiesAllowed: ptr.Bool(true),
				},
			},
		}},
	}

	// Path /
	contentBundle := openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/Bundle", nil))
	respBundle := &openapi3.ResponseRef{Ref: "#/components/schemas/Bundle"}
	respErr := &openapi3.ResponseRef{Ref: "#/components/schemas/OperationOutcome"}
	responsesBundle := openapi3.Responses{
		"200": respBundle,
		"201": respBundle,
		"400": respErr,
		"403": respErr,
		"404": respErr,
		"405": respErr,
		"409": respErr,
		"422": respErr,
	}

	g.Swagger.Paths["/"] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{Ref: "#/components/parameters/search"},
			},
			Description: "This searches all resources of a particular type using the criteria represented in the parameters.",
			Tags:        []string{"search"},
			Responses: openapi3.Responses{
				"200": respBundle,
				"400": respErr,
				"401": respErr,
				"403": respErr,
				"404": respErr,
			},
		},
		Post: &openapi3.Operation{
			Description: "The create interaction creates a bundle of resources.",
			Tags:        []string{"create"},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  contentBundle,
				},
			},
			Responses: responsesBundle,
		},
		Put: &openapi3.Operation{
			Description: "The update interaction creates or updates a bundle of resources.",
			Tags:        []string{"create", "update"},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  contentBundle,
				},
			},
			Responses: responsesBundle,
		},
	}
}

func (g *Generator) encodeSchema(r io.Reader) error {
	return json.NewDecoder(r).Decode(g.Schema)
}

func (g *Generator) convertSchemas(src []*Type) []*openapi3.SchemaRef {
	dst := make([]*openapi3.SchemaRef, len(src))
	for i, schema := range src {
		dst[i] = g.convertSchema(schema)
	}
	return dst
}

func (g *Generator) convertNamedSchemas(src map[string]*Type, genOps bool) openapi3.Schemas {
	var dst openapi3.Schemas
	if len(src) > 0 {
		dst = make(openapi3.Schemas, len(src))
		for name, schema := range src {
			if strings.HasPrefix(name, "_") {
				if _, ok := src[strings.TrimPrefix(name, "_")]; ok && g.SkipUnderscore {
					continue
				}
			}
			if len(name) == 0 { // || (strings.HasPrefix(name, "_") && g.SkipUnderscore) {
				continue
			}
			dst[name] = g.convertSchema(schema)
			if genOps && unicode.IsUpper([]rune(name)[0]) {
				g.createPathes(name)
			}
		}
	}
	return dst
}

func (g *Generator) convertSchema(src *Type) *openapi3.SchemaRef {
	if src == nil {
		return nil
	}

	// Use standard OpenAPI types
	if src.Ref != "" {
		mapper.Convert(src)
	}

	dst := &openapi3.SchemaRef{
		Ref: strings.Replace(src.Ref, "#/definitions/", "#/components/schemas/", 1),
		Value: &openapi3.Schema{
			Type:         src.Type,
			Description:  src.Description,
			Default:      src.Default,
			Enum:         src.Enum,
			Title:        src.Title,
			Format:       src.Format,
			Pattern:      src.Pattern,
			MinLength:    src.MinLength,
			MaxLength:    src.MaxLength,
			Min:          src.Minimum,
			ExclusiveMin: src.ExclusiveMinimum,
			Max:          src.Maximum,
			ExclusiveMax: src.ExclusiveMaximum,
			MultipleOf:   src.MultipleOf,
			MinItems:     src.MinItems,
			MaxItems:     src.MaxItems,
			UniqueItems:  src.UniqueItems,
			MinProps:     src.MinProperties,
			MaxProps:     src.MaxProperties,
			Items:        g.convertSchema(src.Items),
			OneOf:        g.convertSchemas(src.OneOf),
			AllOf:        g.convertSchemas(src.AllOf),
			AnyOf:        g.convertSchemas(src.AnyOf),
			Properties:   g.convertNamedSchemas(src.Properties, false),
			Required:     src.Required,
		},
	}

	if src.Const != nil {
		dst.Value.Enum = append(dst.Value.Enum, src.Const)
	}

	if len(src.Examples) > 0 {
		dst.Value.Example = src.Examples[0]
	}

	return dst
}

func (g *Generator) createPathes(entity string) {
	content := openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/"+entity, nil))
	respEntity := &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr.String("OK"),
		Content:     openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/"+entity, nil)),
	}}
	respErr := &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr.String("Error"),
		Content:     openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/OperationOutcome", nil)),
	}}
	respBundle := &openapi3.ResponseRef{Value: &openapi3.Response{
		Description: ptr.String("OK"),
		Content:     openapi3.NewContentWithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/Bundle", nil)),
	}}
	responsesEntity := openapi3.Responses{
		"200": respEntity,
		"201": respEntity,
		"400": respErr,
		"403": respErr,
		"404": respErr,
		"405": respErr,
		"409": respErr,
		"422": respErr,
	}
	// GET /<Entity>
	g.Swagger.Paths["/"+entity] = &openapi3.PathItem{
		Get: &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{Ref: "#/components/parameters/search"},
			},
			Description: "This searches all resources of a particular type using the criteria represented in the parameters.",
			Tags:        []string{entity},
			Responses: openapi3.Responses{
				"200": respBundle,
				"400": respErr,
				"401": respErr,
				"403": respErr,
				"404": respErr,
			},
		},
		Post: &openapi3.Operation{
			Description: "The create interaction creates a new resource " + entity + " extension",
			Tags:        []string{entity},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  content,
				},
			},
			Responses: responsesEntity,
		},
		Put: &openapi3.Operation{
			Description: "The update interaction creates or updates a resource " + entity + ".",
			Tags:        []string{entity},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  content,
				},
			},
			Responses: responsesEntity,
		},
	}
	g.Swagger.Paths["/"+entity+"/{id}"] = &openapi3.PathItem{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{Value: &openapi3.Parameter{
				Name:     "id",
				In:       "path",
				Required: true,
				Schema:   openapi3.NewSchemaRef("", &openapi3.Schema{Type: "string"}),
			}},
		},
		Get: &openapi3.Operation{
			Description: "This searches all resources of a particular type using the criteria represented in the parameters.",
			Tags:        []string{entity},
			Responses: openapi3.Responses{
				"200": respEntity,
				"400": respErr,
				"401": respErr,
				"403": respErr,
				"404": respErr,
			},
		},
		Post: &openapi3.Operation{
			Description: "The create interaction creates a new resource " + entity + " extension",
			Tags:        []string{entity},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  content,
				},
			},
			Responses: responsesEntity,
		},
		Put: &openapi3.Operation{
			Description: "The update interaction creates or updates a resource " + entity + ".",
			Tags:        []string{entity},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  content,
				},
			},
			Responses: responsesEntity,
		},
		Patch: &openapi3.Operation{
			Description: "The patch interaction patches a resource " + entity + ".",
			Tags:        []string{entity},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  content,
				},
			},
			Responses: responsesEntity,
		},
		Delete: &openapi3.Operation{
			Description: "The patch interaction patches a resource " + entity + ".",
			Tags:        []string{entity},
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Required: true,
					Content:  content,
				},
			},
			Responses: openapi3.Responses{
				"200": &openapi3.ResponseRef{Value: &openapi3.Response{
					Description: ptr.String("OK"),
				}},
				"202": &openapi3.ResponseRef{Value: &openapi3.Response{
					Description: ptr.String("OK"),
				}},
				"204": &openapi3.ResponseRef{Value: &openapi3.Response{
					Description: ptr.String("OK"),
				}},
				"400": respErr,
				"403": respErr,
				"404": respErr,
				"405": respErr,
				"409": respErr,
				"412": respErr,
			},
		},
	}
}
