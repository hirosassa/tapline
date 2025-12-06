package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleConversationStart(t *testing.T) {
	tmpDir := t.TempDir()

	if os.Getenv("TEST_CONVERSATION_START") == "1" {
		handleConversationStart()
		return
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHandleConversationStart")
	cmd.Env = append(os.Environ(), "TEST_CONVERSATION_START=1", "HOME="+tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected success, got error: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "session_start") {
		t.Errorf("Expected session_start in output, got: %s", output)
	}

	if !strings.Contains(string(output), "claude-code") {
		t.Errorf("Expected claude-code service in output, got: %s", output)
	}
}

func TestHandleConversationEnd(t *testing.T) {
	tmpDir := t.TempDir()

	if os.Getenv("TEST_CONVERSATION_END") == "1" {
		handleConversationEnd()
		return
	}

	sessionDir := filepath.Join(tmpDir, ".tapline")
	os.MkdirAll(sessionDir, 0o750)
	sessionFile := filepath.Join(sessionDir, "session_id")
	os.WriteFile(sessionFile, []byte("test-session-id"), 0o600)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHandleConversationEnd")
	cmd.Env = append(os.Environ(), "TEST_CONVERSATION_END=1", "HOME="+tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected success, got error: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "session_end") {
		t.Errorf("Expected session_end in output, got: %s", output)
	}
}

func TestHandleUserPrompt(t *testing.T) {
	tmpDir := t.TempDir()

	if os.Getenv("TEST_USER_PROMPT") == "1" {
		handleUserPrompt([]string{"test", "prompt", "text"})
		return
	}

	sessionDir := filepath.Join(tmpDir, ".tapline")
	os.MkdirAll(sessionDir, 0o750)
	sessionFile := filepath.Join(sessionDir, "session_id")
	os.WriteFile(sessionFile, []byte("test-session-id"), 0o600)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHandleUserPrompt")
	cmd.Env = append(os.Environ(), "TEST_USER_PROMPT=1", "HOME="+tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected success, got error: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "test prompt text") {
		t.Errorf("Expected prompt text in output, got: %s", output)
	}

	if !strings.Contains(string(output), `"role":"user"`) {
		t.Errorf("Expected user role in output, got: %s", output)
	}
}

func TestHandleUserPrompt_NoArgs(t *testing.T) {
	if os.Getenv("TEST_USER_PROMPT_NO_ARGS") == "1" {
		handleUserPrompt([]string{})
		return
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHandleUserPrompt_NoArgs")
	cmd.Env = append(os.Environ(), "TEST_USER_PROMPT_NO_ARGS=1")
	err := cmd.Run()

	if err == nil {
		t.Error("Expected command to exit with error when no args provided")
	}
}

func TestHandleAssistantResponse(t *testing.T) {
	tmpDir := t.TempDir()

	if os.Getenv("TEST_ASSISTANT_RESPONSE") == "1" {
		handleAssistantResponse([]string{"test", "response", "text"})
		return
	}

	sessionDir := filepath.Join(tmpDir, ".tapline")
	os.MkdirAll(sessionDir, 0o750)
	sessionFile := filepath.Join(sessionDir, "session_id")
	os.WriteFile(sessionFile, []byte("test-session-id"), 0o600)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHandleAssistantResponse")
	cmd.Env = append(os.Environ(), "TEST_ASSISTANT_RESPONSE=1", "HOME="+tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Expected success, got error: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(string(output), "test response text") {
		t.Errorf("Expected response text in output, got: %s", output)
	}

	if !strings.Contains(string(output), `"role":"assistant"`) {
		t.Errorf("Expected assistant role in output, got: %s", output)
	}
}

func TestHandleAssistantResponse_NoArgs(t *testing.T) {
	if os.Getenv("TEST_ASSISTANT_RESPONSE_NO_ARGS") == "1" {
		handleAssistantResponse([]string{})
		return
	}

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestHandleAssistantResponse_NoArgs")
	cmd.Env = append(os.Environ(), "TEST_ASSISTANT_RESPONSE_NO_ARGS=1")
	err := cmd.Run()

	if err == nil {
		t.Error("Expected command to exit with error when no args provided")
	}
}

func TestGetHostname(t *testing.T) {
	hostname := getHostname()
	if hostname == "" {
		t.Error("Expected non-empty hostname")
	}
}

func TestGetCwd(t *testing.T) {
	cwd := getCwd()
	if cwd == "" {
		t.Error("Expected non-empty cwd")
	}
}
