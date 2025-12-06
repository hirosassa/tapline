.PHONY: build test test-integration test-crash test-gemini test-codex test-all clean install
.PHONY: fmt lint check deps build-all help

BINARY_NAME=tapline

build:
	go build -o $(BINARY_NAME) ./cmd/tapline

build-all:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 ./cmd/tapline
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 ./cmd/tapline
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 ./cmd/tapline
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-linux-arm64 ./cmd/tapline
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe ./cmd/tapline

install:
	go install ./cmd/tapline

test:
	go test -v ./...

test-integration: build
	./test/integration_test.sh

test-crash: build
	./test/crash_resilience_test.sh

test-gemini: build
	./test/gemini_wrapper_test.sh

test-codex: build
	./test/codex_notify_test.sh

test-all: test test-integration test-crash test-gemini test-codex

test-coverage:
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...

lint:
	golangci-lint run --timeout=5m

lint-fix:
	golangci-lint run --timeout=5m --fix

check: fmt lint test

ci: check test-all

deps:
	go mod download
	go mod tidy

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html
	rm -f test/*.log
	rm -f ~/.tapline/session_id

help:
	@echo "Available targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  build            - Build the binary for current platform"
	@echo "  build-all        - Build for multiple platforms (darwin, linux, windows)"
	@echo "  install          - Install binary to GOPATH/bin"
	@echo ""
	@echo "Test targets:"
	@echo "  test             - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-crash       - Run crash resilience tests"
	@echo "  test-gemini      - Run Gemini CLI wrapper tests"
	@echo "  test-codex       - Run Codex CLI notify tests"
	@echo "  test-all         - Run all tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo ""
	@echo "Code quality targets:"
	@echo "  fmt              - Format code with go fmt"
	@echo "  lint             - Run golangci-lint"
	@echo "  lint-fix         - Run golangci-lint with auto-fix"
	@echo "  check            - Run fmt + lint + test"
	@echo "  ci               - Run check + test-all (CI pipeline)"
	@echo ""
	@echo "Other targets:"
	@echo "  deps             - Download and tidy dependencies"
	@echo "  clean            - Remove build artifacts and test files"
	@echo "  help             - Show this help message"
