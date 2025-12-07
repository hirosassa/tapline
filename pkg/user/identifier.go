// Package user provides user identification functionality for logging.
// It supports multiple identification sources including explicit environment variables,
// API key hashes, and system user information.
package user

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

// Identifier contains user identification information and its source.
type Identifier struct {
	UserID   string `json:"user_id"`
	Source   string `json:"source"`
	Hostname string `json:"hostname,omitempty"`
}

// GetIdentifier returns user identification information based on priority:
// 1. TAPLINE_USER_ID environment variable
// 2. API key hash for the service
// 3. System USER environment variable
// 4. "anonymous" as fallback
func GetIdentifier(service string) Identifier {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	if userID := os.Getenv("TAPLINE_USER_ID"); userID != "" {
		return Identifier{
			UserID:   userID,
			Source:   "env",
			Hostname: hostname,
		}
	}

	if apiKeyHash := getAPIKeyHash(service); apiKeyHash != "" {
		return Identifier{
			UserID:   apiKeyHash,
			Source:   "api_key_hash",
			Hostname: hostname,
		}
	}

	if user := os.Getenv("USER"); user != "" {
		return Identifier{
			UserID:   user,
			Source:   "system",
			Hostname: hostname,
		}
	}

	return Identifier{
		UserID:   "anonymous",
		Source:   "anonymous",
		Hostname: hostname,
	}
}

func getAPIKeyHash(service string) string {
	var apiKey string

	switch service {
	case "claude-code", "claude-api":
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	case "gemini-cli":
		if key := os.Getenv("GEMINI_API_KEY"); key != "" {
			apiKey = key
		} else {
			apiKey = os.Getenv("GOOGLE_API_KEY")
		}
	case "chatgpt", "openai":
		apiKey = os.Getenv("OPENAI_API_KEY")
	case "codex-cli":
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	if apiKey == "" {
		return ""
	}

	// Always hash the full API key for consistent security properties.
	hash := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hash[:8])
}
