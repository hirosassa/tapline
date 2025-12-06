#!/bin/bash
# Example jq queries for analyzing Tapline logs

LOGFILE="sample_output.jsonl"

echo "=== All logs (pretty printed) ==="
cat $LOGFILE | jq .
echo ""

echo "=== Filter by session ID ==="
cat $LOGFILE | jq 'select(.session_id=="0e97d08c-b08b-4f5a-92ea-086c36d5818b")'
echo ""

echo "=== Extract only user prompts ==="
cat $LOGFILE | jq 'select(.role=="user") | .content'
echo ""

echo "=== Extract only assistant responses ==="
cat $LOGFILE | jq 'select(.role=="assistant") | .content'
echo ""

echo "=== Get conversation timeline ==="
cat $LOGFILE | jq '{timestamp, role, content: (.content | split("\n")[0] + "...")}'
echo ""

echo "=== Count messages by role ==="
cat $LOGFILE | jq -s 'group_by(.role) | map({role: .[0].role, count: length})'
echo ""

echo "=== Extract session metadata ==="
cat $LOGFILE | jq 'select(.event=="session_start") | .metadata'
echo ""

echo "=== Find all session events ==="
cat $LOGFILE | jq 'select(.event != null) | {timestamp, event, session_id}'
echo ""
