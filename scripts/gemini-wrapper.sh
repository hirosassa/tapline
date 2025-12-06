#!/bin/bash
# Gemini CLI Wrapper Helper Functions for Tapline Logging
#
# TEMPORARY MEASURE: This is a temporary solution until Gemini CLI
# implements native hooks support (see https://github.com/google-gemini/gemini-cli/issues/2779)
#
# This provides shell functions to wrap Gemini CLI for conversation logging via tapline.
# Once Gemini CLI adds hooks support, migrate to using native hooks instead.
#
# Installation:
#   Add this to your ~/.bashrc or ~/.zshrc:
#
#   # Source tapline gemini wrapper
#   if [ -f /path/to/tapline/scripts/gemini-wrapper.sh ]; then
#       source /path/to/tapline/scripts/gemini-wrapper.sh
#   fi
#
# Usage:
#   gemini_with_logging "Your prompt here"
#   # Or use the alias if you want to override 'gemini' command:
#   # alias gemini='gemini_with_logging'

# Check if tapline is available
_tapline_check() {
    if ! command -v tapline >/dev/null 2>&1; then
        return 1
    fi
    return 0
}

# Ensure session is started
_tapline_ensure_session() {
    if ! _tapline_check; then
        return 1
    fi

    # Check if session already exists
    if [ -f ~/.tapline/session_id ]; then
        return 0
    fi

    # Start new session
    tapline conversation_start >/dev/null 2>&1
    return $?
}

# Main wrapper function for Gemini CLI
gemini_with_logging() {
    # Check if real gemini command exists
    if ! command -v gemini >/dev/null 2>&1; then
        echo "Error: 'gemini' command not found in PATH" >&2
        return 1
    fi

    # If tapline is not available, just run gemini normally
    if ! _tapline_check; then
        echo "Warning: tapline not found, logging disabled" >&2
        command gemini "$@"
        return $?
    fi

    # Ensure session is started
    _tapline_ensure_session

    # If arguments provided, log as user prompt
    if [ $# -gt 0 ]; then
        local prompt="$*"
        tapline user_prompt "$prompt" >/dev/null 2>&1 || true
    fi

    # Run gemini and capture output
    local temp_output
    temp_output=$(mktemp)

    # Run gemini, tee output to file and stdout
    if command gemini "$@" 2>&1 | tee "$temp_output"; then
        # Log the response
        if [ -s "$temp_output" ]; then
            local response
            response=$(cat "$temp_output")
            tapline assistant_response "$response" >/dev/null 2>&1 || true
        fi
        rm -f "$temp_output"
        return 0
    else
        local exit_code=$?
        rm -f "$temp_output"
        return $exit_code
    fi
}

# Function to start a new Gemini session with logging
gemini_start_session() {
    if _tapline_check; then
        tapline conversation_start >/dev/null 2>&1
        echo "Started new Gemini session with tapline logging"
    else
        echo "Warning: tapline not found" >&2
    fi
}

# Function to end current Gemini session
gemini_end_session() {
    if _tapline_check; then
        tapline conversation_end >/dev/null 2>&1
        echo "Ended Gemini session"
    fi
}

# Optional: Uncomment to override the 'gemini' command
# alias gemini='gemini_with_logging'
