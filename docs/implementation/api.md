# Open API Specification

The OpenAPI Specification is a standard for defining RESTful APIs. The [specification](/api/openapi.yaml) is written in YAML and defines the endpoints, request and response objects, query and path parameters as well as the HTTP methods that can be used to interact with the API. 

## Code generation with oapi-codegen
The API specification is used to generate boilerplate code for the server using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen). The generated code is placed in the `internal/api` directory. A `types_gen.go` file gets generated defining the request and response objects of the OpenAPI files in Go structs. The `server_gen.go` file contains the Gin router handlers and validation of the query and path parameters.

The `tools` directory scripts to generate the code for the project. In order to use these tools, just run this command:

```bash
cd tools
go generate -tags tools
```
