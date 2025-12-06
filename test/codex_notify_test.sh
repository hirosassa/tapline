#!/bin/bash
# Test script for Codex CLI notify integration

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TAPLINE="$PROJECT_ROOT/tapline"
LOGFILE="codex_notify_test.log"

cleanup() {
    rm -f "$LOGFILE"
    rm -f ~/.tapline/session_id
}

trap cleanup EXIT

echo "=== Testing Codex CLI notify integration ==="

# Test 1: Check tapline binary exists
echo "Test 1: Tapline binary exists"
if [ ! -x "$TAPLINE" ]; then
    echo "FAIL: Tapline binary not found: $TAPLINE"
    exit 1
fi
echo "PASS: Tapline binary exists"

# Test 2: Start a session
echo "Test 2: Start session"
$TAPLINE conversation_start > "$LOGFILE"
if [ ! -f ~/.tapline/session_id ]; then
    echo "FAIL: Session file not created"
    exit 1
fi
session_id=$(cat ~/.tapline/session_id)
echo "PASS: Session started with ID: $session_id"

# Test 3: Simulate agent-turn-complete event
echo "Test 3: Simulate agent-turn-complete event"
EVENT_JSON='{
  "type": "agent-turn-complete",
  "data": {
    "response": "This is a test response from Codex CLI"
  }
}'

echo "$EVENT_JSON" | $TAPLINE notify-codex

# Give it a moment to process
sleep 1

# Verify by checking if we can still access the session
if [ ! -f ~/.tapline/session_id ]; then
    echo "FAIL: Session was unexpectedly cleared"
    exit 1
fi
echo "PASS: Event processed successfully"

# Test 4: Simulate unknown event (should be ignored)
echo "Test 4: Handle unknown event type"
UNKNOWN_EVENT='{
  "type": "unknown-event",
  "data": {}
}'

echo "$UNKNOWN_EVENT" | $TAPLINE notify-codex
echo "PASS: Unknown event handled gracefully"

# Test 5: End session
echo "Test 5: End session"
$TAPLINE conversation_end >> "$LOGFILE"
echo "PASS: Session ended"

echo ""
echo "=== Codex notify tests passed! ==="
echo "Summary:"
echo "  - Tapline binary: YES"
echo "  - Session management: YES"
echo "  - Event processing: YES"
echo "  - Error handling: YES"
echo ""
echo "Configuration:"
echo "  Add this to ~/.codex/config.toml:"
echo "  notify = \"$TAPLINE notify-codex\""
