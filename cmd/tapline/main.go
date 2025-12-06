package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hirosassa/tapline/pkg/logger"
	"github.com/hirosassa/tapline/pkg/session"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: tapline <command> [args...]")
		os.Exit(1)
	}

	command := os.Args[1]

	if command == "version" {
		fmt.Printf("tapline %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	sessionMgr, err := session.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize session manager: %v\n", err)
		os.Exit(1)
	}

	log := logger.NewLogger("claude-code", sessionMgr)

	switch command {
	case "conversation_start":
		handleConversationStart(log, sessionMgr)
	case "conversation_end":
		handleConversationEnd(log, sessionMgr)
	case "user_prompt":
		handleUserPrompt(log, os.Args[2:])
	case "assistant_response":
		handleAssistantResponse(log, os.Args[2:])
	case "wrap-gemini":
		wrapGemini(os.Args[2:])
	case "notify-codex":
		notifyCodex()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func handleConversationStart(log *logger.Logger, sessionMgr *session.Manager) {
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

func handleConversationEnd(log *logger.Logger, sessionMgr *session.Manager) {
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

func handleUserPrompt(log *logger.Logger, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "user_prompt requires prompt argument")
		os.Exit(1)
	}

	sessionID, err := log.SessionManager.GetSessionID()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get session ID: %v\n", err)
		os.Exit(1)
	}

	// Join all arguments to handle spaces in prompt
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

func handleAssistantResponse(log *logger.Logger, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "assistant_response requires response argument")
		os.Exit(1)
	}

	sessionID, err := log.SessionManager.GetSessionID()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get session ID: %v\n", err)
		os.Exit(1)
	}

	// Join all arguments to handle spaces in response
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
