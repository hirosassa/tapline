package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/hirosassa/tapline/pkg/logger"
	"github.com/hirosassa/tapline/pkg/session"
)

type CodexEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

type AgentTurnCompleteData struct {
	Response string `json:"response"`
}

func notifyCodex() {
	if isTerminal(os.Stdin) {
		os.Exit(0)
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		os.Exit(0)
	}

	var event CodexEvent
	if err := json.Unmarshal(data, &event); err != nil {
		os.Exit(0)
	}

	sessionMgr, err := session.NewManager()
	if err != nil {
		os.Exit(0)
	}

	log := logger.NewLogger("codex-cli", sessionMgr)

	switch event.Type {
	case "agent-turn-complete":
		handleAgentTurnComplete(log, sessionMgr, event.Data)
	case "session_start":
		handleSessionStart(log, sessionMgr)
	case "session_end":
		handleSessionEnd(log, sessionMgr)
	default:
		// Unknown event type, ignore
	}

	os.Exit(0)
}

func handleAgentTurnComplete(log *logger.Logger, sessionMgr *session.Manager, data json.RawMessage) {
	// Try multiple possible data structures
	// Format 1: {response: "..."}
	var eventData1 AgentTurnCompleteData
	if err := json.Unmarshal(data, &eventData1); err == nil && eventData1.Response != "" {
		logResponse(log, sessionMgr, eventData1.Response)
		return
	}

	// Format 2: {data: {response: "..."}}
	var eventData2 struct {
		Data AgentTurnCompleteData `json:"data"`
	}
	if err := json.Unmarshal(data, &eventData2); err == nil && eventData2.Data.Response != "" {
		logResponse(log, sessionMgr, eventData2.Data.Response)
		return
	}
}

func logResponse(log *logger.Logger, sessionMgr *session.Manager, response string) {
	if !sessionMgr.HasActiveSession() {
		newSessionID := uuid.New().String()
		if err := sessionMgr.SetSessionID(newSessionID); err != nil {
			return
		}
	}

	sessionID, err := sessionMgr.GetSessionID()
	if err != nil {
		return
	}

	log.LogAssistantResponse(sessionID, response)
}

func handleSessionStart(log *logger.Logger, sessionMgr *session.Manager) {
	newSessionID := uuid.New().String()
	if err := sessionMgr.SetSessionID(newSessionID); err != nil {
		return
	}
	log.LogSessionStart(newSessionID, nil)
}

func handleSessionEnd(log *logger.Logger, sessionMgr *session.Manager) {
	if !sessionMgr.HasActiveSession() {
		return
	}

	sessionID, err := sessionMgr.GetSessionID()
	if err != nil {
		return
	}

	log.LogSessionEnd(sessionID)
	if err := sessionMgr.ClearSession(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to clear session: %v\n", err)
	}
}

func isTerminal(f *os.File) bool {
	fileInfo, err := f.Stat()
	if err != nil {
		return true
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
