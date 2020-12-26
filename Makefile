SHELL := /bin/bash

build:
	golint ./...
	go build

get:
	go get -v -t -d ./...

test:
	mkdir -p coverage
	go test -v ./... -coverprofile=coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html

	export  = mongodb://spyglass:spyglass@localhost:27017/spyglass?authSource=admin

.PHONY: api
api: build
	bash ./scripts/api.sh
