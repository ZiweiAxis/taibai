package taibai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing
type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Response, nil
}

// Create mock response
func newMockResponse(status int, body interface{}) *http.Response {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		Header:     http.Header{},
	}
}

func TestNewClient(t *testing.T) {
	config := &Config{
		ServerAddress: "localhost:8008",
		Token:         "test-token",
		Timeout:       30 * time.Second,
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("Expected client, got nil")
	}

	if client.baseURL != "http://localhost:8008" {
		t.Errorf("Expected baseURL to be 'http://localhost:8008', got '%s'", client.baseURL)
	}

	if client.token != "test-token" {
		t.Errorf("Expected token to be 'test-token', got '%s'", client.token)
	}
}

func TestNewClientWithHTTPS(t *testing.T) {
	config := &Config{
		ServerAddress: "https://localhost:8008",
		Token:         "test-token",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client.baseURL != "https://localhost:8008" {
		t.Errorf("Expected baseURL to be 'https://localhost:8008', got '%s'", client.baseURL)
	}
}

func TestNewClientTrailingSlash(t *testing.T) {
	config := &Config{
		ServerAddress: "http://localhost:8008/",
		Token:         "test-token",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client.baseURL != "http://localhost:8008" {
		t.Errorf("Expected baseURL to be 'http://localhost:8008', got '%s'", client.baseURL)
	}
}

func TestNewClientInvalidConfig(t *testing.T) {
	config := &Config{
		ServerAddress: "",
		Token:         "test-token",
	}

	_, err := NewClient(config)
	if err == nil {
		t.Fatal("Expected error for invalid config, got nil")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				ServerAddress: "localhost:8008",
				Token:         "test-token",
			},
			wantErr: false,
		},
		{
			name: "invalid config - empty server address",
			config: &Config{
				ServerAddress: "",
				Token:         "test-token",
			},
			wantErr: true,
		},
		{
			name: "invalid config - zero timeout",
			config: &Config{
				ServerAddress: "localhost:8008",
				Token:         "test-token",
				Timeout:       0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigDefaultValues(t *testing.T) {
	config := &Config{
		ServerAddress: "localhost:8008",
		Token:         "test-token",
	}

	err := config.Validate()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Timeout <= 0 {
		t.Errorf("Expected Timeout to be set, got %v", config.Timeout)
	}

	if config.MaxIdleConnections <= 0 {
		t.Errorf("Expected MaxIdleConnections to be set, got %v", config.MaxIdleConnections)
	}

	if config.IdleConnTimeout <= 0 {
		t.Errorf("Expected IdleConnTimeout to be set, got %v", config.IdleConnTimeout)
	}
}

func TestClientToken(t *testing.T) {
	config := &Config{
		ServerAddress: "localhost:8008",
		Token:         "test-token",
	}

	client, _ := NewClient(config)

	// Test GetToken
	if client.GetToken() != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", client.GetToken())
	}

	// Test SetToken
	client.SetToken("new-token")
	if client.GetToken() != "new-token" {
		t.Errorf("Expected token 'new-token', got '%s'", client.GetToken())
	}
}

func TestClientSafeClient(t *testing.T) {
	config := &Config{
		ServerAddress: "localhost:8008",
		Token:         "test-token",
	}

	safeClient, err := NewSafeClient(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if safeClient == nil {
		t.Fatal("Expected SafeClient, got nil")
	}

	// Test SetToken
	safeClient.SetToken("new-token")
}

func TestClientClose(t *testing.T) {
	config := &Config{
		ServerAddress: "localhost:8008",
		Token:         "test-token",
	}

	client, _ := NewClient(config)

	// Should not error
	err := client.Close()
	if err != nil {
		t.Errorf("Expected no error on Close, got %v", err)
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		errResp  ErrorResponse
		expected string
	}{
		{
			name: "message field",
			errResp: ErrorResponse{
				Message: "test message",
			},
			expected: "test message",
		},
		{
			name: "error field",
			errResp: ErrorResponse{
				ErrorMsg: "test error",
			},
			expected: "test error",
		},
		{
			name: "empty",
			errResp: ErrorResponse{},
			expected: "unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errResp.Error()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestAPIError(t *testing.T) {
	err := &APIError{
		Code:    404,
		Message: "not found",
	}

	if err.Error() != "not found" {
		t.Errorf("Expected 'not found', got '%s'", err.Error())
	}

	if err.Code != 404 {
		t.Errorf("Expected code 404, got %d", err.Code)
	}
}

func TestConfigError(t *testing.T) {
	err := &ConfigError{
		msg: "test error",
	}

	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", err.Error())
	}
}

// HTTPClient interface tests - verify interface implementation
var _ HTTPClient = (*MockHTTPClient)(nil)

func TestHTTPClientInterface(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{"status": "ok"}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
	}

	ctx := context.Background()
	_, err := client.do(ctx, &Request{
		Method: "GET",
		Path:   "/test",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClientDoBuildURL(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:     "test-token",
	}

	ctx := context.Background()

	// Test with query parameters
	_, err := client.do(ctx, &Request{
		Method: "GET",
		Path:   "/test",
		Query: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClientDoWithBody(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:     "test-token",
	}

	ctx := context.Background()

	// Test with body
	_, err := client.do(ctx, &Request{
		Method: "POST",
		Path:   "/test",
		Body:   map[string]string{"key": "value"},
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClientDoAPIError(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(400, ErrorResponse{
			Code:    400,
			Message: "bad request",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:     "test-token",
	}

	ctx := context.Background()

	_, err := client.do(ctx, &Request{
		Method: "GET",
		Path:   "/test",
	})

	if err == nil {
		t.Fatal("Expected error for API error response")
	}
}

func TestClientDoHTTPError(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(500, "internal server error"),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:     "test-token",
	}

	ctx := context.Background()

	_, err := client.do(ctx, &Request{
		Method: "GET",
		Path:   "/test",
	})

	if err == nil {
		t.Fatal("Expected error for HTTP error response")
	}
}

func TestClientDoJSON(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{"result": "ok"}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:     "test-token",
	}

	ctx := context.Background()

	type Result struct {
		Result string `json:"result"`
	}

	var result Result
	err := client.doJSON(ctx, &Request{
		Method: "GET",
		Path:   "/test",
	}, &result)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Result != "ok" {
		t.Errorf("Expected result 'ok', got '%s'", result.Result)
	}
}

func TestClientGET(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{"result": "ok"}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:     "test-token",
	}

	ctx := context.Background()

	type Result struct {
		Result string `json:"result"`
	}

	var result Result
	err := client.GET(ctx, "/test", map[string]string{"key": "value"}, &result)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClientPOST(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{"result": "ok"}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:     "test-token",
	}

	ctx := context.Background()

	type Result struct {
		Result string `json:"result"`
	}

	var result Result
	err := client.POST(ctx, "/test", map[string]string{"key": "value"}, &result)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClientPUT(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{"result": "ok"}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:     "test-token",
	}

	ctx := context.Background()

	type Result struct {
		Result string `json:"result"`
	}

	var result Result
	err := client.PUT(ctx, "/test", map[string]string{"key": "value"}, &result)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClientDELETE(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{"result": "ok"}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:     "test-token",
	}

	ctx := context.Background()

	type Result struct {
		Result string `json:"result"`
	}

	var result Result
	err := client.DELETE(ctx, "/test", map[string]string{"key": "value"}, &result)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestClientRequestHeaders(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:     "test-token",
	}

	ctx := context.Background()

	_, err := client.do(ctx, &Request{
		Method: "GET",
		Path:   "/test",
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
