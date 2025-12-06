.PHONY: build test test-integration test-crash test-all clean install fmt vet

# Binary name
BINARY_NAME=tapline

# Build the binary
build:
	go build -o $(BINARY_NAME) ./cmd/tapline

# Run unit tests
test:
	go test -v ./...

# Run integration tests
test-integration: build
	./test/integration_test.sh

# Run crash resilience tests
test-crash: build
	./test/crash_resilience_test.sh

# Run all tests (unit + integration + crash)
test-all: test test-integration test-crash

# Run tests with coverage
test-coverage:
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -f test/*.log
	rm -f ~/.tapline/session_id

# Install binary to $GOPATH/bin
install:
	go install ./cmd/tapline

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run all checks (fmt, vet, test)
check: fmt vet test

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 ./cmd/tapline
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 ./cmd/tapline
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 ./cmd/tapline
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-linux-arm64 ./cmd/tapline
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe ./cmd/tapline

# Help
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
