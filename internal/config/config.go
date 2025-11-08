package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the agent configuration
type Config struct {
	Agent      AgentConfig      `yaml:"agent"`
	Master     MasterConfig     `yaml:"master"`
	API        APIConfig        `yaml:"api"`
	Security   SecurityConfig   `yaml:"security"`
	Storage    StorageConfig    `yaml:"storage"`
	Logging    LoggingConfig    `yaml:"logging"`
	Metrics    MetricsConfig    `yaml:"metrics"`
	Health     HealthConfig     `yaml:"health"`
	Plugins    PluginsConfig    `yaml:"plugins"`
	Executor   ExecutorConfig   `yaml:"executor"`
}

// AgentConfig contains agent-specific settings
type AgentConfig struct {
	ID          string            `yaml:"id"`
	Name        string            `yaml:"name"`
	Environment string            `yaml:"environment"`
	Region      string            `yaml:"region"`
	Zone        string            `yaml:"zone"`
	Tags        map[string]string `yaml:"tags"`
	Capabilities []string         `yaml:"capabilities"`
}

// MasterConfig contains master server connection settings
type MasterConfig struct {
	URL                string        `yaml:"url"`
	Token              string        `yaml:"token"`
	ConnectTimeout     time.Duration `yaml:"connect_timeout"`
	HeartbeatInterval  time.Duration `yaml:"heartbeat_interval"`
	ReconnectInterval  time.Duration `yaml:"reconnect_interval"`
	MaxReconnectAttempts int         `yaml:"max_reconnect_attempts"`
	TLSSkipVerify      bool          `yaml:"tls_skip_verify"`
}

// APIConfig contains API server settings
type APIConfig struct {
	HTTP HTTPConfig `yaml:"http"`
	GRPC GRPCConfig `yaml:"grpc"`
}

// HTTPConfig contains HTTP API settings
type HTTPConfig struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	TLS     TLSConfig `yaml:"tls"`
}

// GRPCConfig contains gRPC API settings
type GRPCConfig struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	TLS     TLSConfig `yaml:"tls"`
}

// TLSConfig contains TLS settings
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
	CAFile   string `yaml:"ca_file"`
}

// SecurityConfig contains security settings
type SecurityConfig struct {
	JWT      JWTConfig      `yaml:"jwt"`
	RBAC     RBACConfig     `yaml:"rbac"`
	Audit    AuditConfig    `yaml:"audit"`
	Firewall FirewallConfig `yaml:"firewall"`
}

// JWTConfig contains JWT authentication settings
type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	Expiration time.Duration `yaml:"expiration"`
	Issuer     string        `yaml:"issuer"`
}

// RBACConfig contains role-based access control settings
type RBACConfig struct {
	Enabled     bool   `yaml:"enabled"`
	PolicyFile  string `yaml:"policy_file"`
	DefaultRole string `yaml:"default_role"`
}

// AuditConfig contains audit logging settings
type AuditConfig struct {
	Enabled    bool   `yaml:"enabled"`
	LogFile    string `yaml:"log_file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// FirewallConfig contains firewall settings
type FirewallConfig struct {
	Enabled       bool     `yaml:"enabled"`
	AllowedIPs    []string `yaml:"allowed_ips"`
	BlockedIPs    []string `yaml:"blocked_ips"`
	RateLimiting  bool     `yaml:"rate_limiting"`
	RequestsPerMin int     `yaml:"requests_per_min"`
}

// StorageConfig contains storage settings
type StorageConfig struct {
	DataDir     string `yaml:"data_dir"`
	TempDir     string `yaml:"temp_dir"`
	MaxFileSize int64  `yaml:"max_file_size"`
	Cleanup     CleanupConfig `yaml:"cleanup"`
}

// CleanupConfig contains cleanup settings
type CleanupConfig struct {
	Enabled      bool          `yaml:"enabled"`
	Interval     time.Duration `yaml:"interval"`
	MaxAge       time.Duration `yaml:"max_age"`
	MaxDiskUsage int           `yaml:"max_disk_usage"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// MetricsConfig contains metrics collection settings
type MetricsConfig struct {
	Enabled    bool          `yaml:"enabled"`
	Address    string        `yaml:"address"`
	Port       int           `yaml:"port"`
	Path       string        `yaml:"path"`
	Interval   time.Duration `yaml:"interval"`
	Collectors []string      `yaml:"collectors"`
}

// HealthConfig contains health check settings
type HealthConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Address  string        `yaml:"address"`
	Port     int           `yaml:"port"`
	Path     string        `yaml:"path"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

// PluginsConfig contains plugin system settings
type PluginsConfig struct {
	Enabled   bool              `yaml:"enabled"`
	Directory string            `yaml:"directory"`
	Plugins   map[string]Plugin `yaml:"plugins"`
}

// Plugin represents a plugin configuration
type Plugin struct {
	Enabled bool                   `yaml:"enabled"`
	Path    string                 `yaml:"path"`
	Config  map[string]interface{} `yaml:"config"`
}

// ExecutorConfig contains task executor settings
type ExecutorConfig struct {
	MaxConcurrentTasks int           `yaml:"max_concurrent_tasks"`
	TaskTimeout        time.Duration `yaml:"task_timeout"`
	WorkerPoolSize     int           `yaml:"worker_pool_size"`
	QueueSize          int           `yaml:"queue_size"`
}

// Load loads configuration from file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables
	data = []byte(os.ExpandEnv(string(data)))

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	config.setDefaults()

	return &config, nil
}

// setDefaults sets default values for configuration
func (c *Config) setDefaults() {
	// Agent defaults
	if c.Agent.ID == "" {
		hostname, _ := os.Hostname()
		c.Agent.ID = hostname
	}
	if c.Agent.Name == "" {
		c.Agent.Name = c.Agent.ID
	}
	if c.Agent.Environment == "" {
		c.Agent.Environment = "production"
	}

	// Master defaults
	if c.Master.ConnectTimeout == 0 {
		c.Master.ConnectTimeout = 30 * time.Second
	}
	if c.Master.HeartbeatInterval == 0 {
		c.Master.HeartbeatInterval = 30 * time.Second
	}
	if c.Master.ReconnectInterval == 0 {
		c.Master.ReconnectInterval = 10 * time.Second
	}
	if c.Master.MaxReconnectAttempts == 0 {
		c.Master.MaxReconnectAttempts = 10
	}

	// API defaults
	if c.API.HTTP.Address == "" {
		c.API.HTTP.Address = "0.0.0.0"
	}
	if c.API.HTTP.Port == 0 {
		c.API.HTTP.Port = 8080
	}
	if c.API.GRPC.Address == "" {
		c.API.GRPC.Address = "0.0.0.0"
	}
	if c.API.GRPC.Port == 0 {
		c.API.GRPC.Port = 8443
	}

	// Storage defaults
	if c.Storage.DataDir == "" {
		c.Storage.DataDir = "/opt/ducla/data"
	}
	if c.Storage.TempDir == "" {
		c.Storage.TempDir = "/tmp/ducla"
	}
	if c.Storage.MaxFileSize == 0 {
		c.Storage.MaxFileSize = 100 * 1024 * 1024 // 100MB
	}

	// Logging defaults
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.Logging.Format == "" {
		c.Logging.Format = "json"
	}
	if c.Logging.Output == "" {
		c.Logging.Output = "stdout"
	}

	// Metrics defaults
	if c.Metrics.Address == "" {
		c.Metrics.Address = "0.0.0.0"
	}
	if c.Metrics.Port == 0 {
		c.Metrics.Port = 9090
	}
	if c.Metrics.Path == "" {
		c.Metrics.Path = "/metrics"
	}
	if c.Metrics.Interval == 0 {
		c.Metrics.Interval = 15 * time.Second
	}

	// Health defaults
	if c.Health.Address == "" {
		c.Health.Address = "0.0.0.0"
	}
	if c.Health.Port == 0 {
		c.Health.Port = 8081
	}
	if c.Health.Path == "" {
		c.Health.Path = "/health"
	}
	if c.Health.Interval == 0 {
		c.Health.Interval = 30 * time.Second
	}
	if c.Health.Timeout == 0 {
		c.Health.Timeout = 5 * time.Second
	}

	// Executor defaults
	if c.Executor.MaxConcurrentTasks == 0 {
		c.Executor.MaxConcurrentTasks = 10
	}
	if c.Executor.TaskTimeout == 0 {
		c.Executor.TaskTimeout = 30 * time.Minute
	}
	if c.Executor.WorkerPoolSize == 0 {
		c.Executor.WorkerPoolSize = 5
	}
	if c.Executor.QueueSize == 0 {
		c.Executor.QueueSize = 100
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Master.URL == "" {
		return fmt.Errorf("master.url is required")
	}
	if c.Master.Token == "" {
		return fmt.Errorf("master.token is required")
	}
	if c.Agent.ID == "" {
		return fmt.Errorf("agent.id is required")
	}
	return nil
}