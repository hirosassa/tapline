# Gemini CLI Integration

This document explains how to use tapline with Google's Gemini CLI for conversation logging.

## Status: Temporary Solution

**IMPORTANT**: This integration uses a wrapper script as a temporary measure. Once Gemini CLI implements native hooks support (tracked in [google-gemini/gemini-cli#2779](https://github.com/google-gemini/gemini-cli/issues/2779)), we recommend migrating to the native hooks system for better integration.

## Prerequisites

- [Gemini CLI](https://github.com/google-gemini/gemini-cli) installed and configured
- tapline installed and in your PATH

## Installation

### Step 1: Source the Wrapper Script

Add the following to your shell configuration file (`~/.bashrc`, `~/.zshrc`, etc.):

```bash
# Source tapline gemini wrapper
if [ -f /path/to/tapline/scripts/gemini-wrapper.sh ]; then
    source /path/to/tapline/scripts/gemini-wrapper.sh
fi
```

Replace `/path/to/tapline` with the actual path to your tapline installation.

### Step 2: Reload Your Shell

```bash
source ~/.bashrc  # or ~/.zshrc
```

## Usage

### Option 1: Use Wrapper Functions (Recommended)

Use the `gemini_with_logging` function instead of the regular `gemini` command:

```bash
# Start a new session
gemini_start_session

# Run gemini with logging
gemini_with_logging "Explain quantum computing"

# End the session
gemini_end_session
```

### Option 2: Override the gemini Command

If you want to automatically log all `gemini` commands, uncomment the alias line in the wrapper script:

```bash
# In scripts/gemini-wrapper.sh, uncomment this line:
alias gemini='gemini_with_logging'
```

Then use `gemini` normally - all interactions will be logged automatically.

## How It Works

The wrapper script:

1. **Session Management**: Automatically starts a tapline session if one doesn't exist
2. **Input Logging**: Logs your prompts as user messages
3. **Output Capture**: Captures Gemini's responses and logs them as assistant messages
4. **Transparent Pass-through**: All Gemini CLI functionality works normally

## Log Format

Logs are written in JSON Lines format to stdout:

```json
{"time":"2025-12-06T22:00:00.123456+09:00","level":"INFO","msg":"conversation","service":"gemini-cli","session_id":"uuid","role":"user","content":"Your prompt"}
{"time":"2025-12-06T22:00:01.234567+09:00","level":"INFO","msg":"conversation","service":"gemini-cli","session_id":"uuid","role":"assistant","content":"Gemini's response"}
```

## Limitations

This wrapper approach has some limitations:

1. **Interactive Sessions**: Works best with single-prompt commands. Interactive multi-turn conversations within one `gemini` invocation are captured as a single response.

2. **Tool Usage**: Internal tool calls and function executions are not individually logged - only the final output is captured.

3. **Streaming**: Real-time streaming output is preserved but logged only after completion.

4. **Exit Codes**: Gemini CLI exit codes are properly forwarded.

## Troubleshooting

### Wrapper Not Found

If you see "Warning: tapline not found", ensure:
- tapline is installed and in your PATH
- Run `which tapline` to verify

### Gemini Command Not Found

If the wrapper can't find `gemini`:
- Ensure Gemini CLI is installed: `npm install -g @google/gemini-cli`
- Verify it's in PATH: `which gemini`

### Logs Not Being Created

Check that:
- The wrapper script is sourced in your shell
- You're using `gemini_with_logging` or have enabled the alias
- tapline is working: `tapline version`

## Migration Path

When Gemini CLI adds native hooks support:

1. Review the [hooks implementation](https://github.com/google-gemini/gemini-cli/issues/2779)
2. Create a Gemini CLI hooks configuration similar to Claude Code's `.claude/hooks.json`
3. Remove the wrapper script and aliases
4. Update your configuration to use native hooks

The native hooks will provide:
- Better integration with Gemini CLI's internals
- Access to individual tool calls and function executions
- More accurate event timing
- No wrapper overhead

## Testing

Test the wrapper with:

```bash
./test/gemini_wrapper_test.sh
```

This creates a mock `gemini` command and verifies the logging functionality.

## References

- [Gemini CLI GitHub Repository](https://github.com/google-gemini/gemini-cli)
- [Gemini CLI Hooks Feature Request](https://github.com/google-gemini/gemini-cli/issues/2779)
- [Tapline Documentation](../README.md)
