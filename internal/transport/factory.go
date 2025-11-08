package transport

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ducla/cloud-agent/internal/config"
	"github.com/sirupsen/logrus"
)

// New creates a new transport instance based on the master URL scheme
func New(cfg config.MasterConfig, logger *logrus.Logger) (Transport, error) {
	parsedURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid master URL: %w", err)
	}

	transportConfig := TransportConfig{
		URL:                  cfg.URL,
		Token:                cfg.Token,
		ConnectTimeout:       cfg.ConnectTimeout,
		HeartbeatInterval:    cfg.HeartbeatInterval,
		ReconnectInterval:    cfg.ReconnectInterval,
		MaxReconnectAttempts: cfg.MaxReconnectAttempts,
		TLSSkipVerify:        cfg.TLSSkipVerify,
	}

	switch strings.ToLower(parsedURL.Scheme) {
	case "ws", "wss":
		transportConfig.Protocol = "websocket"
		return NewWebSocketTransport(transportConfig, logger)
	case "grpc", "grpcs":
		transportConfig.Protocol = "grpc"
		return NewGRPCTransport(transportConfig, logger)
	case "http", "https":
		// Default to WebSocket for HTTP URLs
		transportConfig.Protocol = "websocket"
		// Convert HTTP(S) to WS(S)
		if parsedURL.Scheme == "http" {
			transportConfig.URL = strings.Replace(cfg.URL, "http://", "ws://", 1)
		} else {
			transportConfig.URL = strings.Replace(cfg.URL, "https://", "wss://", 1)
		}
		return NewWebSocketTransport(transportConfig, logger)
	default:
		return nil, fmt.Errorf("unsupported protocol scheme: %s", parsedURL.Scheme)
	}
}