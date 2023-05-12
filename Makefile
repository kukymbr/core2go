all:
	make prepare
	make lint
	make test
	make build

prepare:
	go mod tidy
	go vet ./...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
	golangci-lint run ./...

test:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm -f coverage.out

build:
	go build ./...
	go mod tidy

clean:
	go clean