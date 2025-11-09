package fileops

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/duclacloud/DUCLA-CLOUD-AGENT/internal/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Manager manages file operations
type Manager struct {
	config config.StorageConfig
	logger *logrus.Logger

	// Transfer management
	mu        sync.RWMutex
	transfers map[string]*Transfer
	
	// Cleanup
	cleanupTicker *time.Ticker
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// Transfer represents a file transfer operation
type Transfer struct {
	ID            string                 `json:"id"`
	Type          TransferType           `json:"type"`
	Status        TransferStatus         `json:"status"`
	SourcePath    string                 `json:"source_path"`
	DestPath      string                 `json:"dest_path"`
	Size          int64                  `json:"size"`
	Transferred   int64                  `json:"transferred"`
	Progress      float64                `json:"progress"`
	Checksum      string                 `json:"checksum"`
	StartedAt     time.Time              `json:"started_at"`
	CompletedAt   time.Time              `json:"completed_at,omitempty"`
	Error         string                 `json:"error,omitempty"`
	Metadata      map[string]interface{} `json:"metadata"`
	
	// Internal
	ctx    context.Context
	cancel context.CancelFunc
}

// TransferType represents the type of transfer
type TransferType string

const (
	TransferTypeUpload   TransferType = "upload"
	TransferTypeDownload TransferType = "download"
	TransferTypeCopy     TransferType = "copy"
	TransferTypeMove     TransferType = "move"
)

// TransferStatus represents the status of a transfer
type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "pending"
	TransferStatusRunning   TransferStatus = "running"
	TransferStatusCompleted TransferStatus = "completed"
	TransferStatusFailed    TransferStatus = "failed"
	TransferStatusCancelled TransferStatus = "cancelled"
)

// Operation represents a file operation request
type Operation struct {
	Type       OperationType          `json:"type"`
	SourcePath string                 `json:"source_path"`
	DestPath   string                 `json:"dest_path"`
	Mode       os.FileMode            `json:"mode,omitempty"`
	Recursive  bool                   `json:"recursive"`
	Overwrite  bool                   `json:"overwrite"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// OperationType represents the type of file operation
type OperationType string

const (
	OperationTypeUpload   OperationType = "upload"
	OperationTypeDownload OperationType = "download"
	OperationTypeCopy     OperationType = "copy"
	OperationTypeMove     OperationType = "move"
	OperationTypeDelete   OperationType = "delete"
	OperationTypeList     OperationType = "list"
	OperationTypeStat     OperationType = "stat"
	OperationTypeChmod    OperationType = "chmod"
	OperationTypeChown    OperationType = "chown"
)

// FileInfo represents file information
type FileInfo struct {
	Name    string      `json:"name"`
	Path    string      `json:"path"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"mod_time"`
	IsDir   bool        `json:"is_dir"`
}

// New creates a new file operations manager
func New(cfg config.StorageConfig, logger *logrus.Logger) (*Manager, error) {
	manager := &Manager{
		config:    cfg,
		logger:    logger,
		transfers: make(map[string]*Transfer),
	}

	// Create directories if they don't exist
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	if err := os.MkdirAll(cfg.TempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return manager, nil
}

// Start starts the file operations manager
func (m *Manager) Start(ctx context.Context) error {
	m.logger.Info("Starting file operations manager")

	m.ctx, m.cancel = context.WithCancel(ctx)

	// Start cleanup routine if enabled
	if m.config.Cleanup.Enabled {
		m.cleanupTicker = time.NewTicker(m.config.Cleanup.Interval)
		m.wg.Add(1)
		go m.cleanupLoop()
	}

	m.logger.Info("File operations manager started")
	return nil
}

// Stop stops the file operations manager
func (m *Manager) Stop(ctx context.Context) error {
	m.logger.Info("Stopping file operations manager")

	// Cancel all active transfers
	m.mu.Lock()
	for _, transfer := range m.transfers {
		if transfer.cancel != nil {
			transfer.cancel()
		}
	}
	m.mu.Unlock()

	// Stop cleanup routine
	if m.cleanupTicker != nil {
		m.cleanupTicker.Stop()
	}

	// Cancel context
	if m.cancel != nil {
		m.cancel()
	}

	// Wait for goroutines
	m.wg.Wait()

	m.logger.Info("File operations manager stopped")
	return nil
}

// Name returns the service name
func (m *Manager) Name() string {
	return "fileops"
}

// ExecuteOperation executes a file operation
func (m *Manager) ExecuteOperation(ctx context.Context, op *Operation) (map[string]interface{}, error) {
	m.logger.WithFields(logrus.Fields{
		"type":   op.Type,
		"source": op.SourcePath,
		"dest":   op.DestPath,
	}).Info("Executing file operation")

	switch op.Type {
	case OperationTypeUpload:
		return m.handleUpload(ctx, op)
	case OperationTypeDownload:
		return m.handleDownload(ctx, op)
	case OperationTypeCopy:
		return m.handleCopy(ctx, op)
	case OperationTypeMove:
		return m.handleMove(ctx, op)
	case OperationTypeDelete:
		return m.handleDelete(ctx, op)
	case OperationTypeList:
		return m.handleList(ctx, op)
	case OperationTypeStat:
		return m.handleStat(ctx, op)
	case OperationTypeChmod:
		return m.handleChmod(ctx, op)
	case OperationTypeChown:
		return m.handleChown(ctx, op)
	default:
		return nil, fmt.Errorf("unsupported operation type: %s", op.Type)
	}
}

// handleUpload handles file upload operation
func (m *Manager) handleUpload(ctx context.Context, op *Operation) (map[string]interface{}, error) {
	transferID := uuid.New().String()
	
	transfer := &Transfer{
		ID:         transferID,
		Type:       TransferTypeUpload,
		Status:     TransferStatusPending,
		SourcePath: op.SourcePath,
		DestPath:   op.DestPath,
		StartedAt:  time.Now(),
		Metadata:   op.Metadata,
	}

	transfer.ctx, transfer.cancel = context.WithCancel(ctx)

	// Store transfer
	m.mu.Lock()
	m.transfers[transferID] = transfer
	m.mu.Unlock()

	// Execute upload
	go m.executeUpload(transfer)

	return map[string]interface{}{
		"transfer_id": transferID,
		"status":      transfer.Status,
	}, nil
}

// handleDownload handles file download operation
func (m *Manager) handleDownload(ctx context.Context, op *Operation) (map[string]interface{}, error) {
	transferID := uuid.New().String()
	
	transfer := &Transfer{
		ID:         transferID,
		Type:       TransferTypeDownload,
		Status:     TransferStatusPending,
		SourcePath: op.SourcePath,
		DestPath:   op.DestPath,
		StartedAt:  time.Now(),
		Metadata:   op.Metadata,
	}

	transfer.ctx, transfer.cancel = context.WithCancel(ctx)

	// Store transfer
	m.mu.Lock()
	m.transfers[transferID] = transfer
	m.mu.Unlock()

	// Execute download
	go m.executeDownload(transfer)

	return map[string]interface{}{
		"transfer_id": transferID,
		"status":      transfer.Status,
	}, nil
}

// handleCopy handles file copy operation
func (m *Manager) handleCopy(ctx context.Context, op *Operation) (map[string]interface{}, error) {
	// Validate paths
	if err := m.validatePath(op.SourcePath); err != nil {
		return nil, fmt.Errorf("invalid source path: %w", err)
	}
	if err := m.validatePath(op.DestPath); err != nil {
		return nil, fmt.Errorf("invalid dest path: %w", err)
	}

	// Check if source exists
	srcInfo, err := os.Stat(op.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("source file not found: %w", err)
	}

	// Copy file or directory
	var bytesCopied int64
	if srcInfo.IsDir() {
		if !op.Recursive {
			return nil, fmt.Errorf("source is a directory, use recursive flag")
		}
		bytesCopied, err = m.copyDir(op.SourcePath, op.DestPath)
	} else {
		bytesCopied, err = m.copyFile(op.SourcePath, op.DestPath)
	}

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"source":       op.SourcePath,
		"destination":  op.DestPath,
		"bytes_copied": bytesCopied,
		"is_dir":       srcInfo.IsDir(),
	}, nil
}

// handleMove handles file move operation
func (m *Manager) handleMove(ctx context.Context, op *Operation) (map[string]interface{}, error) {
	// Validate paths
	if err := m.validatePath(op.SourcePath); err != nil {
		return nil, fmt.Errorf("invalid source path: %w", err)
	}
	if err := m.validatePath(op.DestPath); err != nil {
		return nil, fmt.Errorf("invalid dest path: %w", err)
	}

	// Move file
	if err := os.Rename(op.SourcePath, op.DestPath); err != nil {
		return nil, fmt.Errorf("failed to move file: %w", err)
	}

	return map[string]interface{}{
		"source":      op.SourcePath,
		"destination": op.DestPath,
	}, nil
}

// handleDelete handles file delete operation
func (m *Manager) handleDelete(ctx context.Context, op *Operation) (map[string]interface{}, error) {
	// Validate path
	if err := m.validatePath(op.SourcePath); err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// Check if file exists
	info, err := os.Stat(op.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Delete file or directory
	if info.IsDir() && !op.Recursive {
		return nil, fmt.Errorf("path is a directory, use recursive flag")
	}

	if err := os.RemoveAll(op.SourcePath); err != nil {
		return nil, fmt.Errorf("failed to delete: %w", err)
	}

	return map[string]interface{}{
		"path":    op.SourcePath,
		"deleted": true,
	}, nil
}

// handleList handles directory listing operation
func (m *Manager) handleList(ctx context.Context, op *Operation) (map[string]interface{}, error) {
	// Validate path
	if err := m.validatePath(op.SourcePath); err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// Read directory
	entries, err := os.ReadDir(op.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Build file list
	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Name:    entry.Name(),
			Path:    filepath.Join(op.SourcePath, entry.Name()),
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
			IsDir:   entry.IsDir(),
		})
	}

	return map[string]interface{}{
		"path":  op.SourcePath,
		"count": len(files),
		"files": files,
	}, nil
}

// handleStat handles file stat operation
func (m *Manager) handleStat(ctx context.Context, op *Operation) (map[string]interface{}, error) {
	// Validate path
	if err := m.validatePath(op.SourcePath); err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// Get file info
	info, err := os.Stat(op.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	fileInfo := FileInfo{
		Name:    info.Name(),
		Path:    op.SourcePath,
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
	}

	return map[string]interface{}{
		"file": fileInfo,
	}, nil
}

// handleChmod handles chmod operation
func (m *Manager) handleChmod(ctx context.Context, op *Operation) (map[string]interface{}, error) {
	// Validate path
	if err := m.validatePath(op.SourcePath); err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// Change mode
	if err := os.Chmod(op.SourcePath, op.Mode); err != nil {
		return nil, fmt.Errorf("failed to chmod: %w", err)
	}

	return map[string]interface{}{
		"path": op.SourcePath,
		"mode": op.Mode.String(),
	}, nil
}

// handleChown handles chown operation
func (m *Manager) handleChown(ctx context.Context, op *Operation) (map[string]interface{}, error) {
	// Validate path
	if err := m.validatePath(op.SourcePath); err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	// Get UID and GID from metadata
	uid, _ := op.Metadata["uid"].(int)
	gid, _ := op.Metadata["gid"].(int)

	// Change ownership
	if err := os.Chown(op.SourcePath, uid, gid); err != nil {
		return nil, fmt.Errorf("failed to chown: %w", err)
	}

	return map[string]interface{}{
		"path": op.SourcePath,
		"uid":  uid,
		"gid":  gid,
	}, nil
}

// executeUpload executes file upload
func (m *Manager) executeUpload(transfer *Transfer) {
	transfer.Status = TransferStatusRunning

	// Implementation would handle actual upload from master server
	// This is a placeholder
	m.logger.WithField("transfer_id", transfer.ID).Info("Upload not yet implemented")

	transfer.Status = TransferStatusCompleted
	transfer.CompletedAt = time.Now()
}

// executeDownload executes file download
func (m *Manager) executeDownload(transfer *Transfer) {
	transfer.Status = TransferStatusRunning

	// Implementation would handle actual download to master server
	// This is a placeholder
	m.logger.WithField("transfer_id", transfer.ID).Info("Download not yet implemented")

	transfer.Status = TransferStatusCompleted
	transfer.CompletedAt = time.Now()
}

// copyFile copies a single file
func (m *Manager) copyFile(src, dst string) (int64, error) {
	sourceFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destFile.Close()

	bytesCopied, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return 0, err
	}

	// Copy permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return bytesCopied, err
	}
	return bytesCopied, os.Chmod(dst, sourceInfo.Mode())
}

// copyDir copies a directory recursively
func (m *Manager) copyDir(src, dst string) (int64, error) {
	var totalBytes int64

	return totalBytes, filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		bytes, err := m.copyFile(path, destPath)
		totalBytes += bytes
		return err
	})
}

// validatePath validates a file path
func (m *Manager) validatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	// Check if path is absolute
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute")
	}

	// Additional validation can be added here
	return nil
}

// CalculateChecksum calculates file checksum
func (m *Manager) CalculateChecksum(path string, algorithm string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	switch algorithm {
	case "md5":
		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			return "", err
		}
		return hex.EncodeToString(hash.Sum(nil)), nil
	case "sha256":
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			return "", err
		}
		return hex.EncodeToString(hash.Sum(nil)), nil
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

// GetTransfer retrieves a transfer by ID
func (m *Manager) GetTransfer(transferID string) (*Transfer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	transfer, exists := m.transfers[transferID]
	if !exists {
		return nil, fmt.Errorf("transfer not found: %s", transferID)
	}

	return transfer, nil
}

// CancelTransfer cancels a transfer
func (m *Manager) CancelTransfer(transferID string) error {
	m.mu.RLock()
	transfer, exists := m.transfers[transferID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("transfer not found: %s", transferID)
	}

	if transfer.cancel != nil {
		transfer.cancel()
	}

	transfer.Status = TransferStatusCancelled
	return nil
}

// cleanupLoop performs periodic cleanup
func (m *Manager) cleanupLoop() {
	defer m.wg.Done()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-m.cleanupTicker.C:
			m.performCleanup()
		}
	}
}

// performCleanup performs cleanup of old files
func (m *Manager) performCleanup() {
	m.logger.Debug("Performing cleanup")

	// Clean temp directory
	if err := m.cleanupDirectory(m.config.TempDir, m.config.Cleanup.MaxAge); err != nil {
		m.logger.WithError(err).Error("Failed to cleanup temp directory")
	}

	// Clean completed transfers
	m.cleanupTransfers()
}

// cleanupDirectory cleans up old files in a directory
func (m *Manager) cleanupDirectory(dir string, maxAge time.Duration) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if time.Since(info.ModTime()) > maxAge {
			m.logger.WithField("file", path).Debug("Removing old file")
			return os.Remove(path)
		}

		return nil
	})
}

// cleanupTransfers removes old completed transfers
func (m *Manager) cleanupTransfers() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, transfer := range m.transfers {
		if transfer.Status == TransferStatusCompleted || transfer.Status == TransferStatusFailed {
			if time.Since(transfer.CompletedAt) > 24*time.Hour {
				delete(m.transfers, id)
			}
		}
	}
}

// ParseOperation parses operation data from a map
func ParseOperation(data map[string]interface{}) (*Operation, error) {
	op := &Operation{
		Metadata: make(map[string]interface{}),
	}

	if opType, ok := data["type"].(string); ok {
		op.Type = OperationType(opType)
	}

	if sourcePath, ok := data["source_path"].(string); ok {
		op.SourcePath = sourcePath
	}

	if destPath, ok := data["dest_path"].(string); ok {
		op.DestPath = destPath
	}

	if recursive, ok := data["recursive"].(bool); ok {
		op.Recursive = recursive
	}

	if overwrite, ok := data["overwrite"].(bool); ok {
		op.Overwrite = overwrite
	}

	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		op.Metadata = metadata
	}

	return op, nil
}