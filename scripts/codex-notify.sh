#!/bin/bash
# Codex CLI Notify Handler for Tapline Logging
#
# This script is called by Codex CLI when events occur (e.g., agent-turn-complete)
# and logs them to tapline.
#
# Installation:
#   1. Make this script executable: chmod +x codex-notify.sh
#   2. Add to ~/.codex/config.toml:
#      notify = "/path/to/tapline/scripts/codex-notify.sh"
#
# Usage:
#   This script is automatically called by Codex CLI when events occur.
#   It receives event data via stdin as JSON.

set -euo pipefail

# Check if tapline is available
if ! command -v tapline >/dev/null 2>&1; then
    exit 0
fi

# Check if stdin has data
if [ -t 0 ]; then
    exit 0
fi

# Read event data from stdin
EVENT_DATA=$(cat)

# Parse event type and data
# Codex CLI sends JSON event data to notify scripts
EVENT_TYPE=$(echo "$EVENT_DATA" | jq -r '.type // "unknown"' 2>/dev/null || echo "unknown")

case "$EVENT_TYPE" in
    "agent-turn-complete")
        # Extract response from event data
        RESPONSE=$(echo "$EVENT_DATA" | jq -r '.data.response // .response // ""' 2>/dev/null || echo "")

        if [ -n "$RESPONSE" ]; then
            # Ensure session exists
            if [ ! -f ~/.tapline/session_id ]; then
                tapline conversation_start >/dev/null 2>&1 || true
            fi

            # Log assistant response
            tapline assistant_response "$RESPONSE" >/dev/null 2>&1 || true
        fi
        ;;

    "session_start")
        tapline conversation_start >/dev/null 2>&1 || true
        ;;

    "session_end")
        tapline conversation_end >/dev/null 2>&1 || true
        ;;

    *)
        # Unknown event type, ignore
        ;;
esac

exit 0
