package taibai

import (
	"context"
	"time"
)

// MessageAPI handles message-related operations
type MessageAPI struct {
	client *Client
}

// SendMessageRequest represents a message to be sent
type SendMessageRequest struct {
	// RoomID is the room where the message will be sent
	RoomID string `json:"room_id"`

	// Content is the message content
	Content string `json:"content"`

	// MessageType is the type of message (e.g., "text", "html", "m.image")
	MessageType string `json:"msgtype,omitempty"`

	// Format is the format of the message (e.g., "plain", "html")
	Format string `json:"format,omitempty"`

	// Body is the alternative plain text body
	Body string `json:"body,omitempty"`

	// URL is the URL for media messages
	URL string `json:"url,omitempty"`

	// Info is the encryption info for media
	Info *EncryptionInfo `json:"info,omitempty"`

	// Sender is the sender of the message (optional, defaults to authenticated user)
	Sender string `json:"sender,omitempty"`

	// Timestamp is the timestamp of the message (optional)
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// EncryptionInfo contains encryption details for media
type EncryptionInfo struct {
	// Key is the encryption key
	Key string `json:"key,omitempty"`

	// IV is the initialization vector
	IV string `json:"iv,omitempty"`

	// Hashes contains the hashes of the encrypted data
	Hashes map[string]string `json:"hashes,omitempty"`

	// Version is the encryption version
	Version string `json:"v,omitempty"`
}

// SendMessageResponse represents the response from sending a message
type SendMessageResponse struct {
	// EventID is the unique identifier of the sent message
	EventID string `json:"event_id"`

	// RoomID is the room where the message was sent
	RoomID string `json:"room_id,omitempty"`

	// Sender is the sender of the message
	Sender string `json:"sender,omitempty"`

	// Timestamp is the timestamp of the message
	Timestamp int64 `json:"timestamp,omitempty"`
}

// MessageEvent represents a message event
type MessageEvent struct {
	// EventID is the unique identifier of the event
	EventID string `json:"event_id"`

	// RoomID is the room where the event occurred
	RoomID string `json:"room_id"`

	// Sender is the sender of the event
	Sender string `json:"sender"`

	// Type is the type of the event
	Type string `json:"type"`

	// Timestamp is the timestamp of the event
	Timestamp int64 `json:"timestamp"`

	// Content is the content of the message
	Content map[string]interface{} `json:"content,omitempty"`
}

// SendMessage sends a message to a room
func (m *MessageAPI) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	// Set default message type
	if req.MessageType == "" {
		req.MessageType = "m.text"
	}

	// Set default format
	if req.Format == "" {
		req.Format = "plain"
	}

	// Set body if not provided
	if req.Body == "" {
		req.Body = req.Content
	}

	result := &SendMessageResponse{}
	err := m.client.POST(ctx, "/_matrix/client/r0/rooms/"+req.RoomID+"/send/m.room.message", req, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SendTextMessage sends a plain text message to a room
func (m *MessageAPI) SendTextMessage(ctx context.Context, roomID, content string) (*SendMessageResponse, error) {
	return m.SendMessage(ctx, &SendMessageRequest{
		RoomID:      roomID,
		Content:     content,
		MessageType: "m.text",
		Format:      "plain",
	})
}

// SendHTMLMessage sends an HTML message to a room
func (m *MessageAPI) SendHTMLMessage(ctx context.Context, roomID, content, html string) (*SendMessageResponse, error) {
	return m.SendMessage(ctx, &SendMessageRequest{
		RoomID:      roomID,
		Content:     content,
		Body:        content,
		Format:      "html",
		MessageType: "m.text",
	})
}

// SendImageMessage sends an image message to a room
func (m *MessageAPI) SendImageMessage(ctx context.Context, roomID, url, info string) (*SendMessageResponse, error) {
	return m.SendMessage(ctx, &SendMessageRequest{
		RoomID:      roomID,
		URL:         url,
		MessageType: "m.image",
		Info: &EncryptionInfo{
			Version: "v1",
		},
	})
}

// GetMessage retrieves a specific message from a room
func (m *MessageAPI) GetMessage(ctx context.Context, roomID, eventID string) (*MessageEvent, error) {
	path := "/_matrix/client/r0/rooms/" + roomID + "/event/" + eventID
	result := &MessageEvent{}
	err := m.client.GET(ctx, path, nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetRoomMessages retrieves messages from a room
func (m *MessageAPI) GetRoomMessages(ctx context.Context, roomID string, limit int, from, to string) (*MessagesResponse, error) {
	query := map[string]string{
		"limit":  "20",
		"dir":    "b",
	}
	if limit > 0 {
		query["limit"] = string(rune(limit))
	}
	if from != "" {
		query["from"] = from
	}
	if to != "" {
		query["to"] = to
	}

	result := &MessagesResponse{}
	err := m.client.GET(ctx, "/_matrix/client/r0/rooms/"+roomID+"/messages", query, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MessagesResponse represents a paginated list of messages
type MessagesResponse struct {
	// Chunk contains the message events
	Chunk []MessageEvent `json:"chunk"`

	// Start is the token for the start of the chunk
	Start string `json:"start"`

	// End is the token for the end of the chunk
	End string `json:"end"`

	// State contains the state events at the start of the chunk
	State []MessageEvent `json:"state,omitempty"`
}

// RedactMessage redacts a message in a room
func (m *MessageAPI) RedactMessage(ctx context.Context, roomID, eventID string, reason string) error {
	path := "/_matrix/client/r0/rooms/" + roomID + "/redact/" + eventID

	body := map[string]string{}
	if reason != "" {
		body["reason"] = reason
	}

	return m.client.PUT(ctx, path, body, nil)
}
