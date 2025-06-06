FROM gcr.io/distroless/base-debian11 as final
ENTRYPOINT ["/rosinacli"]

FROM golang:1.23 AS builder
WORKDIR /app
COPY "go.mod" "go.mod"
COPY "go.sum" "go.sum"
RUN go mod download
RUN go mod verify

COPY "internal" "internal"
COPY "pkg" "pkg"
COPY "cmd/rosinacli" "cmd/rosinacli"

RUN CGO_ENABLED=0 go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /tmp/rosinacli \
	cmd/rosinacli/main.go

FROM final
COPY --from=builder /tmp/rosinacli /rosinacli

