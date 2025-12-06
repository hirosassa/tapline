package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/hirosassa/tapline/pkg/session"
)

func TestLogger_LogUserPrompt(t *testing.T) {
	sessionMgr, err := session.NewManager()
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	testSessionID := "test-session-123"
	if err := sessionMgr.SetSessionID(testSessionID); err != nil {
		t.Fatalf("Failed to set session ID: %v", err)
	}
	defer sessionMgr.ClearSession()

	// Capture output
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := &Logger{
		slogger:        slog.New(handler),
		Service:        "claude-code",
		SessionManager: sessionMgr,
	}

	// Log user prompt
	logger.LogUserPrompt(testSessionID, "Hello, Claude!")

	// Parse output
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	// Verify fields
	if result["service"] != "claude-code" {
		t.Errorf("Expected service 'claude-code', got %v", result["service"])
	}
	if result["session_id"] != testSessionID {
		t.Errorf("Expected session_id %s, got %v", testSessionID, result["session_id"])
	}
	if result["role"] != "user" {
		t.Errorf("Expected role 'user', got %v", result["role"])
	}
	if result["content"] != "Hello, Claude!" {
		t.Errorf("Expected content 'Hello, Claude!', got %v", result["content"])
	}
}

func TestLogger_LogAssistantResponse(t *testing.T) {
	sessionMgr, err := session.NewManager()
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	testSessionID := "test-session-456"
	if err := sessionMgr.SetSessionID(testSessionID); err != nil {
		t.Fatalf("Failed to set session ID: %v", err)
	}
	defer sessionMgr.ClearSession()

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := &Logger{
		slogger:        slog.New(handler),
		Service:        "claude-code",
		SessionManager: sessionMgr,
	}

	logger.LogAssistantResponse(testSessionID, "How can I help?")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	if result["role"] != "assistant" {
		t.Errorf("Expected role 'assistant', got %v", result["role"])
	}
	if result["content"] != "How can I help?" {
		t.Errorf("Expected content 'How can I help?', got %v", result["content"])
	}
}

func TestLogger_LogSessionStart(t *testing.T) {
	sessionMgr, err := session.NewManager()
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	testSessionID := "test-session-789"

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := &Logger{
		slogger:        slog.New(handler),
		Service:        "claude-code",
		SessionManager: sessionMgr,
	}

	metadata := map[string]string{
		"hostname": "test-host",
		"cwd":      "/test/path",
	}

	logger.LogSessionStart(testSessionID, metadata)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	if result["role"] != "system" {
		t.Errorf("Expected role 'system', got %v", result["role"])
	}
	if result["event"] != "session_start" {
		t.Errorf("Expected event 'session_start', got %v", result["event"])
	}

	// Check metadata
	metadataResult, ok := result["metadata"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected metadata to be a map")
	}
	if metadataResult["hostname"] != "test-host" {
		t.Errorf("Expected hostname 'test-host', got %v", metadataResult["hostname"])
	}
}

func TestLogger_LogSessionEnd(t *testing.T) {
	sessionMgr, err := session.NewManager()
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	testSessionID := "test-session-end"

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := &Logger{
		slogger:        slog.New(handler),
		Service:        "claude-code",
		SessionManager: sessionMgr,
	}

	logger.LogSessionEnd(testSessionID)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	if result["event"] != "session_end" {
		t.Errorf("Expected event 'session_end', got %v", result["event"])
	}
	if result["role"] != "system" {
		t.Errorf("Expected role 'system', got %v", result["role"])
	}
}

func TestNewLogger(t *testing.T) {
	sessionMgr, err := session.NewManager()
	if err != nil {
		t.Fatalf("Failed to create session manager: %v", err)
	}

	logger := NewLogger("test-service", sessionMgr)

	if logger.Service != "test-service" {
		t.Errorf("Expected service 'test-service', got %s", logger.Service)
	}
	if logger.SessionManager != sessionMgr {
		t.Error("SessionManager not properly set")
	}
	if logger.slogger == nil {
		t.Error("slogger should not be nil")
	}
}
