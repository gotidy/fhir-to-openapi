package generator

import "github.com/getkin/kin-openapi/openapi3"

type ParameterLocation string

const (
	InQuery ParameterLocation = "query"
	InPath  ParameterLocation = "path"
)

func NewParameterWithSchema(in ParameterLocation, name string, required bool, schema *openapi3.SchemaRef) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{Value: &openapi3.Parameter{
		Name:     name,
		In:       string(name),
		Required: false,
		Schema:   schema,
	}}
}

func NewParameterRef(refName string) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{Ref: "#/components/parameters/" + refName}
}

func NewRequestBodyWithContent(content openapi3.Content, required bool) *openapi3.RequestBodyRef {
	return &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Required: required,
			Content:  content,
		},
	}
}

func NewRequestBodyRef(refName string) *openapi3.RequestBodyRef {
	return &openapi3.RequestBodyRef{Ref: "#/components/requestBodies/" + refName}
}

func NewParameterInQuery(in ParameterLocation, name string, required bool, typ string) *openapi3.ParameterRef {
	return NewParameterWithSchema(InQuery, name, required, &openapi3.SchemaRef{Value: &openapi3.Schema{Type: typ}})
}

func NewParameterInPath(in ParameterLocation, name string, required bool, typ string) *openapi3.ParameterRef {
	return NewParameterWithSchema(InPath, name, required, NewSchema(typ))
}

func NewSchema(typ string) *openapi3.SchemaRef {
	return &openapi3.SchemaRef{Value: &openapi3.Schema{Type: typ}}
}

func NewSchemaWithFormat(typ string, format string) *openapi3.SchemaRef {
	return &openapi3.SchemaRef{Value: &openapi3.Schema{Type: typ, Format: format}}
}

func NewSchemaRef(refName string) *openapi3.SchemaRef {
	return &openapi3.SchemaRef{Ref: "#/components/schemas/" + refName}
}

func NewSchemaString() *openapi3.SchemaRef {
	return NewSchema("string")
}

func NewSchemaInteger() *openapi3.SchemaRef {
	return NewSchema("integer")
}

func NewSchemaNumber() *openapi3.SchemaRef {
	return NewSchema("number")
}

func NewSchemaBoolean() *openapi3.SchemaRef {
	return NewSchema("boolean")
}

func NewSchemaURI() *openapi3.SchemaRef {
	return NewSchemaWithFormat("string", "uri")
}

func NewSchemaEmail() *openapi3.SchemaRef {
	return NewSchemaWithFormat("string", "email")
}

func NewSchemaUUID() *openapi3.SchemaRef {
	return NewSchemaWithFormat("string", "uuid")
}

func NewSchemaDate() *openapi3.SchemaRef {
	return NewSchemaWithFormat("string", "date")
}

func NewSchemaDateTime() *openapi3.SchemaRef {
	return NewSchemaWithFormat("string", "date-time")
}
