FROM gcr.io/distroless/base-debian11 as final
ENTRYPOINT ["/rosina"]

FROM golang:1.23 AS builder
WORKDIR /app
COPY "go.mod" "go.mod"
COPY "go.sum" "go.sum"
RUN go mod download
RUN go mod verify

COPY "internal" "internal"
COPY "pkg" "pkg"
COPY "cmd/rosina" "cmd/rosina"

RUN go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /tmp/rosina \
	cmd/rosina/main.go

FROM final
COPY --from=builder /tmp/rosina /rosina

