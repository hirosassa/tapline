package session

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	sessionDirName  = ".tapline"
	sessionFileName = "session_id"
)

// Manager handles session ID persistence
type Manager struct {
	sessionDir  string
	sessionFile string
}

// NewManager creates a new session manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	sessionDir := filepath.Join(homeDir, sessionDirName)
	sessionFile := filepath.Join(sessionDir, sessionFileName)

	// Ensure session directory exists
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}

	return &Manager{
		sessionDir:  sessionDir,
		sessionFile: sessionFile,
	}, nil
}

// GetSessionID retrieves the current session ID
func (m *Manager) GetSessionID() (string, error) {
	data, err := os.ReadFile(m.sessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no active session found")
		}
		return "", fmt.Errorf("failed to read session file: %w", err)
	}

	sessionID := string(data)
	if sessionID == "" {
		return "", fmt.Errorf("session file is empty")
	}

	return sessionID, nil
}

// SetSessionID stores a new session ID
func (m *Manager) SetSessionID(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	if err := os.WriteFile(m.sessionFile, []byte(sessionID), 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// ClearSession removes the current session ID
func (m *Manager) ClearSession() error {
	if err := os.Remove(m.sessionFile); err != nil {
		if os.IsNotExist(err) {
			return nil // Already cleared
		}
		return fmt.Errorf("failed to remove session file: %w", err)
	}

	return nil
}

// HasActiveSession checks if there's an active session
func (m *Manager) HasActiveSession() bool {
	_, err := m.GetSessionID()
	return err == nil
}
