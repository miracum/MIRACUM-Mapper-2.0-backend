//go:build tools
// +build tools

package tools

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=oapi-codegen-server.cfg.yaml ../api/openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=oapi-codegen-types.cfg.yaml ../api/openapi.yaml
