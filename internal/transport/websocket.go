package transport

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// WebSocketTransport implements Transport interface using WebSocket
type WebSocketTransport struct {
	config TransportConfig
	logger *logrus.Logger

	// Connection management
	mu         sync.RWMutex
	conn       *websocket.Conn
	connected  bool
	connInfo   ConnectionInfo
	
	// Message handling
	sendChan    chan *Message
	receiveChan chan *Message
	errorChan   chan error
	
	// Lifecycle management
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewWebSocketTransport creates a new WebSocket transport
func NewWebSocketTransport(config TransportConfig, logger *logrus.Logger) (*WebSocketTransport, error) {
	transport := &WebSocketTransport{
		config:      config,
		logger:      logger,
		sendChan:    make(chan *Message, 100),
		receiveChan: make(chan *Message, 100),
		errorChan:   make(chan error, 10),
		connInfo: ConnectionInfo{
			MasterURL: config.URL,
			Protocol:  "websocket",
		},
	}

	return transport, nil
}

// Connect establishes WebSocket connection to master server
func (t *WebSocketTransport) Connect(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.connected {
		return nil
	}

	t.logger.WithField("url", t.config.URL).Info("Connecting to master server via WebSocket")

	// Parse URL
	u, err := url.Parse(t.config.URL)
	if err != nil {
		return &TransportError{
			Code:    ErrCodeConnectionFailed,
			Message: "Invalid URL",
			Err:     err,
		}
	}

	// Setup WebSocket dialer
	dialer := websocket.Dialer{
		HandshakeTimeout: t.config.ConnectTimeout,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: t.config.TLSSkipVerify,
		},
	}

	// Setup headers for authentication
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+t.config.Token)
	headers.Set("User-Agent", "Ducla-Cloud-Agent/1.0.0")

	// Establish connection
	conn, resp, err := dialer.DialContext(ctx, u.String(), headers)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusUnauthorized {
			return &TransportError{
				Code:    ErrCodeAuthenticationFailed,
				Message: "Authentication failed",
				Err:     err,
			}
		}
		return &TransportError{
			Code:    ErrCodeConnectionFailed,
			Message: "Failed to connect to master server",
			Err:     err,
		}
	}

	t.conn = conn
	t.connected = true
	t.connInfo.Connected = true
	t.connInfo.ConnectedAt = time.Now()
	t.connInfo.MessagesSent = 0
	t.connInfo.MessagesRecv = 0
	t.connInfo.Errors = 0

	// Start connection context
	t.ctx, t.cancel = context.WithCancel(ctx)

	// Start message handling goroutines
	t.wg.Add(3)
	go t.sendLoop()
	go t.receiveLoop()
	go t.pingLoop()

	t.logger.Info("Connected to master server successfully")
	return nil
}

// Disconnect closes the WebSocket connection
func (t *WebSocketTransport) Disconnect() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.connected {
		return nil
	}

	t.logger.Info("Disconnecting from master server")

	// Cancel context to stop goroutines
	if t.cancel != nil {
		t.cancel()
	}

	// Close WebSocket connection
	if t.conn != nil {
		t.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		t.conn.Close()
		t.conn = nil
	}

	// Wait for goroutines to finish
	t.wg.Wait()

	t.connected = false
	t.connInfo.Connected = false
	t.connInfo.DisconnectedAt = time.Now()

	t.logger.Info("Disconnected from master server")
	return nil
}

// SendMessage sends a message to the master server
func (t *WebSocketTransport) SendMessage(message *Message) error {
	if !t.IsConnected() {
		return &TransportError{
			Code:    ErrCodeDisconnected,
			Message: "Not connected to master server",
		}
	}

	// Set message metadata
	if message.ID == "" {
		message.ID = uuid.New().String()
	}
	message.Timestamp = time.Now()

	select {
	case t.sendChan <- message:
		return nil
	case <-time.After(5 * time.Second):
		return &TransportError{
			Code:    ErrCodeTimeout,
			Message: "Send channel is full",
		}
	}
}

// ReceiveMessage receives a message from the master server
func (t *WebSocketTransport) ReceiveMessage(ctx context.Context) (*Message, error) {
	if !t.IsConnected() {
		return nil, &TransportError{
			Code:    ErrCodeDisconnected,
			Message: "Not connected to master server",
		}
	}

	select {
	case message := <-t.receiveChan:
		return message, nil
	case err := <-t.errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// IsConnected returns whether the transport is connected
func (t *WebSocketTransport) IsConnected() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.connected
}

// GetConnectionInfo returns connection information
func (t *WebSocketTransport) GetConnectionInfo() ConnectionInfo {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.connInfo
}

// sendLoop handles outgoing messages
func (t *WebSocketTransport) sendLoop() {
	defer t.wg.Done()

	for {
		select {
		case <-t.ctx.Done():
			return
		case message := <-t.sendChan:
			if err := t.sendMessage(message); err != nil {
				t.logger.WithError(err).Error("Failed to send message")
				t.mu.Lock()
				t.connInfo.Errors++
				t.mu.Unlock()
				
				select {
				case t.errorChan <- err:
				default:
				}
			} else {
				t.mu.Lock()
				t.connInfo.MessagesSent++
				t.mu.Unlock()
			}
		}
	}
}

// receiveLoop handles incoming messages
func (t *WebSocketTransport) receiveLoop() {
	defer t.wg.Done()

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			message, err := t.receiveMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					t.logger.Info("WebSocket connection closed normally")
					return
				}
				
				t.logger.WithError(err).Error("Failed to receive message")
				t.mu.Lock()
				t.connInfo.Errors++
				t.mu.Unlock()
				
				select {
				case t.errorChan <- err:
				default:
				}
				
				// Reconnect on error
				go t.reconnect()
				return
			}

			t.mu.Lock()
			t.connInfo.MessagesRecv++
			t.mu.Unlock()

			select {
			case t.receiveChan <- message:
			case <-t.ctx.Done():
				return
			}
		}
	}
}

// pingLoop sends periodic ping messages to keep connection alive
func (t *WebSocketTransport) pingLoop() {
	defer t.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-t.ctx.Done():
			return
		case <-ticker.C:
			if err := t.ping(); err != nil {
				t.logger.WithError(err).Error("Failed to send ping")
			}
		}
	}
}

// sendMessage sends a single message over WebSocket
func (t *WebSocketTransport) sendMessage(message *Message) error {
	t.mu.RLock()
	conn := t.conn
	t.mu.RUnlock()

	if conn == nil {
		return &TransportError{
			Code:    ErrCodeDisconnected,
			Message: "Connection is nil",
		}
	}

	data, err := json.Marshal(message)
	if err != nil {
		return &TransportError{
			Code:    ErrCodeInvalidMessage,
			Message: "Failed to marshal message",
			Err:     err,
		}
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return &TransportError{
			Code:    ErrCodeSendFailed,
			Message: "Failed to write message",
			Err:     err,
		}
	}

	return nil
}

// receiveMessage receives a single message from WebSocket
func (t *WebSocketTransport) receiveMessage() (*Message, error) {
	t.mu.RLock()
	conn := t.conn
	t.mu.RUnlock()

	if conn == nil {
		return nil, &TransportError{
			Code:    ErrCodeDisconnected,
			Message: "Connection is nil",
		}
	}

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	messageType, data, err := conn.ReadMessage()
	if err != nil {
		return nil, &TransportError{
			Code:    ErrCodeReceiveFailed,
			Message: "Failed to read message",
			Err:     err,
		}
	}

	if messageType != websocket.TextMessage {
		return nil, &TransportError{
			Code:    ErrCodeInvalidMessage,
			Message: fmt.Sprintf("Unexpected message type: %d", messageType),
		}
	}

	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		return nil, &TransportError{
			Code:    ErrCodeInvalidMessage,
			Message: "Failed to unmarshal message",
			Err:     err,
		}
	}

	return &message, nil
}

// ping sends a ping message to keep connection alive
func (t *WebSocketTransport) ping() error {
	t.mu.RLock()
	conn := t.conn
	t.mu.RUnlock()

	if conn == nil {
		return &TransportError{
			Code:    ErrCodeDisconnected,
			Message: "Connection is nil",
		}
	}

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	return conn.WriteMessage(websocket.PingMessage, nil)
}

// reconnect attempts to reconnect to the master server
func (t *WebSocketTransport) reconnect() {
	t.logger.Info("Attempting to reconnect to master server")

	attempts := 0
	for attempts < t.config.MaxReconnectAttempts {
		attempts++
		
		t.logger.WithField("attempt", attempts).Info("Reconnecting...")
		
		// Disconnect first
		t.Disconnect()
		
		// Wait before reconnecting
		time.Sleep(t.config.ReconnectInterval)
		
		// Attempt to connect
		ctx, cancel := context.WithTimeout(context.Background(), t.config.ConnectTimeout)
		err := t.Connect(ctx)
		cancel()
		
		if err == nil {
			t.logger.Info("Reconnected successfully")
			return
		}
		
		t.logger.WithError(err).WithField("attempt", attempts).Error("Reconnection failed")
	}
	
	t.logger.Error("Max reconnection attempts reached, giving up")
}