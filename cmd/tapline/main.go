package main

import (
	"fmt"
	"os"
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

	switch command {
	case "conversation_start":
		handleConversationStart()
	case "conversation_end":
		handleConversationEnd()
	case "user_prompt":
		handleUserPrompt(os.Args[2:])
	case "assistant_response":
		handleAssistantResponse(os.Args[2:])
	case "wrap-gemini":
		wrapGemini(os.Args[2:])
	case "notify-codex":
		notifyCodex()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
