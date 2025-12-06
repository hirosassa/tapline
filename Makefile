.PHONY: build test test-integration test-crash test-all clean install fmt vet

BINARY_NAME=tapline

build:
	go build -o $(BINARY_NAME) ./cmd/tapline

test:
	go test -v ./...

test-integration: build
	./test/integration_test.sh

test-crash: build
	./test/crash_resilience_test.sh

test-all: test test-integration test-crash

test-coverage:
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -f test/*.log
	rm -f ~/.tapline/session_id

install:
	go install ./cmd/tapline

fmt:
	go fmt ./...

vet:
	go vet ./...

deps:
	go mod download
	go mod tidy

check: fmt vet test

build-all:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 ./cmd/tapline
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 ./cmd/tapline
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 ./cmd/tapline
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-linux-arm64 ./cmd/tapline
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe ./cmd/tapline

help:
	@echo "Available targets:"
	@echo "  build            - Build the binary"
	@echo "  test             - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-crash       - Run crash resilience tests"
	@echo "  test-all         - Run all tests (unit + integration + crash)"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  clean            - Remove build artifacts and test files"
	@echo "  install          - Install binary to GOPATH/bin"
	@echo "  fmt              - Format code"
	@echo "  vet              - Run go vet"
	@echo "  deps             - Download and tidy dependencies"
	@echo "  check            - Run fmt, vet, and test"
	@echo "  build-all        - Build for multiple platforms"
	@echo "  help             - Show this help message"
