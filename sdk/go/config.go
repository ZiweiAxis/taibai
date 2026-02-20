package taibai

import (
	"time"
)

// Config holds the SDK configuration
type Config struct {
	// ServerAddress is the base URL of the Taibai API server
	ServerAddress string

	// Token is the authentication token
	Token string

	// Timeout for HTTP requests (default: 30 seconds)
	Timeout time.Duration

	// Max max idle connections in the connection pool (default: 10)
	MaxIdleConnections int

	// IdleConnTimeout timeout for idle connections (default: 90 seconds)
	IdleConnTimeout time.Duration

	// TLSConfig TLS configuration (optional)
	// TLSConfig *tls.Config
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		Timeout:            30 * time.Second,
		MaxIdleConnections: 10,
		IdleConnTimeout:    90 * time.Second,
	}
}

// Validate checks if the config is valid
func (c *Config) Validate() error {
	if c.ServerAddress == "" {
		return ErrInvalidServerAddress
	}
	if c.Timeout <= 0 {
		c.Timeout = 30 * time.Second
	}
	if c.MaxIdleConnections <= 0 {
		c.MaxIdleConnections = 10
	}
	if c.IdleConnTimeout <= 0 {
		c.IdleConnTimeout = 90 * time.Second
	}
	return nil
}

// Error messages
var (
	ErrInvalidServerAddress = &ConfigError{"invalid server address: server address is required"}
	ErrUnauthorized         = &APIError{Code: 401, Message: "unauthorized"}
	ErrForbidden            = &APIError{Code: 403, Message: "forbidden"}
	ErrNotFound             = &APIError{Code: 404, Message: "not found"}
	ErrInternalServerError = &APIError{Code: 500, Message: "internal server error"}
)

type ConfigError struct {
	msg string
}

func (e *ConfigError) Error() string {
	return e.msg
}

// APIError represents an API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}
