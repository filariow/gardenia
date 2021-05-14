GO := go
GOARCH := arm
GOOS := linux

.PHONY: trybuild build

trybuild: **.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o /dev/null cmd/valvectl/main.go

build: **.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -trimpath -ldflags="-s -w" -o bin/valvectl cmd/valvectl/main.go

ci: build
	scp bin/valvectl rpi:/home/fra/valvectl

