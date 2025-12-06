// Package logger provides conversation logging functionality using structured logging.
// It supports logging user prompts, assistant responses, and session lifecycle events
// in JSON format to stdout with immediate flushing for durability.
package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/hirosassa/tapline/pkg/session"
)

// Logger handles conversation logging using slog
type Logger struct {
	slogger        *slog.Logger
	SessionManager *session.Manager
	Service        string
}

// NewLogger creates a new Logger instance with slog JSON handler
func NewLogger(service string, sessionMgr *session.Manager) *Logger {
	// Create JSON handler that writes to stdout
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return &Logger{
		slogger:        slog.New(handler),
		Service:        service,
		SessionManager: sessionMgr,
	}
}

// LogUserPrompt logs a user prompt
func (l *Logger) LogUserPrompt(sessionID, content string) {
	l.slogger.Info("conversation",
		slog.String("service", l.Service),
		slog.String("session_id", sessionID),
		slog.String("role", "user"),
		slog.String("content", content),
	)
	//nolint:errcheck // Sync errors are not critical for logging
	os.Stdout.Sync()
}

// LogAssistantResponse logs an assistant response
func (l *Logger) LogAssistantResponse(sessionID, content string) {
	l.slogger.Info("conversation",
		slog.String("service", l.Service),
		slog.String("session_id", sessionID),
		slog.String("role", "assistant"),
		slog.String("content", content),
	)
	//nolint:errcheck // Sync errors are not critical for logging
	os.Stdout.Sync()
}

// LogSessionStart logs a session start event
func (l *Logger) LogSessionStart(sessionID string, metadata map[string]string) {
	attrs := []slog.Attr{
		slog.String("service", l.Service),
		slog.String("session_id", sessionID),
		slog.String("role", "system"),
		slog.String("content", ""),
		slog.String("event", "session_start"),
	}

	// Add metadata
	if len(metadata) > 0 {
		metadataAttrs := make([]any, 0, len(metadata)*2)
		for k, v := range metadata {
			metadataAttrs = append(metadataAttrs, k, v)
		}
		attrs = append(attrs, slog.Group("metadata", metadataAttrs...))
	}

	l.slogger.LogAttrs(context.TODO(), slog.LevelInfo, "conversation", attrs...)
	//nolint:errcheck // Sync errors are not critical for logging
	os.Stdout.Sync()
}

// LogSessionEnd logs a session end event
func (l *Logger) LogSessionEnd(sessionID string) {
	l.slogger.Info("conversation",
		slog.String("service", l.Service),
		slog.String("session_id", sessionID),
		slog.String("role", "system"),
		slog.String("content", ""),
		slog.String("event", "session_end"),
	)
	//nolint:errcheck // Sync errors are not critical for logging
	os.Stdout.Sync()
}

// Adapter interface for future service implementations
type Adapter interface {
	// ParseEvent parses service-specific events into log attributes
	ParseEvent(eventType string, data interface{}) ([]slog.Attr, error)

	ServiceName() string
}
