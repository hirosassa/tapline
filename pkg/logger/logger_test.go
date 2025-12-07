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

	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := &Logger{
		slogger:        slog.New(handler),
		Service:        "claude-code",
		SessionManager: sessionMgr,
		UserID:         "test-user-123",
		UserSource:     "env",
		Hostname:       "test-hostname",
		GitRepoURL:     "https://github.com/test/repo.git",
		GitRepoName:    "test/repo",
		GitBranch:      "main",
	}

	logger.LogUserPrompt(testSessionID, "Hello, Claude!")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	if result["service"] != "claude-code" {
		t.Errorf("Expected service 'claude-code', got %v", result["service"])
	}
	if result["session_id"] != testSessionID {
		t.Errorf("Expected session_id %s, got %v", testSessionID, result["session_id"])
	}
	if result["user_id"] != "test-user-123" {
		t.Errorf("Expected user_id 'test-user-123', got %v", result["user_id"])
	}
	if result["user_source"] != "env" {
		t.Errorf("Expected user_source 'env', got %v", result["user_source"])
	}
	if result["hostname"] != "test-hostname" {
		t.Errorf("Expected hostname 'test-hostname', got %v", result["hostname"])
	}
	if result["git_repo_url"] != "https://github.com/test/repo.git" {
		t.Errorf("Expected git_repo_url 'https://github.com/test/repo.git', got %v", result["git_repo_url"])
	}
	if result["git_repo_name"] != "test/repo" {
		t.Errorf("Expected git_repo_name 'test/repo', got %v", result["git_repo_name"])
	}
	if result["git_branch"] != "main" {
		t.Errorf("Expected git_branch 'main', got %v", result["git_branch"])
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
		UserID:         "test-user-456",
		UserSource:     "api_key_hash",
		Hostname:       "test-host",
	}

	logger.LogAssistantResponse(testSessionID, "How can I help?")

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	if result["user_id"] != "test-user-456" {
		t.Errorf("Expected user_id 'test-user-456', got %v", result["user_id"])
	}
	if result["user_source"] != "api_key_hash" {
		t.Errorf("Expected user_source 'api_key_hash', got %v", result["user_source"])
	}
	if result["hostname"] != "test-host" {
		t.Errorf("Expected hostname 'test-host', got %v", result["hostname"])
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
		UserID:         "test-user-789",
		UserSource:     "system",
		Hostname:       "test-machine",
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

	if result["user_id"] != "test-user-789" {
		t.Errorf("Expected user_id 'test-user-789', got %v", result["user_id"])
	}
	if result["user_source"] != "system" {
		t.Errorf("Expected user_source 'system', got %v", result["user_source"])
	}
	if result["hostname"] != "test-machine" {
		t.Errorf("Expected hostname 'test-machine', got %v", result["hostname"])
	}
	if result["role"] != "system" {
		t.Errorf("Expected role 'system', got %v", result["role"])
	}
	if result["event"] != "session_start" {
		t.Errorf("Expected event 'session_start', got %v", result["event"])
	}

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
		UserID:         "test-user-end",
		UserSource:     "anonymous",
		Hostname:       "end-host",
	}

	logger.LogSessionEnd(testSessionID)

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal log output: %v", err)
	}

	if result["user_id"] != "test-user-end" {
		t.Errorf("Expected user_id 'test-user-end', got %v", result["user_id"])
	}
	if result["user_source"] != "anonymous" {
		t.Errorf("Expected user_source 'anonymous', got %v", result["user_source"])
	}
	if result["hostname"] != "end-host" {
		t.Errorf("Expected hostname 'end-host', got %v", result["hostname"])
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
	if logger.UserID == "" {
		t.Error("UserID should not be empty")
	}
	if logger.UserSource == "" {
		t.Error("UserSource should not be empty")
	}
	if logger.Hostname == "" {
		t.Error("Hostname should not be empty")
	}

	t.Logf("Git info: URL=%s, Name=%s, Branch=%s",
		logger.GitRepoURL, logger.GitRepoName, logger.GitBranch)
}
