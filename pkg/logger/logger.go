// Package logger provides conversation logging functionality using structured logging.
// It supports logging user prompts, assistant responses, and session lifecycle events
// in JSON format to stdout with immediate flushing for durability.
package logger

import (
	"log/slog"
	"os"

	"github.com/hirosassa/tapline/pkg/git"
	"github.com/hirosassa/tapline/pkg/session"
	"github.com/hirosassa/tapline/pkg/user"
)

// Logger handles conversation logging using slog
type Logger struct {
	slogger        *slog.Logger
	SessionManager *session.Manager
	Service        string
	UserID         string
	UserSource     string
	Hostname       string
	GitRepoURL     string
	GitRepoName    string
	GitBranch      string
}

// NewLogger creates a new Logger instance with slog JSON handler
func NewLogger(service string, sessionMgr *session.Manager) *Logger {
	// Create JSON handler that writes to stdout
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	userInfo := user.GetIdentifier(service)

	var gitRepoURL, gitRepoName, gitBranch string
	if repoInfo, err := git.GetRepoInfo(); err == nil {
		gitRepoURL = repoInfo.OriginURL
		gitRepoName = repoInfo.RepoName
		gitBranch = repoInfo.Branch
	} else {
		// Log at debug level; git info is optional, so we don't fail
		tmpLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		tmpLogger.Debug("Failed to retrieve git repo info", slog.Any("error", err))
	}

	return &Logger{
		slogger:        slog.New(handler),
		Service:        service,
		SessionManager: sessionMgr,
		UserID:         userInfo.UserID,
		UserSource:     userInfo.Source,
		Hostname:       userInfo.Hostname,
		GitRepoURL:     gitRepoURL,
		GitRepoName:    gitRepoName,
		GitBranch:      gitBranch,
	}
}

// appendGitAttrs appends Git-related attributes to the slice if they are set
func (l *Logger) appendGitAttrs(attrs []any) []any {
	if l.GitRepoURL != "" {
		attrs = append(attrs, slog.String("git_repo_url", l.GitRepoURL))
	}
	if l.GitRepoName != "" {
		attrs = append(attrs, slog.String("git_repo_name", l.GitRepoName))
	}
	if l.GitBranch != "" {
		attrs = append(attrs, slog.String("git_branch", l.GitBranch))
	}
	return attrs
}

// LogUserPrompt logs a user prompt
func (l *Logger) LogUserPrompt(sessionID, content string) {
	attrs := []any{
		slog.String("service", l.Service),
		slog.String("session_id", sessionID),
		slog.String("user_id", l.UserID),
		slog.String("user_source", l.UserSource),
		slog.String("hostname", l.Hostname),
	}

	attrs = l.appendGitAttrs(attrs)

	attrs = append(attrs,
		slog.String("role", "user"),
		slog.String("content", content),
	)

	l.slogger.Info("conversation", attrs...)
	//nolint:errcheck // Sync errors are not critical for logging
	os.Stdout.Sync()
}

// LogAssistantResponse logs an assistant response
func (l *Logger) LogAssistantResponse(sessionID, content string) {
	attrs := []any{
		slog.String("service", l.Service),
		slog.String("session_id", sessionID),
		slog.String("user_id", l.UserID),
		slog.String("user_source", l.UserSource),
		slog.String("hostname", l.Hostname),
	}

	attrs = l.appendGitAttrs(attrs)

	attrs = append(attrs,
		slog.String("role", "assistant"),
		slog.String("content", content),
	)

	l.slogger.Info("conversation", attrs...)
	//nolint:errcheck // Sync errors are not critical for logging
	os.Stdout.Sync()
}

// LogSessionStart logs a session start event
func (l *Logger) LogSessionStart(sessionID string, metadata map[string]string) {
	attrs := []any{
		slog.String("service", l.Service),
		slog.String("session_id", sessionID),
		slog.String("user_id", l.UserID),
		slog.String("user_source", l.UserSource),
		slog.String("hostname", l.Hostname),
	}

	attrs = l.appendGitAttrs(attrs)

	attrs = append(attrs,
		slog.String("role", "system"),
		slog.String("content", ""),
		slog.String("event", "session_start"),
	)

	// Add metadata
	if len(metadata) > 0 {
		metadataAttrs := make([]any, 0, len(metadata)*2)
		for k, v := range metadata {
			metadataAttrs = append(metadataAttrs, k, v)
		}
		attrs = append(attrs, slog.Group("metadata", metadataAttrs...))
	}

	l.slogger.Info("conversation", attrs...)
	//nolint:errcheck // Sync errors are not critical for logging
	os.Stdout.Sync()
}

// LogSessionEnd logs a session end event
func (l *Logger) LogSessionEnd(sessionID string) {
	attrs := []any{
		slog.String("service", l.Service),
		slog.String("session_id", sessionID),
		slog.String("user_id", l.UserID),
		slog.String("user_source", l.UserSource),
		slog.String("hostname", l.Hostname),
	}

	attrs = l.appendGitAttrs(attrs)

	attrs = append(attrs,
		slog.String("role", "system"),
		slog.String("content", ""),
		slog.String("event", "session_end"),
	)

	l.slogger.Info("conversation", attrs...)
	//nolint:errcheck // Sync errors are not critical for logging
	os.Stdout.Sync()
}

// Adapter interface for future service implementations
type Adapter interface {
	// ParseEvent parses service-specific events into log attributes
	ParseEvent(eventType string, data interface{}) ([]slog.Attr, error)

	ServiceName() string
}
