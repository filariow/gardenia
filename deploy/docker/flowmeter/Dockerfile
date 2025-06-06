FROM gcr.io/distroless/base-debian11 as final
ENTRYPOINT ["/flowmeter"]

FROM golang:1.24.4 AS builder
WORKDIR /app
COPY "go.mod" "go.mod"
COPY "go.sum" "go.sum"
RUN go mod download
RUN go mod verify

COPY "cmd/flowmeter" "cmd/flowmeter"

RUN go build \
	-trimpath \
	-ldflags="-s -w" \
	-o bin/flowmeter \
	cmd/flowmeter/main.go

FROM final
COPY --from=builder /app/bin/flowmeter /flowmeter

