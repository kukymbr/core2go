all:
	make prepare
	make test

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

clean:
	go clean