package config

import "time"

// Config holds runtime settings for sdk, server, and sidecar modes.
type Config struct {
	Mode                 string
	PolicyBackend        PolicyBackendConfig
	CounterBackend       CounterBackendConfig
	CacheTTL             time.Duration
	RefreshInterval      time.Duration
	DefaultFailureMode   string
	MetricsEnabled       bool
	ShadowMetricsEnabled bool
	DecisionLogEnabled   bool
}

// PolicyBackendConfig selects and configures the policy storage backend.
type PolicyBackendConfig struct {
	Driver           string
	DSN              string
	TableName        string
	UseLegacyAdapter bool
}

// CounterBackendConfig selects and configures the hot path counter backend.
type CounterBackendConfig struct {
	Driver   string
	Address  string
	Password string
	DB       int
	Prefix   string
}
