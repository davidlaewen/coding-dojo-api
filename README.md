# Coding Dojo API
Reference implementation of a toy CRUD API in Go.

The server code and types are generated from an OpenAPI spec in `app/api/coding-dojo-api.yaml` using [oapi-codegen](https://github.com/deepmap/oapi-codegen).

## Setup
Install `oapi-codegen` for generating code from OpenAPI specification:
```sh
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
```

The code for the API types and server is generated from `user-registry.yml` in `/app/api` with:
```sh
$ oapi-codegen -generate types  -o types.gen.go  -package api user-registry.yml
$ oapi-codegen -generate server -o server.gen.go -package api user-registry.yml
```
Or using the config `.yml` files:

```sh
$ oapi-codegen --config server.cfg.yml user-registry.yml
$ oapi-codegen --config types.cfg.yml user-registry.yml
```

The generated code in `types.gen.go` defines struct types corresponding to the schemas given in the OpenAPI spec.
The code in `server.gen.go` defines an interface for the HTTP server which is implemented in `app/server.go`.
