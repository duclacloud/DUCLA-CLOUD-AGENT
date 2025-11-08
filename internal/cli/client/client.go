package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

func (c *Client) GetHealth(ctx context.Context) (*HealthResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/health", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &health, nil
}

type MetricsResponse struct {
	System  map[string]interface{} `json:"system"`
	Process map[string]interface{} `json:"process"`
	Agent   map[string]interface{} `json:"agent"`
}

func (c *Client) GetMetrics(ctx context.Context) (*MetricsResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/metrics", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var metrics MetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &metrics, nil
}

type AgentInfo struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Environment  string            `json:"environment"`
	Region       string            `json:"region"`
	Zone         string            `json:"zone"`
	Tags         map[string]string `json:"tags"`
	Capabilities []string          `json:"capabilities"`
}

func (c *Client) GetAgentInfo(ctx context.Context) (*AgentInfo, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/agent/info", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info AgentInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &info, nil
}

type AgentStatus struct {
	Status         string    `json:"status"`
	Uptime         string    `json:"uptime"`
	ConnectedTo    string    `json:"connected_to"`
	LastHeartbeat  time.Time `json:"last_heartbeat"`
	ActiveTasks    int       `json:"active_tasks"`
	CompletedTasks int       `json:"completed_tasks"`
	FailedTasks    int       `json:"failed_tasks"`
}

func (c *Client) GetAgentStatus(ctx context.Context) (*AgentStatus, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/agent/status", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var status AgentStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &status, nil
}

func (c *Client) ListAgents(ctx context.Context) ([]AgentInfo, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/agents", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var agents []AgentInfo
	if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return agents, nil
}

type Task struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	Payload   map[string]interface{} `json:"payload"`
	Result    map[string]interface{} `json:"result,omitempty"`
	Error     string                 `json:"error,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func (c *Client) ListTasks(ctx context.Context, status string) ([]Task, error) {
	path := "/api/v1/tasks"
	if status != "" {
		path += "?status=" + status
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return tasks, nil
}

func (c *Client) GetTask(ctx context.Context, taskID string) (*Task, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/tasks/"+taskID, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

type ExecuteTaskRequest struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
	Timeout string                 `json:"timeout"`
}

type ExecuteTaskResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}

func (c *Client) ExecuteTask(ctx context.Context, taskType, payload, timeout string) (*ExecuteTaskResponse, error) {
	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &payloadMap); err != nil {
		return nil, fmt.Errorf("invalid payload JSON: %w", err)
	}

	req := ExecuteTaskRequest{
		Type:    taskType,
		Payload: payloadMap,
		Timeout: timeout,
	}

	resp, err := c.doRequest(ctx, "POST", "/api/v1/tasks", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ExecuteTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func (c *Client) CancelTask(ctx context.Context, taskID string) error {
	resp, err := c.doRequest(ctx, "DELETE", "/api/v1/tasks/"+taskID, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type FileInfo struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	IsDir     bool      `json:"is_dir"`
	ModTime   time.Time `json:"mod_time"`
	Mode      string    `json:"mode"`
}

func (c *Client) ListFiles(ctx context.Context, path string) ([]FileInfo, error) {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/files?path="+path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var files []FileInfo
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return files, nil
}

func (c *Client) UploadFile(ctx context.Context, file io.Reader, remotePath string) error {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/files?path="+remotePath, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) DownloadFile(ctx context.Context, remotePath string, dest io.Writer) error {
	resp, err := c.doRequest(ctx, "GET", "/api/v1/files/download?path="+remotePath, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(dest, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *Client) DeleteFile(ctx context.Context, remotePath string) error {
	resp, err := c.doRequest(ctx, "DELETE", "/api/v1/files?path="+remotePath, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
