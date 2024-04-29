SHELL := /bin/bash

build:
	go vet ./...
	go build

mod:
	go mod tidy

test:
	mkdir -p coverage
	go test -v ./... -coverprofile=coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html

db:
	docker-compose up -d

.PHONY: api
api: build
	bash ./scripts/api.sh
