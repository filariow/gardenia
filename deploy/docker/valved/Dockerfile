FROM gcr.io/distroless/base-debian11 as final
ENTRYPOINT ["/valved"]

FROM golang:1.16 as builder
WORKDIR /app
COPY "go.mod" "go.mod"
COPY "go.sum" "go.sum"
RUN go mod download
RUN go mod verify

COPY "cmd/valved" "cmd/valved"
COPY "internal" "internal"
COPY "pkg" "pkg"

RUN go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /tmp/valved \
	cmd/valved/main.go

FROM final
COPY --from=builder /tmp/valved /valved

