# Tapline Quick Start Guide

## Installation

### Option 1: Using the install script (recommended)

```bash
./scripts/install.sh
```

### Option 2: Manual installation

```bash
# Build the binary
make build

# Install to GOPATH/bin
make install

# Or copy manually to your PATH
cp tapline /usr/local/bin/
```

### Option 3: Build from source

```bash
go build -o tapline ./cmd/tapline
```

## Setup for Claude Code

1. The `.claude/hooks.json` file is already configured in this repository
2. Build and install the binary (see above)
3. Start using Claude Code in this directory - logs will be automatically collected!

## Testing the Installation

### Quick Test

Run the following commands to test:

```bash
# Start a conversation
tapline conversation_start

# Log a user prompt
tapline user_prompt "Hello, world!"

# Log an assistant response
tapline assistant_response "Hello! How can I help?"

# End the conversation
tapline conversation_end
```

Each command will output a JSON line to stdout.

### Comprehensive Test Suite

Run all tests including crash resilience:

```bash
# Run all tests
make test-all

# Or run individually
make test              # Unit tests
make test-integration  # Integration tests
make test-crash        # Crash resilience tests
```

## Collecting Logs

**Important:** Tapline logs are immediately flushed to disk after each event, ensuring no data loss even if the process crashes or the terminal is closed unexpectedly.

### Option 1: Redirect to file

```bash
# Claude Code will automatically call tapline via hooks
# Redirect stdout to a log file when running Claude Code
claude > conversation.log 2>&1
```

### Option 2: Use tee to see logs in real-time

```bash
claude 2>&1 | tee conversation.log
```

### Option 3: Process logs in real-time with jq

```bash
claude 2>&1 | grep -E '^{' | jq .
```

## Analyzing Logs

See `examples/jq_queries.sh` for common analysis patterns.

### Quick examples:

```bash
# View all logs pretty-printed
cat conversation.log | grep -E '^{' | jq .

# Extract user prompts only
cat conversation.log | grep -E '^{' | jq 'select(.role=="user") | .content'

# Count messages by role
cat conversation.log | grep -E '^{' | jq -s 'group_by(.role) | map({role: .[0].role, count: length})'

# Filter by specific session
cat conversation.log | grep -E '^{' | jq 'select(.session_id=="YOUR-SESSION-ID")'
```

## Understanding the Log Format

Each log entry is a JSON object with these fields:

- `timestamp`: ISO 8601 timestamp
- `service`: Service identifier (e.g., "claude-code")
- `session_id`: UUID for the conversation session
- `role`: "user", "assistant", or "system"
- `content`: The message content
- `metadata`: Optional metadata (e.g., hostname, cwd)
- `event`: Optional event type (e.g., "session_start", "session_end")

### Example log entry:

```json
{
  "timestamp": "2025-12-06T16:19:02.935561+09:00",
  "service": "claude-code",
  "session_id": "0e97d08c-b08b-4f5a-92ea-086c36d5818b",
  "role": "system",
  "content": "",
  "metadata": {
    "cwd": "/Users/username/project",
    "hostname": "my-laptop"
  },
  "event": "session_start"
}
```

## Session Management

- Sessions are automatically created when a conversation starts
- The session ID is stored in `~/.tapline/session_id`
- Sessions are cleared when a conversation ends
- All logs within a conversation share the same session ID

## Troubleshooting

### Check if hooks are working

```bash
# Look for hook output in Claude Code stderr
claude 2>&1 | grep tapline
```

### Verify session file

```bash
# Check current session
cat ~/.tapline/session_id

# Manually clear session if stuck
rm ~/.tapline/session_id
```

### Test commands manually

```bash
# Each command should output valid JSON
tapline conversation_start
tapline user_prompt "test"
tapline conversation_end
```

## Next Steps

- Integrate with log aggregation systems (Elasticsearch, Splunk, etc.)
- Set up automated log rotation
- Create dashboards for conversation analytics
- Extend to other AI services (Gemini CLI, ChatGPT, etc.)

## Development

```bash
# Run tests
make test

# Format code
make fmt

# Run all checks
make check

# Build for multiple platforms
make build-all
```

For more details, see the main [README.md](README.md).
