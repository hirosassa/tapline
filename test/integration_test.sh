#!/bin/bash
# Integration test to verify immediate log flushing

set -e

TAPLINE="./tapline"
LOGFILE="test_output.log"

# Clean up function
cleanup() {
    rm -f "$LOGFILE"
    rm -f ~/.tapline/session_id
}

trap cleanup EXIT

echo "=== Testing immediate log flushing ==="

# Test 1: Start conversation and verify log is written immediately
echo "Test 1: Conversation start flushes immediately"
$TAPLINE conversation_start > "$LOGFILE"
if [ ! -s "$LOGFILE" ]; then
    echo "FAIL: Log file is empty after conversation_start"
    exit 1
fi
if ! grep -q "session_start" "$LOGFILE"; then
    echo "FAIL: session_start event not found in log"
    exit 1
fi
echo "PASS: Conversation start logged immediately"

# Test 2: User prompt is written immediately
echo "Test 2: User prompt flushes immediately"
$TAPLINE user_prompt "Test message" >> "$LOGFILE"
if ! grep -q "Test message" "$LOGFILE"; then
    echo "FAIL: User prompt not found in log"
    exit 1
fi
line_count=$(wc -l < "$LOGFILE")
if [ "$line_count" -lt 2 ]; then
    echo "FAIL: Expected at least 2 log lines"
    exit 1
fi
echo "PASS: User prompt logged immediately"

# Test 3: Assistant response is written immediately
echo "Test 3: Assistant response flushes immediately"
$TAPLINE assistant_response "Test response" >> "$LOGFILE"
if ! grep -q "Test response" "$LOGFILE"; then
    echo "FAIL: Assistant response not found in log"
    exit 1
fi
echo "PASS: Assistant response logged immediately"

# Test 4: Conversation end is written immediately
echo "Test 4: Conversation end flushes immediately"
$TAPLINE conversation_end >> "$LOGFILE"
if ! grep -q "session_end" "$LOGFILE"; then
    echo "FAIL: session_end event not found in log"
    exit 1
fi
echo "PASS: Conversation end logged immediately"

# Test 5: Verify all logs are valid JSON
echo "Test 5: All logs are valid JSON"
line_num=0
while IFS= read -r line; do
    line_num=$((line_num + 1))
    if ! echo "$line" | jq . > /dev/null 2>&1; then
        echo "FAIL: Line $line_num is not valid JSON"
        echo "Line content: $line"
        exit 1
    fi
done < "$LOGFILE"
echo "PASS: All $line_num log lines are valid JSON"

# Test 6: Verify session consistency
echo "Test 6: Session ID is consistent"
session_ids=$(jq -r '.session_id' < "$LOGFILE" | sort -u)
session_count=$(echo "$session_ids" | wc -l | tr -d ' ')
if [ "$session_count" -ne 1 ]; then
    echo "FAIL: Expected 1 unique session ID, found $session_count"
    echo "Session IDs: $session_ids"
    exit 1
fi
echo "PASS: Single consistent session ID across all logs"

echo ""
echo "=== All tests passed! ==="
echo "Total log lines: $line_num"
echo "Session ID: $session_ids"
