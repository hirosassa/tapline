package session

import (
	"testing"

	"github.com/google/uuid"
)

func TestManager_SessionLifecycle(t *testing.T) {
	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Clean up any existing session
	mgr.ClearSession()

	// Test: No active session initially
	if mgr.HasActiveSession() {
		t.Error("Expected no active session initially")
	}

	// Test: Set session ID
	testSessionID := uuid.New().String()
	if err := mgr.SetSessionID(testSessionID); err != nil {
		t.Fatalf("Failed to set session ID: %v", err)
	}

	// Test: Has active session after setting
	if !mgr.HasActiveSession() {
		t.Error("Expected active session after SetSessionID")
	}

	// Test: Get session ID
	retrievedID, err := mgr.GetSessionID()
	if err != nil {
		t.Fatalf("Failed to get session ID: %v", err)
	}
	if retrievedID != testSessionID {
		t.Errorf("Session ID mismatch: got %s, want %s", retrievedID, testSessionID)
	}

	// Test: Clear session
	if err := mgr.ClearSession(); err != nil {
		t.Fatalf("Failed to clear session: %v", err)
	}

	// Test: No active session after clearing
	if mgr.HasActiveSession() {
		t.Error("Expected no active session after ClearSession")
	}
}

func TestManager_SetEmptySessionID(t *testing.T) {
	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	err = mgr.SetSessionID("")
	if err == nil {
		t.Error("Expected error when setting empty session ID")
	}
}

func TestManager_GetSessionIDWithoutSetting(t *testing.T) {
	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Clean up any existing session
	mgr.ClearSession()

	_, err = mgr.GetSessionID()
	if err == nil {
		t.Error("Expected error when getting session ID without setting it first")
	}
}

func TestManager_ConcurrentSessions(t *testing.T) {
	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	defer mgr.ClearSession()

	// Set first session
	session1 := uuid.New().String()
	if err := mgr.SetSessionID(session1); err != nil {
		t.Fatalf("Failed to set first session: %v", err)
	}

	// Verify first session
	retrieved, err := mgr.GetSessionID()
	if err != nil {
		t.Fatalf("Failed to get first session: %v", err)
	}
	if retrieved != session1 {
		t.Errorf("First session mismatch: got %s, want %s", retrieved, session1)
	}

	// Overwrite with second session
	session2 := uuid.New().String()
	if err := mgr.SetSessionID(session2); err != nil {
		t.Fatalf("Failed to set second session: %v", err)
	}

	// Verify second session overwrote the first
	retrieved, err = mgr.GetSessionID()
	if err != nil {
		t.Fatalf("Failed to get second session: %v", err)
	}
	if retrieved != session2 {
		t.Errorf("Second session mismatch: got %s, want %s", retrieved, session2)
	}
	if retrieved == session1 {
		t.Error("Second session should have overwritten the first")
	}
}
