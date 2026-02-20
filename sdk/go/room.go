package taibai

import (
	"context"
	"time"
)

// RoomAPI handles room-related operations
type RoomAPI struct {
	client *Client
}

// CreateRoomRequest represents a request to create a room
type CreateRoomRequest struct {
	// Name is the name of the room
	Name string `json:"name,omitempty"`

	// Topic is the topic of the room
	Topic string `json:"topic,omitempty"`

	// RoomAliasName is the alias of the room (e.g., "my-room")
	RoomAliasName string `json:"room_alias_name,omitempty"`

	// Visibility is the visibility of the room ("public" or "private")
	Visibility string `json:"visibility,omitempty"`

	// Preset is the room preset ("private_chat", "public_chat", "trusted_private_chat")
	Preset string `json:"preset,omitempty"`

	// Invite is a list of user IDs to invite
	Invite []string `json:"invite,omitempty"`

	// Invite3PID is a list of third-party invites
	Invite3PID []ThirdPartyInvite `json:"invite_3pid,omitempty"`

	// RoomVersion is the version of the room
	RoomVersion string `json:"room_version,omitempty"`

	// CreationContent contains additional creation content
	CreationContent map[string]interface{} `json:"creation_content,omitempty"`

	// InitialState contains initial state events
	InitialState []StateEvent `json:"initial_state,omitempty"`

	// PowerLevelContentOverride overrides the default power levels
	PowerLevelContentOverride *PowerLevels `json:"power_level_content_override,omitempty"`

	// JoinRule is the join rule of the room ("public", "knock", "invite", "private")
	JoinRule string `json:"join_rule,omitempty"`

	// GuestCanJoin indicates if guests can join
	GuestCanJoin bool `json:"guest_can_join,omitempty"`

	// IsDirect indicates if this is a direct message room
	IsDirect bool `json:"is_direct,omitempty"`
}

// ThirdPartyInvite represents a third-party invite
type ThirdPartyInvite struct {
	// Medium is the medium of the invite (e.g., "email")
	Medium string `json:"medium"`

	// Address is the address of the invite
	Address string `json:"address"`

	// ValidSinceMs is the valid since timestamp in milliseconds
	ValidSinceMs int64 `json:"valid_since_ms,omitempty"`
}

// StateEvent represents a state event
type StateEvent struct {
	// Type is the type of the event
	Type string `json:"type"`

	// StateKey is the key of the state
	StateKey string `json:"state_key,omitempty"`

	// Content is the content of the event
	Content interface{} `json:"content"`

	// RedactIfUntrusted indicates if the event should be redacted if not trusted
	RedactIfUntrusted bool `json:"redact_if_untrusted,omitempty"`

	// VerifiesIdentity indicates if the event verifies identity
	VerifiesIdentity bool `json:"verifies_identity,omitempty"`
}

// PowerLevels represents the power levels in a room
type PowerLevels struct {
	// Users overrides the power levels for specific users
	Users map[string]int `json:"users,omitempty"`

	// UsersDefault is the default power level for users
	UsersDefault int `json:"users_default,omitempty"`

	// Events overrides the power levels for specific events
	Events map[string]int `json:"events,omitempty"`

	// EventsDefault is the default power level for events
	EventsDefault int `json:"events_default,omitempty"`

	// StateDefault is the default power level for state events
	StateDefault int `json:"state_default,omitempty"`

	// Ban is the power level required to ban users
	Ban int `json:"ban,omitempty"`

	// Kick is the power level required to kick users
	Kick int `json:"kick,omitempty"`

	// Redact is the power level required to redact events
	Redact int `json:"redact,omitempty"`

	// Invite is the power level required to invite users
	Invite int `json:"invite,omitempty"`
}

// CreateRoomResponse represents the response from creating a room
type CreateRoomResponse struct {
	// RoomID is the unique identifier of the created room
	RoomID string `json:"room_id"`

	// RoomAlias is the alias of the room (if set)
	RoomAlias string `json:"room_alias,omitempty"`
}

// Room represents a room
type Room struct {
	// RoomID is the unique identifier of the room
	RoomID string `json:"room_id"`

	// Name is the name of the room
	Name string `json:"name,omitempty"`

	// Topic is the topic of the room
	Topic string `json:"topic,omitempty"`

	// AvatarURL is the avatar URL of the room
	AvatarURL string `json:"avatar_url,omitempty"`

	// CanonicalAlias is the canonical alias of the room
	CanonicalAlias string `json:"canonical_alias,omitempty"`

	// JoinRule is the join rule of the room
	JoinRule string `json:"join_rule,omitempty"`

	// GuestCanJoin indicates if guests can join
	GuestCanJoin bool `json:"guest_can_join,omitempty"`

	// MemberCount is the number of members in the room
	MemberCount int `json:"num_joined_members,omitempty"`

	// Topic is the topic of the room
	WorldReadable bool `json:"world_readable,omitempty"`
}

// JoinRoomRequest represents a request to join a room
type JoinRoomRequest struct {
	// ServerName is the server to use to join the room
	ServerName string `json:"server_name,omitempty"`

	// ThirdPartySigned is the third-party signed data
	ThirdPartySigned map[string]string `json:"third_party_signed,omitempty"`
}

// JoinRoom joins a room by ID or alias
func (r *RoomAPI) JoinRoom(ctx context.Context, roomIDOrAlias string, req *JoinRoomRequest) (*JoinRoomResponse, error) {
	if req == nil {
		req = &JoinRoomRequest{}
	}

	result := &JoinRoomResponse{}
	err := r.client.POST(ctx, "/_matrix/client/r0/join/"+roomIDOrAlias, req, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// JoinRoomResponse represents the response from joining a room
type JoinRoomResponse struct {
	// RoomID is the room that was joined
	RoomID string `json:"room_id"`
}

// LeaveRoomRequest represents a request to leave a room
type LeaveRoomRequest struct {
	// Reason is the reason for leaving
	Reason string `json:"reason,omitempty"`
}

// LeaveRoom leaves a room
func (r *RoomAPI) LeaveRoom(ctx context.Context, roomID string, req *LeaveRoomRequest) error {
	if req == nil {
		req = &LeaveRoomRequest{}
	}

	return r.client.POST(ctx, "/_matrix/client/r0/rooms/"+roomID+"/leave", req, nil)
}

// InviteUserRequest represents a request to invite a user to a room
type InviteUserRequest struct {
	// UserID is the user ID to invite
	UserID string `json:"user_id"`

	// Reason is the reason for the invitation
	Reason string `json:"reason,omitempty"`
}

// InviteUser invites a user to a room
func (r *RoomAPI) InviteUser(ctx context.Context, roomID string, req *InviteUserRequest) error {
	return r.client.POST(ctx, "/_matrix/client/r0/rooms/"+roomID+"/invite", req, nil)
}

// KickUserRequest represents a request to kick a user from a room
type KickUserRequest struct {
	// UserID is the user ID to kick
	UserID string `json:"user_id"`

	// Reason is the reason for the kick
	Reason string `json:"reason,omitempty"`
}

// KickUser kicks a user from a room
func (r *RoomAPI) KickUser(ctx context.Context, roomID string, req *KickUserRequest) error {
	return r.client.POST(ctx, "/_matrix/client/r0/rooms/"+roomID+"/kick", req, nil)
}

// BanUserRequest represents a request to ban a user from a room
type BanUserRequest struct {
	// UserID is the user ID to ban
	UserID string `json:"user_id"`

	// Reason is the reason for the ban
	Reason string `json:"reason,omitempty"`
}

// BanUser bans a user from a room
func (r *RoomAPI) BanUser(ctx context.Context, roomID string, req *BanUserRequest) error {
	return r.client.POST(ctx, "/_matrix/client/r0/rooms/"+roomID+"/ban", req, nil)
}

// UnbanUserRequest represents a request to unban a user from a room
type UnbanUserRequest struct {
	// UserID is the user ID to unban
	UserID string `json:"user_id"`

	// Reason is the reason for the unban
	Reason string `json:"reason,omitempty"`
}

// UnbanUser unbans a user from a room
func (r *RoomAPI) UnbanUser(ctx context.Context, roomID string, req *UnbanUserRequest) error {
	return r.client.POST(ctx, "/_matrix/client/r0/rooms/"+roomID+"/unban", req, nil)
}

// GetRoomState gets the state of a room
func (r *RoomAPI) GetRoomState(ctx context.Context, roomID, eventType, stateKey string) (interface{}, error) {
	path := "/_matrix/client/r0/rooms/" + roomID + "/state/" + eventType
	if stateKey != "" {
		path += "/" + stateKey
	}

	var result interface{}
	err := r.client.GET(ctx, path, nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetRoomMembers gets the members of a room
func (r *RoomAPI) GetRoomMembers(ctx context.Context, roomID string, at string) (*RoomMembersResponse, error) {
	query := map[string]string{}
	if at != "" {
		query["at"] = at
	}

	result := &RoomMembersResponse{}
	err := r.client.GET(ctx, "/_matrix/client/r0/rooms/"+roomID+"/members", query, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RoomMembersResponse represents the response from getting room members
type RoomMembersResponse struct {
	// Chunk contains the member events
	Chunk []MemberEvent `json:"chunk"`
}

// MemberEvent represents a member event
type MemberEvent struct {
	// EventID is the unique identifier of the event
	EventID string `json:"event_id"`

	// RoomID is the room where the event occurred
	RoomID string `json:"room_id"`

	// Type is the type of the event
	Type string `json:"type"`

	// Sender is the sender of the event
	Sender string `json:"sender"`

	// StateKey is the key of the state
	StateKey string `json:"state_key"`

	// Content is the content of the event
	Content MemberContent `json:"content"`

	// Timestamp is the timestamp of the event
	Timestamp int64 `json:"timestamp,omitempty"`
}

// MemberContent represents the content of a member event
type MemberContent struct {
	// Membership is the membership state ("join", "leave", "ban", "invite")
	Membership string `json:"membership"`

	// DisplayName is the display name of the user
	DisplayName string `json:"displayname,omitempty"`

	// AvatarURL is the avatar URL of the user
	AvatarURL string `json:"avatar_url,omitempty"`

	// Reason is the reason for the membership change
	Reason string `json:"reason,omitempty"`
}

// GetRoom gets the information of a room
func (r *RoomAPI) GetRoom(ctx context.Context, roomID string) (*Room, error) {
	result := &Room{}
	err := r.client.GET(ctx, "/_matrix/client/r0/rooms/"+roomID, nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateRoom creates a new room
func (r *RoomAPI) CreateRoom(ctx context.Context, req *CreateRoomRequest) (*CreateRoomResponse, error) {
	if req == nil {
		req = &CreateRoomRequest{}
	}

	// Set default visibility
	if req.Visibility == "" {
		req.Visibility = "private"
	}

	// Set default preset
	if req.Preset == "" {
		req.Preset = "private_chat"
	}

	result := &CreateRoomResponse{}
	err := r.client.POST(ctx, "/_matrix/client/r0/createRoom", req, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreatePublicRoom creates a new public room
func (r *RoomAPI) CreatePublicRoom(ctx context.Context, name, topic, alias string) (*CreateRoomResponse, error) {
	return r.CreateRoom(ctx, &CreateRoomRequest{
		Name:            name,
		Topic:           topic,
		RoomAliasName:   alias,
		Visibility:      "public",
		Preset:          "public_chat",
	})
}

// CreatePrivateRoom creates a new private room
func (r *RoomAPI) CreatePrivateRoom(ctx context.Context, name string, invite []string) (*CreateRoomResponse, error) {
	return r.CreateRoom(ctx, &CreateRoomRequest{
		Name:     name,
		Invite:   invite,
		Visibility: "private",
		Preset:   "private_chat",
	})
}

// SetRoomName sets the name of a room
func (r *RoomAPI) SetRoomName(ctx context.Context, roomID, name string) error {
	body := map[string]string{
		"name": name,
	}
	return r.client.PUT(ctx, "/_matrix/client/r0/rooms/"+roomID+"/state/m.room.name", body, nil)
}

// SetRoomTopic sets the topic of a room
func (r *RoomAPI) SetRoomTopic(ctx context.Context, roomID, topic string) error {
	body := map[string]string{
		"topic": topic,
	}
	return r.client.PUT(ctx, "/_matrix/client/r0/rooms/"+roomID+"/state/m.room.topic", body, nil)
}

// SetRoomAvatar sets the avatar of a room
func (r *RoomAPI) SetRoomAvatar(ctx context.Context, roomID, avatarURL string) error {
	body := map[string]string{
		"url": avatarURL,
	}
	return r.client.PUT(ctx, "/_matrix/client/r0/rooms/"+roomID+"/state/m.room.avatar", body, nil)
}

// GetJoinedRooms gets the rooms that the user has joined
func (r *RoomAPI) GetJoinedRooms(ctx context.Context) (*JoinedRoomsResponse, error) {
	result := &JoinedRoomsResponse{}
	err := r.client.GET(ctx, "/_matrix/client/r0/joined_rooms", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// JoinedRoomsResponse represents the response from getting joined rooms
type JoinedRoomsResponse struct {
	// JoinedRooms is a list of room IDs
	JoinedRooms []string `json:"joined_rooms"`
}

// GetRoomPowerLevels gets the power levels of a room
func (r *RoomAPI) GetRoomPowerLevels(ctx context.Context, roomID string) (*PowerLevels, error) {
	result := &PowerLevels{}
	err := r.client.GET(ctx, "/_matrix/client/r0/rooms/"+roomID+"/state/m.room.power_levels", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SetRoomPowerLevels sets the power levels of a room
func (r *RoomAPI) SetRoomPowerLevels(ctx context.Context, roomID string, levels *PowerLevels) error {
	return r.client.PUT(ctx, "/_matrix/client/r0/rooms/"+roomID+"/state/m.room.power_levels", levels, nil)
}

// GetRoomAliases gets the aliases of a room
func (r *RoomAPI) GetRoomAliases(ctx context.Context, roomID string) (*RoomAliasesResponse, error) {
	result := &RoomAliasesResponse{}
	err := r.client.GET(ctx, "/_matrix/client/r0/rooms/"+roomID+"/aliases", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RoomAliasesResponse represents the response from getting room aliases
type RoomAliasesResponse struct {
	// Aliases is a list of room aliases
	Aliases []string `json:"aliases"`
}

// GetUserRooms gets the rooms that a user has joined (admin API)
func (r *RoomAPI) GetUserRooms(ctx context.Context, userID string) (*JoinedRoomsResponse, error) {
	result := &JoinedRoomsResponse{}
	err := r.client.GET(ctx, "/_matrix/client/r0/admin/rooms/"+userID+"/joined_rooms", nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetRoomDetails gets the details of a room (admin API)
func (r *RoomAPI) GetRoomDetails(ctx context.Context, roomID string) (*RoomDetailsResponse, error) {
	result := &RoomDetailsResponse{}
	err := r.client.GET(ctx, "/_matrix/client/r0/admin/rooms/"+roomID, nil, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RoomDetailsResponse represents the response from getting room details
type RoomDetailsResponse struct {
	// RoomID is the room ID
	RoomID string `json:"room_id"`

	// Name is the room name
	Name string `json:"name"`

	// Avatar is the room avatar
	Avatar string `json:"avatar"`

	// Topic is the room topic
	Topic string `json:"topic"`

	// Creator is the creator of the room
	Creator string `json:"creator"`

	// GuestCanJoin indicates if guests can join
	GuestCanJoin bool `json:"guest_can_join"`

	// JoinRule is the join rule
	JoinRule string `json:"join_rule"`

	// Private indicates if the room is private
	Private bool `json:"private"`

	// JoinedMembers is the number of joined members
	JoinedMembers int `json:"joined_members"`

	// JoinedLocalMembers is the number of joined local members
	JoinedLocalMembers int `json:"joined_local_members"`

	// CreateAt is the creation time
	CreateAt time.Time `json:"create_at"`
}

// DeleteRoom deletes a room (admin API)
func (r *RoomAPI) DeleteRoom(ctx context.Context, roomID string, purge bool) error {
	body := map[string]interface{}{
		"purge": purge,
	}
	return r.client.DELETE(ctx, "/_matrix/client/r0/admin/rooms/"+roomID, nil, nil)
}

// ForgetRoom forgets a room
func (r *RoomAPI) ForgetRoom(ctx context.Context, roomID string) error {
	return r.client.POST(ctx, "/_matrix/client/r0/rooms/"+roomID+"/forget", nil, nil)
}
