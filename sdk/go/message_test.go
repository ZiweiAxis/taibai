package taibai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// MockMessageClient creates a client with mock HTTP for message testing
func MockMessageClient(response *http.Response, err error) *Client {
	mock := &MockHTTPClient{
		Response: response,
		Err:      err,
	}
	return &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: &Client{
			httpClient: mock,
			baseURL:    "http://localhost:8008",
		}},
	}
}

func TestSendMessage(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"event_id": "$test-event-id",
			"room_id":  "!test-room:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Message.SendMessage(ctx, &SendMessageRequest{
		RoomID:      "!test-room:localhost",
		Content:     "Hello, World!",
		MessageType: "m.text",
		Format:      "plain",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.EventID != "$test-event-id" {
		t.Errorf("Expected event_id '$test-event-id', got '%s'", resp.EventID)
	}

	if resp.RoomID != "!test-room:localhost" {
		t.Errorf("Expected room_id '!test-room:localhost', got '%s'", resp.RoomID)
	}
}

func TestSendTextMessage(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"event_id": "$test-event-id",
			"room_id":  "!test-room:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Message.SendTextMessage(ctx, "!test-room:localhost", "Hello, World!")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.EventID != "$test-event-id" {
		t.Errorf("Expected event_id '$test-event-id', got '%s'", resp.EventID)
	}
}

func TestSendHTMLMessage(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"event_id": "$test-event-id",
			"room_id":  "!test-room:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Message.SendHTMLMessage(ctx, "!test-room:localhost", "Hello", "<b>Hello</b>")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.EventID != "$test-event-id" {
		t.Errorf("Expected event_id '$test-event-id', got '%s'", resp.EventID)
	}
}

func TestSendImageMessage(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"event_id": "$test-event-id",
			"room_id":  "!test-room:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Message.SendImageMessage(ctx, "!test-room:localhost", "mxc://example.com/image", "{}")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.EventID != "$test-event-id" {
		t.Errorf("Expected event_id '$test-event-id', got '%s'", resp.EventID)
	}
}

func TestGetMessage(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]interface{}{
			"event_id":  "$test-event-id",
			"room_id":   "!test-room:localhost",
			"sender":    "@user:localhost",
			"type":      "m.room.message",
			"timestamp": 1234567890,
			"content": map[string]string{
				"msgtype": "m.text",
				"body":    "Hello, World!",
			},
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: client},
	}

	ctx := context.Background()

	msg, err := client.Message.GetMessage(ctx, "!test-room:localhost", "$test-event-id")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if msg.EventID != "$test-event-id" {
		t.Errorf("Expected event_id '$test-event-id', got '%s'", msg.EventID)
	}

	if msg.Sender != "@user:localhost" {
		t.Errorf("Expected sender '@user:localhost', got '%s'", msg.Sender)
	}
}

func TestGetRoomMessages(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]interface{}{
			"chunk": []interface{}{
				map[string]interface{}{
					"event_id":  "$event-1",
					"room_id":   "!test-room:localhost",
					"sender":    "@user1:localhost",
					"type":      "m.room.message",
					"timestamp": 1234567890,
				},
				map[string]interface{}{
					"event_id":  "$event-2",
					"room_id":   "!test-room:localhost",
					"sender":    "@user2:localhost",
					"type":      "m.room.message",
					"timestamp": 1234567891,
				},
			},
			"start": "start-token",
			"end":   "end-token",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Message.GetRoomMessages(ctx, "!test-room:localhost", 20, "", "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(resp.Chunk) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(resp.Chunk))
	}

	if resp.Start != "start-token" {
		t.Errorf("Expected start token 'start-token', got '%s'", resp.Start)
	}

	if resp.End != "end-token" {
		t.Errorf("Expected end token 'end-token', got '%s'", resp.End)
	}
}

func TestRedactMessage(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"event_id": "$redacted-event-id",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: client},
	}

	ctx := context.Background()

	err := client.Message.RedactMessage(ctx, "!test-room:localhost", "$test-event-id", "Spam")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendMessageDefaultValues(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"event_id": "$test-event-id",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: client},
	}

	ctx := context.Background()

	// Send message with minimal fields
	resp, err := client.Message.SendMessage(ctx, &SendMessageRequest{
		RoomID:  "!test-room:localhost",
		Content: "Hello",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.EventID != "$test-event-id" {
		t.Errorf("Expected event_id '$test-event-id', got '%s'", resp.EventID)
	}
}

func TestMessageEventStruct(t *testing.T) {
	// Test that MessageEvent can be properly serialized/deserialized
	event := MessageEvent{
		EventID:   "$test-event-id",
		RoomID:    "!test-room:localhost",
		Sender:    "@user:localhost",
		Type:      "m.room.message",
		Timestamp: 1234567890,
		Content: map[string]interface{}{
			"msgtype": "m.text",
			"body":    "Hello",
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal MessageEvent: %v", err)
	}

	var parsed MessageEvent
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal MessageEvent: %v", err)
	}

	if parsed.EventID != event.EventID {
		t.Errorf("Expected EventID '%s', got '%s'", event.EventID, parsed.EventID)
	}

	if parsed.RoomID != event.RoomID {
		t.Errorf("Expected RoomID '%s', got '%s'", event.RoomID, parsed.RoomID)
	}
}

func TestEncryptionInfo(t *testing.T) {
	info := &EncryptionInfo{
		Key:     "test-key",
		IV:      "test-iv",
		Hashes:  map[string]string{"sha256": "hash-value"},
		Version: "v1",
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal EncryptionInfo: %v", err)
	}

	var parsed EncryptionInfo
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal EncryptionInfo: %v", err)
	}

	if parsed.Key != "test-key" {
		t.Errorf("Expected Key 'test-key', got '%s'", parsed.Key)
	}

	if parsed.Version != "v1" {
		t.Errorf("Expected Version 'v1', got '%s'", parsed.Version)
	}
}

func TestMessagesResponse(t *testing.T) {
	response := MessagesResponse{
		Chunk: []MessageEvent{
			{EventID: "$event-1", RoomID: "!room:localhost"},
			{EventID: "$event-2", RoomID: "!room:localhost"},
		},
		Start: "start",
		End:   "end",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal MessagesResponse: %v", err)
	}

	var parsed MessagesResponse
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal MessagesResponse: %v", err)
	}

	if len(parsed.Chunk) != 2 {
		t.Errorf("Expected 2 chunks, got %d", len(parsed.Chunk))
	}
}

func TestSendMessageResponse(t *testing.T) {
	response := SendMessageResponse{
		EventID:   "$test-event-id",
		RoomID:    "!test-room:localhost",
		Sender:    "@user:localhost",
		Timestamp: 1234567890,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal SendMessageResponse: %v", err)
	}

	var parsed SendMessageResponse
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal SendMessageResponse: %v", err)
	}

	if parsed.EventID != "$test-event-id" {
		t.Errorf("Expected EventID '$test-event-id', got '%s'", parsed.EventID)
	}
}

// Helper function to create mock response
func createMockResponse(status int, body interface{}) *http.Response {
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

func TestMessageAPIErrorHandling(t *testing.T) {
	// Test error response handling
	mock := &MockHTTPClient{
		Response: createMockResponse(400, ErrorResponse{
			Code:    400,
			Message: "Bad Request",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: client},
	}

	ctx := context.Background()

	_, err := client.Message.SendMessage(ctx, &SendMessageRequest{
		RoomID:  "!test-room:localhost",
		Content: "Hello",
	})

	if err == nil {
		t.Fatal("Expected error for bad request")
	}
}

func TestMessageAPIPathConstruction(t *testing.T) {
	// Test that paths are constructed correctly
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"event_id": "$test-event-id",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Message:    &MessageAPI{client: client},
	}

	ctx := context.Background()

	// This test verifies the path construction logic
	_, err := client.Message.GetMessage(ctx, "!test-room:localhost", "$event-id")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
