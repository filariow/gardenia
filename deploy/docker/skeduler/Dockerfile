FROM gcr.io/distroless/static as final
ENTRYPOINT ["/go/bin/skeduler"]

FROM golang:1.18 as builder
WORKDIR /skeduler
COPY "go.mod" "go.mod"
COPY "go.sum" "go.sum"
RUN go mod download
RUN go mod verify

COPY "cmd/skeduler" "cmd/skeduler"
COPY "internal" "internal"
COPY "pkg" "pkg"

RUN go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /tmp/skeduler \
	cmd/skeduler/main.go

FROM final
COPY --from=builder /tmp/skeduler /go/bin/skeduler

