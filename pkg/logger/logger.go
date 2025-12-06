package logger

import (
	"log/slog"
	"os"

	"github.com/hirosassa/tapline/pkg/session"
)

// Logger handles conversation logging using slog
type Logger struct {
	slogger        *slog.Logger
	Service        string
	SessionManager *session.Manager
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
	os.Stdout.Sync() // Ensure immediate flush to disk
}

// LogAssistantResponse logs an assistant response
func (l *Logger) LogAssistantResponse(sessionID, content string) {
	l.slogger.Info("conversation",
		slog.String("service", l.Service),
		slog.String("session_id", sessionID),
		slog.String("role", "assistant"),
		slog.String("content", content),
	)
	os.Stdout.Sync() // Ensure immediate flush to disk
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

	l.slogger.LogAttrs(nil, slog.LevelInfo, "conversation", attrs...)
	os.Stdout.Sync() // Ensure immediate flush to disk
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
	os.Stdout.Sync() // Ensure immediate flush to disk
}

// Adapter interface for future service implementations
type Adapter interface {
	// ParseEvent parses service-specific events into log attributes
	ParseEvent(eventType string, data interface{}) ([]slog.Attr, error)

	// ServiceName returns the name of the service
	ServiceName() string
}
