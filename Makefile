build:
	golint ./...
	go build

get:
	go get -v -t -d ./...

test:
	mkdir -p coverage
	go test -v ./... -coverprofile=coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
