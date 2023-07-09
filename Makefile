GO := go
GOARCH := arm
GOOS := linux

TARGET_RPI ?= rpi4

.PHONY: trybuild build protos

trybuild:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o /dev/null cmd/valvectl/main.go

build:
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

.PHONY: flowmeter
flowmeter:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -trimpath -ldflags="-s -w" -o bin/flowmeter cmd/flowmeter/main.go

.PHONY: flowmeter-rsync
flowmeter-rsync: flowmeter
	rsync ./bin/flowmeter root@$(TARGET_RPI):/usr/local/bin/flowmeter
