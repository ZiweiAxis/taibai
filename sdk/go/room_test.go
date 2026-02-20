package taibai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestCreateRoom(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"room_id": "!test-room:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Room.CreateRoom(ctx, &CreateRoomRequest{
		Name:       "Test Room",
		Topic:      "A test room",
		Visibility: "private",
		Preset:    "private_chat",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.RoomID != "!test-room:localhost" {
		t.Errorf("Expected room_id '!test-room:localhost', got '%s'", resp.RoomID)
	}
}

func TestCreateRoomWithAlias(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"room_id":   "!test-room:localhost",
			"room_alias": "#test-alias:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Room.CreateRoom(ctx, &CreateRoomRequest{
		Name:            "Test Room",
		RoomAliasName:   "test-alias",
		Visibility:      "public",
		Preset:          "public_chat",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.RoomID != "!test-room:localhost" {
		t.Errorf("Expected room_id '!test-room:localhost', got '%s'", resp.RoomID)
	}

	if resp.RoomAlias != "#test-alias:localhost" {
		t.Errorf("Expected room_alias '#test-alias:localhost', got '%s'", resp.RoomAlias)
	}
}

func TestCreatePublicRoom(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"room_id": "!test-room:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Room.CreatePublicRoom(ctx, "Public Room", "A public room", "public-room")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.RoomID != "!test-room:localhost" {
		t.Errorf("Expected room_id '!test-room:localhost', got '%s'", resp.RoomID)
	}
}

func TestCreatePrivateRoom(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"room_id": "!test-room:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Room.CreatePrivateRoom(ctx, "Private Room", []string{"@user1:localhost", "@user2:localhost"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.RoomID != "!test-room:localhost" {
		t.Errorf("Expected room_id '!test-room:localhost', got '%s'", resp.RoomID)
	}
}

func TestJoinRoom(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"room_id": "!test-room:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Room.JoinRoom(ctx, "!test-room:localhost", nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.RoomID != "!test-room:localhost" {
		t.Errorf("Expected room_id '!test-room:localhost', got '%s'", resp.RoomID)
	}
}

func TestJoinRoomWithServerName(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"room_id": "!test-room:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Room.JoinRoom(ctx, "#test-alias:localhost", &JoinRoomRequest{
		ServerName: "localhost",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.RoomID != "!test-room:localhost" {
		t.Errorf("Expected room_id '!test-room:localhost', got '%s'", resp.RoomID)
	}
}

func TestLeaveRoom(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.LeaveRoom(ctx, "!test-room:localhost", nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestLeaveRoomWithReason(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.LeaveRoom(ctx, "!test-room:localhost", &LeaveRoomRequest{
		Reason: "No longer needed",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestInviteUser(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.InviteUser(ctx, "!test-room:localhost", &InviteUserRequest{
		UserID: "@user:localhost",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestInviteUserWithReason(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.InviteUser(ctx, "!test-room:localhost", &InviteUserRequest{
		UserID: "@user:localhost",
		Reason: "Please join our room",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestKickUser(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.KickUser(ctx, "!test-room:localhost", &KickUserRequest{
		UserID: "@user:localhost",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBanUser(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.BanUser(ctx, "!test-room:localhost", &BanUserRequest{
		UserID: "@user:localhost",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUnbanUser(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.UnbanUser(ctx, "!test-room:localhost", &UnbanUserRequest{
		UserID: "@user:localhost",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetRoom(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]interface{}{
			"room_id":      "!test-room:localhost",
			"name":         "Test Room",
			"topic":        "A test room",
			"num_joined_members": 5,
			"join_rule":   "public",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	room, err := client.Room.GetRoom(ctx, "!test-room:localhost")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if room.RoomID != "!test-room:localhost" {
		t.Errorf("Expected room_id '!test-room:localhost', got '%s'", room.RoomID)
	}

	if room.Name != "Test Room" {
		t.Errorf("Expected name 'Test Room', got '%s'", room.Name)
	}

	if room.MemberCount != 5 {
		t.Errorf("Expected member count 5, got %d", room.MemberCount)
	}
}

func TestGetRoomState(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"name": "Test Room",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	_, err := client.Room.GetRoomState(ctx, "!test-room:localhost", "m.room.name", "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetRoomMembers(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]interface{}{
			"chunk": []interface{}{
				map[string]interface{}{
					"event_id":   "$event-1",
					"room_id":    "!test-room:localhost",
					"type":       "m.room.member",
					"sender":     "@user1:localhost",
					"state_key":  "@user1:localhost",
					"content": map[string]string{
						"membership": "join",
					},
				},
				map[string]interface{}{
					"event_id":   "$event-2",
					"room_id":    "!test-room:localhost",
					"type":       "m.room.member",
					"sender":     "@user2:localhost",
					"state_key":  "@user2:localhost",
					"content": map[string]string{
						"membership": "join",
					},
				},
			},
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Room.GetRoomMembers(ctx, "!test-room:localhost", "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(resp.Chunk) != 2 {
		t.Errorf("Expected 2 members, got %d", len(resp.Chunk))
	}
}

func TestSetRoomName(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.SetRoomName(ctx, "!test-room:localhost", "New Room Name")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSetRoomTopic(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.SetRoomTopic(ctx, "!test-room:localhost", "New Topic")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSetRoomAvatar(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.SetRoomAvatar(ctx, "!test-room:localhost", "mxc://example.com/avatar")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetJoinedRooms(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]interface{}{
			"joined_rooms": []string{
				"!room1:localhost",
				"!room2:localhost",
			},
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Room.GetJoinedRooms(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(resp.JoinedRooms) != 2 {
		t.Errorf("Expected 2 rooms, got %d", len(resp.JoinedRooms))
	}
}

func TestGetRoomPowerLevels(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]interface{}{
			"users": map[string]int{
				"@admin:localhost": 100,
				"@user:localhost":   50,
			},
			"users_default":    0,
			"events":           map[string]int{},
			"events_default":   0,
			"state_default":    50,
			"ban":              50,
			"kick":             50,
			"redact":           50,
			"invite":           0,
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	levels, err := client.Room.GetRoomPowerLevels(ctx, "!test-room:localhost")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if levels.Users["@admin:localhost"] != 100 {
		t.Errorf("Expected admin power level 100, got %d", levels.Users["@admin:localhost"])
	}
}

func TestSetRoomPowerLevels(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.SetRoomPowerLevels(ctx, "!test-room:localhost", &PowerLevels{
		Users: map[string]int{
			"@admin:localhost": 100,
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetRoomAliases(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]interface{}{
			"aliases": []string{
				"#test-alias:localhost",
				"#another-alias:localhost",
			},
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	resp, err := client.Room.GetRoomAliases(ctx, "!test-room:localhost")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(resp.Aliases) != 2 {
		t.Errorf("Expected 2 aliases, got %d", len(resp.Aliases))
	}
}

func TestForgetRoom(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.ForgetRoom(ctx, "!test-room:localhost")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteRoom(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, nil),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	err := client.Room.DeleteRoom(ctx, "!test-room:localhost", true)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateRoomDefaultValues(t *testing.T) {
	mock := &MockHTTPClient{
		Response: newMockResponse(200, map[string]string{
			"room_id": "!test-room:localhost",
		}),
	}

	client := &Client{
		httpClient: mock,
		baseURL:    "http://localhost:8008",
		token:      "test-token",
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	// Create room with nil request
	resp, err := client.Room.CreateRoom(ctx, nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.RoomID != "!test-room:localhost" {
		t.Errorf("Expected room_id '!test-room:localhost', got '%s'", resp.RoomID)
	}
}

func TestRoomStruct(t *testing.T) {
	room := Room{
		RoomID:         "!test-room:localhost",
		Name:           "Test Room",
		Topic:          "A test room",
		JoinRule:       "public",
		GuestCanJoin:   false,
		MemberCount:   5,
	}

	data, err := json.Marshal(room)
	if err != nil {
		t.Fatalf("Failed to marshal Room: %v", err)
	}

	var parsed Room
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal Room: %v", err)
	}

	if parsed.RoomID != room.RoomID {
		t.Errorf("Expected RoomID '%s', got '%s'", room.RoomID, parsed.RoomID)
	}
}

func TestCreateRoomRequest(t *testing.T) {
	req := CreateRoomRequest{
		Name:            "Test Room",
		Topic:           "A test room",
		RoomAliasName:   "test-room",
		Visibility:      "public",
		Preset:          "public_chat",
		Invite:          []string{"@user1:localhost", "@user2:localhost"},
		IsDirect:        false,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal CreateRoomRequest: %v", err)
	}

	var parsed CreateRoomRequest
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal CreateRoomRequest: %v", err)
	}

	if parsed.Name != req.Name {
		t.Errorf("Expected Name '%s', got '%s'", req.Name, parsed.Name)
	}

	if len(parsed.Invite) != 2 {
		t.Errorf("Expected 2 invites, got %d", len(parsed.Invite))
	}
}

func TestPowerLevels(t *testing.T) {
	levels := PowerLevels{
		Users:           map[string]int{"@admin:localhost": 100},
		UsersDefault:    0,
		Events:          map[string]int{},
		EventsDefault:   0,
		StateDefault:    50,
		Ban:             50,
		Kick:            50,
		Redact:          50,
		Invite:          0,
	}

	data, err := json.Marshal(levels)
	if err != nil {
		t.Fatalf("Failed to marshal PowerLevels: %v", err)
	}

	var parsed PowerLevels
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal PowerLevels: %v", err)
	}

	if parsed.Users["@admin:localhost"] != 100 {
		t.Errorf("Expected admin power 100, got %d", parsed.Users["@admin:localhost"])
	}
}

func TestMemberContent(t *testing.T) {
	content := MemberContent{
		Membership:  "join",
		DisplayName: "Test User",
		AvatarURL:   "mxc://example.com/avatar",
		Reason:      "",
	}

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("Failed to marshal MemberContent: %v", err)
	}

	var parsed MemberContent
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("Failed to unmarshal MemberContent: %v", err)
	}

	if parsed.Membership != "join" {
		t.Errorf("Expected membership 'join', got '%s'", parsed.Membership)
	}
}

func TestRoomAPIErrorHandling(t *testing.T) {
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
		Room:       &RoomAPI{client: client},
	}

	ctx := context.Background()

	_, err := client.Room.CreateRoom(ctx, &CreateRoomRequest{
		Name: "Test Room",
	})

	if err == nil {
		t.Fatal("Expected error for bad request")
	}
}
