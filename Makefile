
ALL_PACKAGES=$(shell go list ./... | grep -v "mocks")

all: fmt vet test

up: download-dataset
	@docker compose up

down:
	@docker compose down

.PHONY: test
test:
	@echo "== Running all tests =="
	@# https://github.com/golang/go/issues/61229#issuecomment-1988965927
	@go test -race -ldflags=-extldflags=-Wl,-w -failfast -coverprofile=coverage.out $(ALL_PACKAGES)
	@$(MAKE) coverage

coverage:
	@echo "\nTotal Coverage:" $(shell go tool cover -func=coverage.out | tail -n 1 | awk '{print $$3}')

fmt:
	@echo "== Running go fmt =="
	go fmt ./...

vet:
	@echo "== Running go vet =="
	go vet ./...

mockery:
	@type go > /dev/null 2>&1 || (echo "Please install Go first" && exit 1)
	@if [ ! -f "./bin/mockery" ]; then \
		GOBIN=$(shell pwd)/bin go install github.com/vektra/mockery/v2@v2.44.1; \
	fi
	@./bin/mockery

download-dataset:
	@echo "== Downloading dataset chicago_taxi_trips_2020.parquet =="
	@mkdir -p dataset
	@if [ ! -f "dataset/chicago_taxi_trips_2020.parquet" ]; then \
	  curl -L --progress-bar "https://drive.usercontent.google.com/download?id=1qaCwmmr6fBTzGFH171-vdZbPNeFIgiBs&export=download&confirm=y" -o "dataset/chicago_taxi_trips_2020.parquet"; \
	else \
	  echo "dataset/chicago_taxi_trips_2020.parquet already exists\n"; \
	fi
