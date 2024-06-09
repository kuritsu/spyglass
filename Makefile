SHELL := /bin/bash
.DEFAULT_GOAL := build
export MONGODB_CONNECTIONSTRING ?= mongodb://spyglass:spyglass@localhost:27017/spyglass?authSource=admin

.PHONY: api build mod test db ui

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

api: build
	./spyglass server -v DEBUG

sch:
	./spyglass scheduler -v DEBUG -l shell

ui-build:
	cd ui; ember build --environment=production

ui:
	./spyglass ui
