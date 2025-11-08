package transport

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCTransport implements Transport interface using gRPC
type GRPCTransport struct {
	config TransportConfig
	logger *logrus.Logger

	// Connection management
	mu         sync.RWMutex
	conn       *grpc.ClientConn
	client     AgentServiceClient
	stream     AgentService_StreamClient
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

// NewGRPCTransport creates a new gRPC transport
func NewGRPCTransport(config TransportConfig, logger *logrus.Logger) (*GRPCTransport, error) {
	transport := &GRPCTransport{
		config:      config,
		logger:      logger,
		sendChan:    make(chan *Message, 100),
		receiveChan: make(chan *Message, 100),
		errorChan:   make(chan error, 10),
		connInfo: ConnectionInfo{
			MasterURL: config.URL,
			Protocol:  "grpc",
		},
	}

	return transport, nil
}

// Connect establishes gRPC connection to master server
func (t *GRPCTransport) Connect(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.connected {
		return nil
	}

	t.logger.WithField("url", t.config.URL).Info("Connecting to master server via gRPC")

	// Setup gRPC connection options
	var opts []grpc.DialOption

	// Setup credentials
	if t.config.TLSSkipVerify {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})))
	} else {
		// Check if URL uses secure scheme
		if t.config.URL[:5] == "grpcs" {
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
	}

	// Setup connection timeout
	opts = append(opts, grpc.WithBlock())
	
	// Setup keepalive
	opts = append(opts, grpc.WithKeepaliveParams(grpc.KeepaliveParams{
		Time:                30 * time.Second,
		Timeout:             5 * time.Second,
		PermitWithoutStream: true,
	}))

	// Extract address from URL
	address := t.config.URL
	if address[:4] == "grpc" {
		if address[:5] == "grpcs" {
			address = address[8:] // Remove "grpcs://"
		} else {
			address = address[7:] // Remove "grpc://"
		}
	}

	// Establish connection
	connectCtx, cancel := context.WithTimeout(ctx, t.config.ConnectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(connectCtx, address, opts...)
	if err != nil {
		return &TransportError{
			Code:    ErrCodeConnectionFailed,
			Message: "Failed to connect to master server",
			Err:     err,
		}
	}

	// Create client
	client := NewAgentServiceClient(conn)

	// Setup authentication metadata
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + t.config.Token,
		"user-agent":    "Ducla-Cloud-Agent/1.0.0",
	})
	authCtx := metadata.NewOutgoingContext(ctx, md)

	// Create bidirectional stream
	stream, err := client.Stream(authCtx)
	if err != nil {
		conn.Close()
		if status.Code(err) == codes.Unauthenticated {
			return &TransportError{
				Code:    ErrCodeAuthenticationFailed,
				Message: "Authentication failed",
				Err:     err,
			}
		}
		return &TransportError{
			Code:    ErrCodeConnectionFailed,
			Message: "Failed to create stream",
			Err:     err,
		}
	}

	t.conn = conn
	t.client = client
	t.stream = stream
	t.connected = true
	t.connInfo.Connected = true
	t.connInfo.ConnectedAt = time.Now()
	t.connInfo.MessagesSent = 0
	t.connInfo.MessagesRecv = 0
	t.connInfo.Errors = 0

	// Start connection context
	t.ctx, t.cancel = context.WithCancel(ctx)

	// Start message handling goroutines
	t.wg.Add(2)
	go t.sendLoop()
	go t.receiveLoop()

	t.logger.Info("Connected to master server successfully")
	return nil
}

// Disconnect closes the gRPC connection
func (t *GRPCTransport) Disconnect() error {
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

	// Close stream
	if t.stream != nil {
		t.stream.CloseSend()
		t.stream = nil
	}

	// Close gRPC connection
	if t.conn != nil {
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
func (t *GRPCTransport) SendMessage(message *Message) error {
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
func (t *GRPCTransport) ReceiveMessage(ctx context.Context) (*Message, error) {
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
func (t *GRPCTransport) IsConnected() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.connected
}

// GetConnectionInfo returns connection information
func (t *GRPCTransport) GetConnectionInfo() ConnectionInfo {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.connInfo
}

// sendLoop handles outgoing messages
func (t *GRPCTransport) sendLoop() {
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
func (t *GRPCTransport) receiveLoop() {
	defer t.wg.Done()

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			message, err := t.receiveMessage()
			if err != nil {
				if err == io.EOF {
					t.logger.Info("gRPC stream closed")
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

// sendMessage sends a single message over gRPC stream
func (t *GRPCTransport) sendMessage(message *Message) error {
	t.mu.RLock()
	stream := t.stream
	t.mu.RUnlock()

	if stream == nil {
		return &TransportError{
			Code:    ErrCodeDisconnected,
			Message: "Stream is nil",
		}
	}

	// Convert Message to gRPC message
	data, err := json.Marshal(message.Data)
	if err != nil {
		return &TransportError{
			Code:    ErrCodeInvalidMessage,
			Message: "Failed to marshal message data",
			Err:     err,
		}
	}

	grpcMessage := &StreamMessage{
		Id:        message.ID,
		Type:      string(message.Type),
		Timestamp: message.Timestamp.Unix(),
		AgentId:   message.AgentID,
		ReplyTo:   message.ReplyTo,
		Data:      data,
		Metadata:  message.Metadata,
	}

	if err := stream.Send(grpcMessage); err != nil {
		return &TransportError{
			Code:    ErrCodeSendFailed,
			Message: "Failed to send message",
			Err:     err,
		}
	}

	return nil
}

// receiveMessage receives a single message from gRPC stream
func (t *GRPCTransport) receiveMessage() (*Message, error) {
	t.mu.RLock()
	stream := t.stream
	t.mu.RUnlock()

	if stream == nil {
		return nil, &TransportError{
			Code:    ErrCodeDisconnected,
			Message: "Stream is nil",
		}
	}

	grpcMessage, err := stream.Recv()
	if err != nil {
		return nil, &TransportError{
			Code:    ErrCodeReceiveFailed,
			Message: "Failed to receive message",
			Err:     err,
		}
	}

	// Convert gRPC message to Message
	var data map[string]interface{}
	if err := json.Unmarshal(grpcMessage.Data, &data); err != nil {
		return nil, &TransportError{
			Code:    ErrCodeInvalidMessage,
			Message: "Failed to unmarshal message data",
			Err:     err,
		}
	}

	message := &Message{
		ID:        grpcMessage.Id,
		Type:      MessageType(grpcMessage.Type),
		Timestamp: time.Unix(grpcMessage.Timestamp, 0),
		AgentID:   grpcMessage.AgentId,
		ReplyTo:   grpcMessage.ReplyTo,
		Data:      data,
		Metadata:  grpcMessage.Metadata,
	}

	return message, nil
}

// reconnect attempts to reconnect to the master server
func (t *GRPCTransport) reconnect() {
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