FROM golang:1.22.6-bookworm as builder

WORKDIR /app

ENV GOPATH /go
ENV GOCACHE /gocache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/gocache \
    go build -tags="no_duckdb_arrow" -o ./bin/server .

FROM debian:bookworm-slim

COPY --from=builder /app/bin/server /usr/local/bin/server

CMD ["/usr/local/bin/server"]
