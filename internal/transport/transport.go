package transport

import (
	"context"
	"time"
)

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeHeartbeat           MessageType = "heartbeat"
	MessageTypeTask                MessageType = "task"
	MessageTypeTaskResult          MessageType = "task_result"
	MessageTypeFileOperation       MessageType = "file_operation"
	MessageTypeFileOperationResult MessageType = "file_operation_result"
	MessageTypeHealthCheck         MessageType = "health_check"
	MessageTypeHealthCheckResult   MessageType = "health_check_result"
	MessageTypeConfig              MessageType = "config"
	MessageTypeConfigResult        MessageType = "config_result"
	MessageTypeError               MessageType = "error"
	MessageTypeLog                 MessageType = "log"
	MessageTypeMetrics             MessageType = "metrics"
)

// Message represents a message exchanged between agent and master
type Message struct {
	ID        string                 `json:"id"`
	Type      MessageType            `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	AgentID   string                 `json:"agent_id,omitempty"`
	ReplyTo   string                 `json:"reply_to,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]string      `json:"metadata,omitempty"`
}

// Transport defines the interface for communication with master server
type Transport interface {
	// Connect establishes connection to master server
	Connect(ctx context.Context) error

	// Disconnect closes connection to master server
	Disconnect() error

	// SendMessage sends a message to master server
	SendMessage(message *Message) error

	// ReceiveMessage receives a message from master server
	ReceiveMessage(ctx context.Context) (*Message, error)

	// IsConnected returns whether the transport is connected
	IsConnected() bool

	// GetConnectionInfo returns connection information
	GetConnectionInfo() ConnectionInfo
}

// ConnectionInfo contains information about the connection
type ConnectionInfo struct {
	Connected      bool      `json:"connected"`
	ConnectedAt    time.Time `json:"connected_at,omitempty"`
	DisconnectedAt time.Time `json:"disconnected_at,omitempty"`
	MasterURL      string    `json:"master_url"`
	Protocol       string    `json:"protocol"`
	Latency        int64     `json:"latency_ms"`
	MessagesSent   int64     `json:"messages_sent"`
	MessagesRecv   int64     `json:"messages_recv"`
	Errors         int64     `json:"errors"`
}

// TransportConfig contains transport configuration
type TransportConfig struct {
	URL                  string
	Token                string
	ConnectTimeout       time.Duration
	HeartbeatInterval    time.Duration
	ReconnectInterval    time.Duration
	MaxReconnectAttempts int
	TLSSkipVerify        bool
	Protocol             string // "websocket" or "grpc"
}

// Error types
type TransportError struct {
	Code    string
	Message string
	Err     error
}

func (e *TransportError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *TransportError) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	ErrCodeConnectionFailed   = "CONNECTION_FAILED"
	ErrCodeAuthenticationFailed = "AUTHENTICATION_FAILED"
	ErrCodeTimeout            = "TIMEOUT"
	ErrCodeDisconnected       = "DISCONNECTED"
	ErrCodeInvalidMessage     = "INVALID_MESSAGE"
	ErrCodeSendFailed         = "SEND_FAILED"
	ErrCodeReceiveFailed      = "RECEIVE_FAILED"
)