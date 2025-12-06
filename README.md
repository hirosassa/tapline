# Tapline - AI Conversation Logger

[![CI](https://github.com/hirosassa/tapline/workflows/CI/badge.svg)](https://github.com/hirosassa/tapline/actions?query=workflow%3ACI)
[![Go Report Card](https://goreportcard.com/badge/github.com/hirosassa/tapline)](https://goreportcard.com/report/github.com/hirosassa/tapline)
[![codecov](https://codecov.io/gh/hirosassa/tapline/graph/badge.svg?token=349A1NN1JT)](https://codecov.io/gh/hirosassa/tapline)
[![GitHub release](https://img.shields.io/github/release/hirosassa/tapline.svg)](https://github.com/hirosassa/tapline/releases)

Tapline is a unified conversation logging system for multiple AI chat services (Claude Code, Gemini CLI, ChatGPT, etc.). It outputs structured logs in JSON Lines format to stdout, making it easy to integrate with log aggregation systems.

## Features

- Structured logging using Go's standard `log/slog` library
- JSON Lines format output for easy parsing
- **Crash-resilient logging**: Logs are immediately flushed to disk
- **Process-per-event model**: No buffering between events
- Session ID management for conversation tracking
- Claude Code integration via hooks system
- Extensible adapter pattern for future services
- Minimal dependencies - only `github.com/google/uuid` for session IDs

## Installation

### Download Pre-built Binaries (Recommended)

Download the latest release from the [Releases page](https://github.com/hirosassa/tapline/releases):

```bash
# Linux (amd64)
wget https://github.com/hirosassa/tapline/releases/latest/download/tapline_Linux_x86_64.tar.gz
tar xzf tapline_Linux_x86_64.tar.gz
sudo mv tapline /usr/local/bin/

# macOS (amd64)
wget https://github.com/hirosassa/tapline/releases/latest/download/tapline_Darwin_x86_64.tar.gz
tar xzf tapline_Darwin_x86_64.tar.gz
sudo mv tapline /usr/local/bin/

# macOS (arm64)
wget https://github.com/hirosassa/tapline/releases/latest/download/tapline_Darwin_arm64.tar.gz
tar xzf tapline_Darwin_arm64.tar.gz
sudo mv tapline /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/hirosassa/tapline/cmd/tapline@latest
```

### Build from Source

```bash
git clone https://github.com/hirosassa/tapline.git
cd tapline
make build
sudo cp tapline /usr/local/bin/
```

## Usage

### Supported Services

Tapline currently supports the following AI services:

1. **Claude Code** - Native integration via hooks system
2. **Codex CLI** - Native integration via notify configuration
3. **Gemini CLI** - Wrapper script integration (temporary solution)

See service-specific documentation:
- [Claude Code Integration](#claude-code-integration)
- [Codex CLI Integration](docs/CODEX_CLI.md)
- [Gemini CLI Integration](docs/GEMINI_CLI.md)

### Claude Code Integration

Tapline integrates with Claude Code using the hooks system. The configuration is in `.claude/hooks.json`:

```json
{
  "hooks": {
    "conversation_start": {
      "command": "tapline conversation_start",
      "enabled": true
    },
    "conversation_end": {
      "command": "tapline conversation_end",
      "enabled": true
    },
    "user_prompt_submit": {
      "command": "tapline user_prompt '{{prompt}}'",
      "enabled": true
    },
    "assistant_response": {
      "command": "tapline assistant_response '{{response}}'",
      "enabled": true
    }
  }
}
```

### Manual Usage

```bash
# Start a conversation (creates new session)
tapline conversation_start

# Log user prompt
tapline user_prompt "Hello, how are you?"

# Log assistant response
tapline assistant_response "I'm doing well, thank you!"

# End conversation (clears session)
tapline conversation_end
```

## Log Format

Each log entry is output as a single JSON line to stdout using Go's `log/slog`:

```json
{"time":"2025-12-06T16:26:36.768095+09:00","level":"INFO","msg":"conversation","service":"claude-code","session_id":"c0db0a0f-561b-44d3-b213-13fd3d9c0472","role":"user","content":"Hello!"}
{"time":"2025-12-06T16:26:37.768095+09:00","level":"INFO","msg":"conversation","service":"claude-code","session_id":"c0db0a0f-561b-44d3-b213-13fd3d9c0472","role":"assistant","content":"Hi there!"}
```

### Log Schema

- `time`: ISO 8601 timestamp (automatically added by slog)
- `level`: Log level (always "INFO" for conversation logs)
- `msg`: Message type (always "conversation")
- `service`: Service identifier (e.g., "claude-code", "gemini-cli")
- `session_id`: UUID for the conversation session
- `role`: "user", "assistant", or "system"
- `content`: The message content
- `metadata`: Optional metadata object (for session events)
- `event`: Optional event type (e.g., "session_start", "session_end")

## Session Management

Tapline manages session IDs persistently in `~/.tapline/session_id`. Sessions are:

- Created on `conversation_start`
- Used for all subsequent logs in the same conversation
- Cleared on `conversation_end`

**Crash Resilience:** Even if `conversation_end` is never called (due to crashes or terminal closure), all logs up to that point are preserved. Each hook runs as an independent process that immediately flushes logs to disk with `os.Stdout.Sync()`. 

## Log Processing Examples

### View logs in real-time

```bash
tail -f conversation.log | jq .
```

### Filter by session

```bash
cat conversation.log | jq 'select(.session_id=="550e8400-e29b-41d4-a716-446655440000")'
```

### Extract user prompts only

```bash
cat conversation.log | jq 'select(.role=="user") | .content'
```

### Group by service

```bash
cat conversation.log | jq -s 'group_by(.service) | map({service: .[0].service, count: length})'
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/logger
go test ./pkg/session
```

## Contributing

Contributions are welcome! Please ensure:

1. All tests pass: `make test-all`
2. Code is formatted: `go fmt ./...`
3. Linting passes: `golangci-lint run`
4. Documentation is updated

## License

MIT License - see [LICENSE](LICENSE) file for details
