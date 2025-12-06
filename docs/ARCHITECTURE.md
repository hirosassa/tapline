# Tapline Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        Claude Code                          │
│                                                             │
│  ┌──────────────┐    ┌──────────────┐   ┌──────────────┐  │
│  │ User Input   │───>│ Conversation │──>│ Assistant    │  │
│  │              │    │              │   │ Response     │  │
│  └──────────────┘    └──────────────┘   └──────────────┘  │
│         │                   │                    │         │
│         │                   │                    │         │
│         └───────────────────┴────────────────────┘         │
│                             │                              │
│                    (via hooks.json)                        │
└─────────────────────────────┼───────────────────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │     Tapline      │
                    │   CLI Binary     │
                    └──────────────────┘
                              │
                ┌─────────────┼─────────────┐
                │             │             │
                ▼             ▼             ▼
         ┌──────────┐  ┌──────────┐  ┌──────────┐
         │  Logger  │  │ Session  │  │ Adapter  │
         │  Module  │  │ Manager  │  │Interface │
         └──────────┘  └──────────┘  └──────────┘
                │             │
                │             ▼
                │      ┌─────────────┐
                │      │~/.tapline/  │
                │      │session_id   │
                │      └─────────────┘
                │
                ▼
          ┌──────────┐
          │  stdout  │
          │  (JSON   │
          │  Lines)  │
          └──────────┘
                │
    ┌───────────┼───────────┐
    │           │           │
    ▼           ▼           ▼
┌────────┐  ┌────────┐  ┌────────┐
│  File  │  │  jq    │  │  Log   │
│System  │  │Analysis│  │Aggreg. │
└────────┘  └────────┘  └────────┘
```

## Component Details

### 1. Claude Code Integration

Claude Code integrates with Tapline through its hooks system defined in `.claude/hooks.json`:

- `conversation_start`: Triggered when a conversation begins
- `conversation_end`: Triggered when a conversation ends
- `user_prompt_submit`: Triggered when user submits a prompt
- `assistant_response`: Triggered when assistant responds

### 2. Tapline CLI Binary

The main executable that processes hook events and generates structured logs.

**Commands:**
- `tapline conversation_start` - Creates new session
- `tapline conversation_end` - Ends current session
- `tapline user_prompt <text>` - Logs user message
- `tapline assistant_response <text>` - Logs assistant message

### 3. Logger Module (`pkg/logger`)

Responsible for formatting and outputting conversation entries as JSON Lines.

**Key Types:**
```go
type ConversationEntry struct {
    Timestamp time.Time
    Service   string
    SessionID string
    Role      string
    Content   string
    Metadata  map[string]interface{}
    Event     string
}
```

**Responsibilities:**
- Format log entries consistently
- Output to stdout in JSON Lines format
- Provide adapter interface for future services

### 4. Session Manager (`pkg/session`)

Manages conversation session IDs persistently across invocations.

**Storage Location:** `~/.tapline/session_id`

**Operations:**
- `GetSessionID()` - Retrieve current session
- `SetSessionID(id)` - Store new session
- `ClearSession()` - Remove session file
- `HasActiveSession()` - Check if session exists

### 5. Adapter Interface

Provides extensibility for additional AI services:

```go
type Adapter interface {
    ParseEvent(eventType string, data interface{}) (ConversationEntry, error)
    ServiceName() string
}
```

## Data Flow

### Conversation Start

```
1. Claude Code hook triggered
2. Tapline generates UUID
3. Session Manager stores ID
4. Logger outputs session_start event
```

### User/Assistant Messages

```
1. Claude Code hook triggered with content
2. Session Manager retrieves session ID
3. Logger formats entry with ID
4. JSON line written to stdout
```

### Conversation End

```
1. Claude Code hook triggered
2. Session Manager retrieves session ID
3. Logger outputs session_end event
4. Session Manager clears session file
```

## Log Format Specification

### Standard Entry
```json
{
  "timestamp": "2025-12-06T16:19:02.935561+09:00",
  "service": "claude-code",
  "session_id": "0e97d08c-b08b-4f5a-92ea-086c36d5818b",
  "role": "user|assistant|system",
  "content": "message text"
}
```

### Session Event
```json
{
  "timestamp": "2025-12-06T16:19:02.935561+09:00",
  "service": "claude-code",
  "session_id": "0e97d08c-b08b-4f5a-92ea-086c36d5818b",
  "role": "system",
  "content": "",
  "event": "session_start|session_end",
  "metadata": {
    "hostname": "my-laptop",
    "cwd": "/path/to/project"
  }
}
```

## Future Extensions

### Phase 2: Additional Services

```
┌─────────────────────────────────────────────────┐
│             Tapline Core                        │
│  ┌────────────────────────────────────────┐    │
│  │      Unified Logger Interface          │    │
│  └────────────────────────────────────────┘    │
│           ▲         ▲          ▲               │
│           │         │          │               │
│  ┌────────┴───┐ ┌──┴──────┐ ┌─┴──────────┐    │
│  │ Claude     │ │ Gemini  │ │ ChatGPT    │    │
│  │ Adapter    │ │ Adapter │ │ Adapter    │    │
│  └────────────┘ └─────────┘ └────────────┘    │
└─────────────────────────────────────────────────┘
```

Each adapter implements:
- Service-specific parsing
- Consistent output format
- Error handling
- Session management

### Planned Features

1. **DuckDB Adapter**: Direct database storage option
2. **Log Rotation**: Automatic file rotation and compression
3. **Real-time Streaming**: WebSocket support for live monitoring
4. **Analytics Module**: Built-in conversation metrics
5. **Web UI**: Browser-based log viewer and search

## Design Principles

1. **Separation of Concerns**: Logger, session management, and adapters are independent
2. **Extensibility**: Easy to add new services via adapter pattern
3. **Simplicity**: Structured logging to stdout keeps system simple
4. **Testability**: Each component is independently testable
5. **Performance**: Single binary, minimal overhead, fast execution
6. **Reliability**: Graceful error handling, no data loss

## Testing Strategy

### Unit Tests
- Logger module: Output format validation
- Session Manager: Lifecycle management
- Adapters: Event parsing accuracy

### Integration Tests
- End-to-end conversation flow
- Session persistence across invocations
- Error recovery scenarios

### Manual Testing
- Claude Code hook integration
- Log output verification
- Session ID consistency
