//go:generate oapi-codegen -generate types -package fhir -o types.gen.go fhir.schema.oapi.yaml
//go:generate oapi-codegen -generate client -package fhir -o client.gen.go fhir.schema.oapi.yaml
//go:generate goimports -w ./types.gen.go
//go:generate goimports -w ./client.gen.go
package fhir
