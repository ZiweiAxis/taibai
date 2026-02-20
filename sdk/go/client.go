package taibai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the Taibai HTTP client
type Client struct {
	config     *Config
	httpClient HTTPClient
	baseURL    string
	token      string

	// APIs
	Message *MessageAPI
	Room    *RoomAPI
	User    *UserAPI
	Approval *ApprovalAPI
}

// NewClient creates a new Taibai client
func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); nil != err {
		return nil, err
	}

	// Ensure base URL has scheme
	baseURL := config.ServerAddress
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}

	// Remove trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Create HTTP client with connection pool
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConnections,
		IdleConnTimeout:     config.IdleConnTimeout,
		MaxIdleConnsPerHost: config.MaxIdleConnections,
	}

	httpClient := &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}

	client := &Client{
		config:     config,
		httpClient: httpClient,
		baseURL:    baseURL,
		token:      config.Token,
	}

	// Initialize APIs
	client.Message = &MessageAPI{client: client}
	client.Room = &RoomAPI{client: client}
	client.User = &UserAPI{client: client}
	client.Approval = &ApprovalAPI{client: client}

	return client, nil
}

// Request represents an API request
type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Query   map[string]string
	Headers map[string]string
}

// Response represents an API response
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (e *ErrorResponse) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Error != "" {
		return e.Error
	}
	return "unknown error"
}

// do performs an HTTP request
func (c *Client) do(ctx context.Context, req *Request) (*Response, error) {
	// Build URL
	url := c.baseURL + req.Path

	// Add query parameters
	if len(req.Query) > 0 {
		query := ""
		for key, value := range req.Query {
			if query != "" {
				query += "&"
			}
			query += fmt.Sprintf("%s=%s", key, value)
		}
		url += "?" + query
	}

	// Marshal body
	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Add authentication token
	if c.token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.token)
	}

	// Add custom headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Perform request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for API error
	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
		}
		return nil, &APIError{Code: resp.StatusCode, Message: errResp.Error()}
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

// doJSON performs an HTTP request and unmarshals the response
func (c *Client) doJSON(ctx context.Context, req *Request, result interface{}) error {
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}

	if result != nil && len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// GET performs a GET request
func (c *Client) GET(ctx context.Context, path string, query map[string]string, result interface{}) error {
	return c.doJSON(ctx, &Request{
		Method: http.MethodGet,
		Path:   path,
		Query:  query,
	}, result)
}

// POST performs a POST request
func (c *Client) POST(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.doJSON(ctx, &Request{
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	}, result)
}

// PUT performs a PUT request
func (c *Client) PUT(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.doJSON(ctx, &Request{
		Method: http.MethodPut,
		Path:   path,
		Body:   body,
	}, result)
}

// DELETE performs a DELETE request
func (c *Client) DELETE(ctx context.Context, path string, query map[string]string, result interface{}) error {
	return c.doJSON(ctx, &Request{
		Method: http.MethodDelete,
		Path:   path,
		Query:  query,
	}, result)
}

// SetToken sets the authentication token
func (c *Client) SetToken(token string) {
	c.token = token
}

// GetToken gets the current authentication token
func (c *Client) GetToken() string {
	return c.token
}

// Close closes the client and releases resources
func (c *Client) Close() error {
	if closer, ok := c.httpClient.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

// SafeClient is a thread-safe wrapper around Client
type SafeClient struct {
	client *Client
	mu     sync.RWMutex
}

// NewSafeClient creates a new thread-safe Taibai client
func NewSafeClient(config *Config) (*SafeClient, error) {
	client, err := NewClient(config)
	if err != nil {
		return nil, err
	}
	return &SafeClient{client: client}, nil
}

// Do performs a thread-safe request
func (sc *SafeClient) Do(ctx context.Context, req *Request) (*Response, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.client.do(ctx, req)
}

// SetToken sets the authentication token thread-safely
func (sc *SafeClient) SetToken(token string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.client.token = token
}

// Close closes the client thread-safely
func (sc *SafeClient) Close() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.client.Close()
}
