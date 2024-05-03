SHELL := /bin/bash

build:
	go vet github.com/kuritsu/spyglass
	go build

mod:
	go mod tidy

test:
	mkdir -p coverage
	go test -v ./... -coverprofile=coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html

db:
	docker-compose up -d

db-clean:
	docker-compose down
	rm -rf data

.PHONY: api
api: build
	bash ./scripts/api.sh

ui-build:
	cd ui; ember build --environment=production

.PHONY: ui
ui:
	./spyglass ui
