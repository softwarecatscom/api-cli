package generator

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config openapi-cfg.yaml https://raw.githubusercontent.com/hostinger/api/refs/heads/main/openapi.json

//go:generate go run generator/docs.go
