package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNotifyCodex_NoStdin(t *testing.T) {
	if os.Getenv("TEST_NOTIFY_CODEX_NO_STDIN") == "1" {
		notifyCodex()
		return
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestNotifyCodex_NoStdin")
	cmd.Env = append(os.Environ(), "TEST_NOTIFY_CODEX_NO_STDIN=1")
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected clean exit, got error: %v", err)
	}
}

func TestNotifyCodex_InvalidJSON(t *testing.T) {
	if os.Getenv("TEST_NOTIFY_CODEX_INVALID_JSON") == "1" {
		notifyCodex()
		return
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestNotifyCodex_InvalidJSON")
	cmd.Env = append(os.Environ(), "TEST_NOTIFY_CODEX_INVALID_JSON=1")
	cmd.Stdin = strings.NewReader("invalid json")
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected clean exit even with invalid JSON, got error: %v", err)
	}
}

func TestNotifyCodex_UnknownEventType(t *testing.T) {
	if os.Getenv("TEST_NOTIFY_CODEX_UNKNOWN_EVENT") == "1" {
		notifyCodex()
		return
	}

	event := CodexEvent{
		Type: "unknown-event-type",
	}
	eventJSON, _ := json.Marshal(event)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestNotifyCodex_UnknownEventType")
	cmd.Env = append(os.Environ(), "TEST_NOTIFY_CODEX_UNKNOWN_EVENT=1")
	cmd.Stdin = bytes.NewReader(eventJSON)
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected clean exit with unknown event, got error: %v", err)
	}
}

func TestNotifyCodex_AgentTurnComplete_Format1(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	if os.Getenv("TEST_NOTIFY_CODEX_FORMAT1") == "1" {
		notifyCodex()
		return
	}

	sessionDir := filepath.Join(tmpDir, ".tapline")
	os.MkdirAll(sessionDir, 0o750)
	sessionFile := filepath.Join(sessionDir, "session_id")
	sessionID := "test-session-123"
	os.WriteFile(sessionFile, []byte(sessionID), 0o600)

	eventData := AgentTurnCompleteData{
		Response: "Test response from Codex",
	}
	event := CodexEvent{
		Type: "agent-turn-complete",
	}
	event.Data, _ = json.Marshal(eventData)
	eventJSON, _ := json.Marshal(event)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestNotifyCodex_AgentTurnComplete_Format1")
	cmd.Env = append(os.Environ(), "TEST_NOTIFY_CODEX_FORMAT1=1", "HOME="+tmpDir)
	cmd.Stdin = bytes.NewReader(eventJSON)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected success, got error: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "Test response from Codex") {
		t.Logf("Output: %s", output)
	}
}

func TestNotifyCodex_AgentTurnComplete_Format2(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	if os.Getenv("TEST_NOTIFY_CODEX_FORMAT2") == "1" {
		notifyCodex()
		return
	}

	sessionDir := filepath.Join(tmpDir, ".tapline")
	os.MkdirAll(sessionDir, 0o750)
	sessionFile := filepath.Join(sessionDir, "session_id")
	sessionID := "test-session-456"
	os.WriteFile(sessionFile, []byte(sessionID), 0o600)

	nestedData := struct {
		Data AgentTurnCompleteData `json:"data"`
	}{
		Data: AgentTurnCompleteData{
			Response: "Nested test response",
		},
	}
	event := CodexEvent{
		Type: "agent-turn-complete",
	}
	event.Data, _ = json.Marshal(nestedData)
	eventJSON, _ := json.Marshal(event)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestNotifyCodex_AgentTurnComplete_Format2")
	cmd.Env = append(os.Environ(), "TEST_NOTIFY_CODEX_FORMAT2=1", "HOME="+tmpDir)
	cmd.Stdin = bytes.NewReader(eventJSON)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected success, got error: %v\nOutput: %s", err, output)
	}
}

func TestNotifyCodex_SessionStart(t *testing.T) {
	tmpDir := t.TempDir()

	if os.Getenv("TEST_NOTIFY_CODEX_SESSION_START") == "1" {
		notifyCodex()
		return
	}

	event := CodexEvent{
		Type: "session_start",
	}
	eventJSON, _ := json.Marshal(event)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestNotifyCodex_SessionStart")
	cmd.Env = append(os.Environ(), "TEST_NOTIFY_CODEX_SESSION_START=1", "HOME="+tmpDir)
	cmd.Stdin = bytes.NewReader(eventJSON)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected success, got error: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "session_start") {
		t.Logf("Output: %s", output)
	}
}

func TestNotifyCodex_SessionEnd(t *testing.T) {
	tmpDir := t.TempDir()

	if os.Getenv("TEST_NOTIFY_CODEX_SESSION_END") == "1" {
		notifyCodex()
		return
	}

	sessionDir := filepath.Join(tmpDir, ".tapline")
	os.MkdirAll(sessionDir, 0o750)
	sessionFile := filepath.Join(sessionDir, "session_id")
	os.WriteFile(sessionFile, []byte("test-session-789"), 0o600)

	event := CodexEvent{
		Type: "session_end",
	}
	eventJSON, _ := json.Marshal(event)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestNotifyCodex_SessionEnd")
	cmd.Env = append(os.Environ(), "TEST_NOTIFY_CODEX_SESSION_END=1", "HOME="+tmpDir)
	cmd.Stdin = bytes.NewReader(eventJSON)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected success, got error: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "session_end") {
		t.Logf("Output: %s", output)
	}
}

func TestIsTerminal(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if isTerminal(f) {
		t.Error("Expected regular file to not be terminal")
	}
}

func TestHandleAgentTurnComplete_EmptyResponse(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	if os.Getenv("TEST_HANDLE_EMPTY_RESPONSE") == "1" {
		sessionDir := filepath.Join(tmpDir, ".tapline")
		os.MkdirAll(sessionDir, 0o750)

		notifyCodex()
		return
	}

	eventData := AgentTurnCompleteData{
		Response: "",
	}
	event := CodexEvent{
		Type: "agent-turn-complete",
	}
	event.Data, _ = json.Marshal(eventData)
	eventJSON, _ := json.Marshal(event)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHandleAgentTurnComplete_EmptyResponse")
	cmd.Env = append(os.Environ(), "TEST_HANDLE_EMPTY_RESPONSE=1", "HOME="+tmpDir)
	cmd.Stdin = bytes.NewReader(eventJSON)
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected clean exit with empty response, got error: %v", err)
	}
}

func TestHandleSessionEnd_NoActiveSession(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	if os.Getenv("TEST_SESSION_END_NO_ACTIVE") == "1" {
		notifyCodex()
		return
	}

	event := CodexEvent{
		Type: "session_end",
	}
	eventJSON, _ := json.Marshal(event)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHandleSessionEnd_NoActiveSession")
	cmd.Env = append(os.Environ(), "TEST_SESSION_END_NO_ACTIVE=1", "HOME="+tmpDir)
	cmd.Stdin = bytes.NewReader(eventJSON)
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected clean exit when no session exists, got error: %v", err)
	}
}

func TestCodexEvent_UnmarshalJSON(t *testing.T) {
	jsonStr := `{"type":"test-event","data":{"key":"value"}}`
	var event CodexEvent
	err := json.Unmarshal([]byte(jsonStr), &event)
	if err != nil {
		t.Errorf("Failed to unmarshal CodexEvent: %v", err)
	}
	if event.Type != "test-event" {
		t.Errorf("Expected type 'test-event', got '%s'", event.Type)
	}
}

func TestAgentTurnCompleteData_UnmarshalJSON(t *testing.T) {
	jsonStr := `{"response":"test response"}`
	var data AgentTurnCompleteData
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		t.Errorf("Failed to unmarshal AgentTurnCompleteData: %v", err)
	}
	if data.Response != "test response" {
		t.Errorf("Expected response 'test response', got '%s'", data.Response)
	}
}

func TestNotifyCodex_ReadError(t *testing.T) {
	if os.Getenv("TEST_NOTIFY_CODEX_READ_ERROR") == "1" {
		notifyCodex()
		return
	}

	r, w := io.Pipe()
	w.Close()

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestNotifyCodex_ReadError")
	cmd.Env = append(os.Environ(), "TEST_NOTIFY_CODEX_READ_ERROR=1")
	cmd.Stdin = r
	err := cmd.Run()
	if err != nil {
		t.Errorf("Expected clean exit even with read error, got: %v", err)
	}
}
