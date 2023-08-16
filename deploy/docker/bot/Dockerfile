FROM gcr.io/distroless/base-debian11 as final
ENTRYPOINT ["/bot"]

FROM golang:1.21 as builder
WORKDIR /app
COPY "go.mod" "go.mod"
COPY "go.sum" "go.sum"
RUN go mod download
RUN go mod verify

COPY "internal" "internal"
COPY "pkg" "pkg"
COPY "cmd/bot" "cmd/bot"

RUN CGO_ENABLED=0 go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /tmp/bot \
	cmd/bot/main.go

FROM final
COPY --from=builder /tmp/bot /bot

