package taibai

import (
	"context"
)

// ==================== User API ====================

type UserAPI struct {
	client *Client
}

type GetUserRequest struct {
	UserID string `json:"user_id,omitempty"`
	DID    string `json:"did,omitempty"`
}

type GetUserResponse struct {
	UserID      string `json:"user_id"`
	DID         string `json:"did"`
	DisplayName string `json:"display_name,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
}

type GetUserRolesRequest struct {
	UserID string `json:"user_id"`
	RoomID string `json:"room_id,omitempty"`
}

type GetUserRolesResponse struct {
	Roles []string `json:"roles"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id,omitempty"`
	DID    string `json:"did,omitempty"`
}

func (u *UserAPI) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	resp := &GetUserResponse{}
	err := u.client.GET(ctx, "/api/v1/users/get", map[string]string{
		"user_id": req.UserID,
		"did":    req.DID,
	}, resp)
	return resp, err
}

func (u *UserAPI) GetUserRoles(ctx context.Context, req *GetUserRolesRequest) (*GetUserRolesResponse, error) {
	resp := &GetUserRolesResponse{}
	err := u.client.GET(ctx, "/api/v1/users/roles", map[string]string{
		"user_id": req.UserID,
		"room_id": req.RoomID,
	}, resp)
	return resp, err
}

func (u *UserAPI) ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	resp := &ValidateTokenResponse{}
	err := u.client.POST(ctx, "/api/v1/users/validate-token", req, resp)
	return resp, err
}

// ==================== Approval API ====================

type ApprovalAPI struct {
	client *Client
}

type SendApprovalRequestRequest struct {
	RequestID       string `json:"request_id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	TraceID         string `json:"trace_id"`
	GatewayBaseURL  string `json:"gateway_base_url"`
	AgentDID        string `json:"agent_did"`
	AgentName       string `json:"agent_name"`
	Operation       string `json:"operation"`
	Target          string `json:"target"`
	RiskLevel       string `json:"risk_level"`
	RequesterDID    string `json:"requester_did"`
	RequesterName   string `json:"requester_name"`
}

type SendApprovalRequestResponse struct {
	ApprovalID string `json:"approval_id"`
	Status     string `json:"status"`
}

type ApprovalCallbackRequest struct {
	ApprovalID  string `json:"approval_id"`
	Approved    bool   `json:"approved"`
	ApprovedBy  string `json:"approved_by"`
	Reason      string `json:"reason,omitempty"`
}

func (a *ApprovalAPI) SendApprovalRequest(ctx context.Context, req *SendApprovalRequestRequest) (*SendApprovalRequestResponse, error) {
	resp := &SendApprovalRequestResponse{}
	err := a.client.POST(ctx, "/api/v1/delivery/approval-request", req, resp)
	return resp, err
}

func (a *ApprovalAPI) HandleCallback(ctx context.Context, req *ApprovalCallbackRequest) error {
	return a.client.POST(ctx, "/api/v1/approval/callback", req, nil)
}
