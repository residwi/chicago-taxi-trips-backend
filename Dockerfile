# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder

WORKDIR /app

ENV GOPATH=/go
ENV GOCACHE=/gocache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/ \
  go mod download -x

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=cache,target=/gocache \
  go build -tags="no_duckdb_arrow" -o ./bin/server .

FROM debian:trixie-slim

WORKDIR /

COPY --from=builder /app/bin/server /server

EXPOSE 8080

CMD ["/server"]
