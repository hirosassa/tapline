package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hirosassa/tapline/pkg/logger"
	"github.com/hirosassa/tapline/pkg/session"
)

func handleConversationStart() {
	sessionMgr, err := session.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize session manager: %v\n", err)
		os.Exit(1)
	}

	log := logger.NewLogger("claude-code", sessionMgr)

	sessionID := uuid.New().String()
	if err := sessionMgr.SetSessionID(sessionID); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set session ID: %v\n", err)
		os.Exit(1)
	}

	metadata := map[string]string{
		"hostname": getHostname(),
		"cwd":      getCwd(),
	}

	log.LogSessionStart(sessionID, metadata)
}

func handleConversationEnd() {
	sessionMgr, err := session.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize session manager: %v\n", err)
		os.Exit(1)
	}

	log := logger.NewLogger("claude-code", sessionMgr)

	sessionID, err := sessionMgr.GetSessionID()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get session ID: %v\n", err)
		os.Exit(1)
	}

	log.LogSessionEnd(sessionID)

	if err := sessionMgr.ClearSession(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to clear session: %v\n", err)
	}
}

func handleUserPrompt(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "user_prompt requires prompt argument")
		os.Exit(1)
	}

	sessionMgr, err := session.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize session manager: %v\n", err)
		os.Exit(1)
	}

	log := logger.NewLogger("claude-code", sessionMgr)

	sessionID, err := sessionMgr.GetSessionID()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get session ID: %v\n", err)
		os.Exit(1)
	}

	prompt := args[0]
	if len(args) > 1 {
		prompt = ""
		for i, arg := range args {
			if i > 0 {
				prompt += " "
			}
			prompt += arg
		}
	}

	log.LogUserPrompt(sessionID, prompt)
}

func handleAssistantResponse(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "assistant_response requires response argument")
		os.Exit(1)
	}

	sessionMgr, err := session.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize session manager: %v\n", err)
		os.Exit(1)
	}

	log := logger.NewLogger("claude-code", sessionMgr)

	sessionID, err := sessionMgr.GetSessionID()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get session ID: %v\n", err)
		os.Exit(1)
	}

	response := args[0]
	if len(args) > 1 {
		response = ""
		for i, arg := range args {
			if i > 0 {
				response += " "
			}
			response += arg
		}
	}

	log.LogAssistantResponse(sessionID, response)
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func getCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return cwd
}
