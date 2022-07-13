app/api/server.gen.go: app/api/coding-dojo-api.yaml
	cd app/api && \
		oapi-codegen --config server.cfg.yaml coding-dojo-api.yaml && \
		oapi-codegen --config types.cfg.yaml coding-dojo-api.yaml

build::
	cd app && go build -o bin/main .

run::
	cd app && go run .
