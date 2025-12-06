# Codex CLI Integration

This document explains how to use tapline with OpenAI's Codex CLI for conversation logging.

## Overview

Codex CLI integration uses the `notify` configuration option to trigger tapline logging when Codex events occur. This approach leverages Codex CLI's built-in notification system for clean integration.

## Prerequisites

- [Codex CLI](https://github.com/openai/codex) installed and configured
- tapline installed and in your PATH

## Installation

### Step 1: Configure Codex CLI

Add the notify handler to your Codex CLI configuration file (`~/.codex/config.toml`):

```toml
# Tapline integration
notify = "tapline notify-codex"
```

If `tapline` is not in your PATH, use the absolute path:

```toml
notify = "/path/to/tapline notify-codex"
```

You can use the example configuration as a starting point:

```bash
# Copy example configuration
cp examples/codex-config.toml ~/.codex/config.toml

# Edit and update the command
nano ~/.codex/config.toml
```

### Step 2: Test the Integration

```bash
# Run the test script
./test/codex_notify_test.sh
```

## How It Works

The integration uses Codex CLI's `notify` configuration:

1. **Event Detection**: Codex CLI calls `tapline notify-codex` when events occur
2. **Event Processing**: The command receives event data as JSON via stdin
3. **Logging**: Events are logged to tapline based on event type
4. **Session Management**: Sessions are automatically managed

### Supported Events

Currently, the notify handler supports:

- `agent-turn-complete`: Logged as assistant responses
- `session_start`: Creates a new tapline session
- `session_end`: Ends the current tapline session

## Log Format

Logs are written in JSON Lines format to stdout:

```json
{"time":"2025-12-06T22:00:00.123456+09:00","level":"INFO","msg":"conversation","service":"codex-cli","session_id":"uuid","role":"assistant","content":"Codex's response"}
```

## Configuration Options

### Basic Configuration

Minimal configuration in `~/.codex/config.toml`:

```toml
notify = "tapline notify-codex"
```

### With TUI Notifications

Enable both tapline logging and TUI notifications:

```toml
notify = "tapline notify-codex"

[tui]
notifications = true
```

### Full Example

See `examples/codex-config.toml` for a complete configuration example with additional Codex CLI settings.

## Advantages Over Other Integrations

Compared to Gemini CLI's wrapper approach, Codex CLI integration offers:

1. **Native Integration**: Uses Codex CLI's official `notify` configuration
2. **No Wrapper Needed**: No need to wrap or alias the codex command
3. **Stability**: Less likely to break with Codex CLI updates
4. **Official Support**: Built on documented Codex CLI features
5. **Single Binary**: Everything included in `tapline` binary via `go install`

## Limitations

1. **Event Coverage**: Currently only `agent-turn-complete` is supported by Codex CLI
2. **User Prompts**: User prompts are not captured (only assistant responses)
3. **Tool Calls**: Individual tool calls are not separately logged

## Future Improvements

When Codex CLI adds more event types or hooks:

- User prompt events
- Tool execution events
- Error events
- More detailed event metadata

Check the [Codex CLI GitHub repository](https://github.com/openai/codex) for updates on event support.

## Troubleshooting

### Notify Script Not Called

Verify the configuration:

```bash
# Check config file exists
cat ~/.codex/config.toml

# Verify path in config
grep notify ~/.codex/config.toml

# Check script is executable
ls -la /path/to/tapline/scripts/codex-notify.sh
```

### No Logs Being Created

Check that:
- tapline is in your PATH: `which tapline`
- The notify script has execute permissions
- Codex CLI is actually running (not just planning)

### Session Not Created

The notify script automatically creates a session if needed. If sessions aren't being created:

```bash
# Manually start a session
tapline conversation_start

# Check session file
ls -la ~/.tapline/session_id
```

## Testing

Test the notify integration:

```bash
./test/codex_notify_test.sh
```

This creates mock events and verifies the notify script processes them correctly.

## References

- [Codex CLI Documentation](https://developers.openai.com/codex/cli/)
- [Codex CLI Configuration](https://developers.openai.com/codex/local-config/)
- [Codex CLI GitHub Repository](https://github.com/openai/codex)
- [Tapline Documentation](../README.md)
