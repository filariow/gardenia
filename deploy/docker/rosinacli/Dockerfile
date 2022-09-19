FROM gcr.io/distroless/base-debian11 as final
ENTRYPOINT ["/rosinacli"]

FROM golang:1.18 as builder
WORKDIR /app
COPY "go.mod" "go.mod"
COPY "go.sum" "go.sum"
RUN go mod download
RUN go mod verify

COPY "cmd/rosinacli" "cmd/rosinacli"
COPY "internal" "internal"
COPY "pkg" "pkg"

RUN go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /tmp/rosinacli \
	cmd/rosinacli/main.go

FROM final
COPY --from=builder /tmp/rosinacli /rosinacli

