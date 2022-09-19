FROM gcr.io/distroless/base-debian11 as final
ENTRYPOINT ["/metrics"]

FROM golang:1.18 as builder
WORKDIR /app
COPY "go.mod" "go.mod"
COPY "go.sum" "go.sum"
RUN go mod download
RUN go mod verify

COPY "cmd/metrics" "cmd/metrics"
COPY "pkg" "pkg"

RUN go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /tmp/metrics \
	cmd/metrics/main.go

FROM final
COPY --from=builder /tmp/metrics /metrics

