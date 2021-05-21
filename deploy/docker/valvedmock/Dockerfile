FROM gcr.io/distroless/static as final

FROM golang:1.16 as builder
WORKDIR /valved
COPY "go.mod" "go.mod"
COPY "go.sum" "go.sum"
RUN go mod download
RUN go mod verify

COPY "cmd/valvedmock" "cmd/valvedmock"
COPY "internal" "internal"
COPY "pkg" "pkg"

RUN go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /tmp/valved \
	cmd/valvedmock/main.go

FROM golang:1.16
COPY --from=builder /tmp/valved /go/bin/valved
ENTRYPOINT ["/go/bin/valved"]

