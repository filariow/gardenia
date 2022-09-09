GO := go
GOARCH := arm
GOOS := linux

.PHONY: trybuild build protos

trybuild:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o /dev/null cmd/valvectl/main.go

build: **.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -trimpath -ldflags="-s -w" -o bin/valvectl cmd/valvectl/main.go

ci: build
	scp bin/valvectl rpi:/home/fra/valvectl

protos:
	mkdir -p pkg/valvedprotos
	protoc \
		-I protos \
		--go_opt=paths=source_relative \
		--go_out=pkg/valvedprotos \
		--go-grpc_opt=paths=source_relative \
		--go-grpc_out=pkg/valvedprotos \
		protos/*.proto
	mkdir -p ./fe/gardenia-web/src/grpc
	protoc \
		-I protos \
		--plugin="protoc-gen-ts=./fe/gardenia-web/node_modules/.bin/protoc-gen-ts" \
		--ts_out="service=grpc-web:./fe/gardenia-web/src/grpc" \
		protos/valved.proto

