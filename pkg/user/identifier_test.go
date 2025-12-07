package user

import (
	"os"
	"testing"
)

func TestGetIdentifier_TaplineUserID(t *testing.T) {
	os.Setenv("TAPLINE_USER_ID", "test@example.com")
	defer os.Unsetenv("TAPLINE_USER_ID")

	result := GetIdentifier("claude-code")

	if result.UserID != "test@example.com" {
		t.Errorf("Expected UserID to be 'test@example.com', got '%s'", result.UserID)
	}
	if result.Source != "env" {
		t.Errorf("Expected Source to be 'env', got '%s'", result.Source)
	}
	if result.Hostname == "" {
		t.Error("Expected Hostname to be set")
	}
}

func TestGetIdentifier_APIKeyHash_Claude(t *testing.T) {
	os.Unsetenv("TAPLINE_USER_ID")
	os.Setenv("ANTHROPIC_API_KEY", "sk-ant-test-key-1234567890")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	result := GetIdentifier("claude-code")

	if result.UserID == "" {
		t.Error("Expected UserID to be set from API key hash")
	}
	if result.Source != "api_key_hash" {
		t.Errorf("Expected Source to be 'api_key_hash', got '%s'", result.Source)
	}
	if len(result.UserID) != 16 {
		t.Errorf("Expected UserID to be 16 characters, got %d", len(result.UserID))
	}
}

func TestGetIdentifier_APIKeyHash_Gemini(t *testing.T) {
	os.Setenv("GEMINI_API_KEY", "test-gemini-key-1234567890")
	defer func() {
		os.Unsetenv("TAPLINE_USER_ID")
		os.Unsetenv("ANTHROPIC_API_KEY")
		os.Unsetenv("GEMINI_API_KEY")
	}()

	result := GetIdentifier("gemini-cli")

	if result.UserID == "" {
		t.Error("Expected UserID to be set from API key hash")
	}
	if result.Source != "api_key_hash" {
		t.Errorf("Expected Source to be 'api_key_hash', got '%s'", result.Source)
	}
}

func TestGetIdentifier_APIKeyHash_GeminiGoogleKey(t *testing.T) {
	os.Unsetenv("TAPLINE_USER_ID")
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("GEMINI_API_KEY")
	os.Setenv("GOOGLE_API_KEY", "test-google-key-1234567890")
	defer os.Unsetenv("GOOGLE_API_KEY")

	result := GetIdentifier("gemini-cli")

	if result.UserID == "" {
		t.Error("Expected UserID to be set from API key hash")
	}
	if result.Source != "api_key_hash" {
		t.Errorf("Expected Source to be 'api_key_hash', got '%s'", result.Source)
	}
}

func TestGetIdentifier_SystemUser(t *testing.T) {
	os.Unsetenv("TAPLINE_USER_ID")
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("GOOGLE_API_KEY")
	os.Unsetenv("OPENAI_API_KEY")

	originalUser := os.Getenv("USER")
	os.Setenv("USER", "testuser")
	defer func() {
		if originalUser != "" {
			os.Setenv("USER", originalUser)
		} else {
			os.Unsetenv("USER")
		}
	}()

	result := GetIdentifier("claude-code")

	if result.UserID != "testuser" {
		t.Errorf("Expected UserID to be 'testuser', got '%s'", result.UserID)
	}
	if result.Source != "system" {
		t.Errorf("Expected Source to be 'system', got '%s'", result.Source)
	}
}

func TestGetIdentifier_Anonymous(t *testing.T) {
	os.Unsetenv("TAPLINE_USER_ID")
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("GOOGLE_API_KEY")
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("USER")

	result := GetIdentifier("claude-code")

	if result.UserID != "anonymous" {
		t.Errorf("Expected UserID to be 'anonymous', got '%s'", result.UserID)
	}
	if result.Source != "anonymous" {
		t.Errorf("Expected Source to be 'anonymous', got '%s'", result.Source)
	}
}

func TestGetIdentifier_Priority(t *testing.T) {
	os.Setenv("TAPLINE_USER_ID", "priority@example.com")
	os.Setenv("ANTHROPIC_API_KEY", "sk-ant-test-key")
	os.Setenv("USER", "testuser")
	defer func() {
		os.Unsetenv("TAPLINE_USER_ID")
		os.Unsetenv("ANTHROPIC_API_KEY")
		os.Unsetenv("USER")
	}()

	result := GetIdentifier("claude-code")

	if result.UserID != "priority@example.com" {
		t.Errorf("Expected TAPLINE_USER_ID to take priority, got '%s'", result.UserID)
	}
	if result.Source != "env" {
		t.Errorf("Expected Source to be 'env', got '%s'", result.Source)
	}
}

func TestGetAPIKeyHash_ShortKey(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "short")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	hash := getAPIKeyHash("claude-code")

	if hash == "" {
		t.Error("Expected hash to be generated for short key")
	}
	if len(hash) != 16 {
		t.Errorf("Expected hash to be 16 characters, got %d", len(hash))
	}
}

func TestGetAPIKeyHash_LongKey(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "this-is-a-very-long-api-key-that-exceeds-sixteen-characters")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	hash := getAPIKeyHash("claude-code")

	if hash == "" {
		t.Error("Expected hash to be generated for long key")
	}
	if len(hash) != 16 {
		t.Errorf("Expected hash to be 16 characters, got %d", len(hash))
	}
}

func TestGetAPIKeyHash_DifferentKeys(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "key-1234567890")
	hash1 := getAPIKeyHash("claude-code")
	os.Unsetenv("ANTHROPIC_API_KEY")

	os.Setenv("ANTHROPIC_API_KEY", "key-0987654321")
	hash2 := getAPIKeyHash("claude-code")
	os.Unsetenv("ANTHROPIC_API_KEY")

	if hash1 == hash2 {
		t.Error("Expected different keys to produce different hashes")
	}
}

func TestGetAPIKeyHash_SameKeyDifferentService(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "sk-ant-test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	hash1 := getAPIKeyHash("claude-code")
	hash2 := getAPIKeyHash("claude-api")

	if hash1 != hash2 {
		t.Error("Expected same key to produce same hash for different Claude services")
	}
}

func TestGetAPIKeyHash_UnknownService(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "sk-ant-test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	hash := getAPIKeyHash("unknown-service")

	if hash != "" {
		t.Errorf("Expected empty hash for unknown service, got '%s'", hash)
	}
}

func TestGetAPIKeyHash_NoAPIKey(t *testing.T) {
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("GOOGLE_API_KEY")
	os.Unsetenv("OPENAI_API_KEY")

	hash := getAPIKeyHash("claude-code")

	if hash != "" {
		t.Errorf("Expected empty hash when no API key is set, got '%s'", hash)
	}
}
