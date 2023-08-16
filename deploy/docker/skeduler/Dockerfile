FROM gcr.io/distroless/base-debian11 as final
ENTRYPOINT ["/skeduler"]

FROM golang:1.21 as builder
WORKDIR /app
COPY "go.mod" "go.mod"
COPY "go.sum" "go.sum"
RUN go mod download
RUN go mod verify

COPY "cmd/skeduler" "cmd/skeduler"
COPY "internal" "internal"
COPY "pkg" "pkg"

RUN CGO_ENABLED=0 go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /tmp/skeduler \
	cmd/skeduler/main.go

FROM final
COPY --from=builder /tmp/skeduler /skeduler

