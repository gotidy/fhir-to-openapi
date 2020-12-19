# Generation FHIR R4 OpenAPI 3 specification

[![GoDev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)][godev]

[godev]: https://pkg.go.dev/github.com/gotidy/copy

Generation FHIR 4 OpenAPI 3 specification based on http://hl7.org/fhir/.

## Generation

```sh
# Install fhir-to-openapi
go get -u github.com/gotidy/fhir-to-openapi/cmd/fhir-to-openapi/...

# Generate OpenAPI 3 specification as YAML file.
fhir-to-openapi -i ./fhir.schema.json -o ./fhir.schema.oapi.yaml
# Or as JSON file.
fhir-to-openapi -i ./fhir.schema.json -o ./fhir.schema.oapi.json
```

or

```sh
git clone https://github.com/gotidy/fhir-to-openapi.git
cd fhir-to-openapi
# Generate OpenAPI 3 specification.
make run
# That variant generate OpenAPI 3 specification and generate Go Client and Server Code (uses https://github.com/deepmap/oapi-codegen).
make run-gen
```

## License

[Apache 2.0](https://github.com/gotidy/fhir-to-openapi/blob/master/LICENSE)
