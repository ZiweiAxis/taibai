package taibai

import (
	"encoding/json"
	"fmt"
	"time"
)

// ============ 消息类型定义 ============

// MessageType 消息类型常量
const (
	MessageTypeText     = "text"      // 文本消息
	MessageTypeImage    = "image"     // 图片消息
	MessageTypeVoice    = "voice"     // 语音消息
	MessageTypeVideo    = "video"     // 视频消息
	MessageTypeFile     = "file"      // 文件消息
	MessageTypeCard     = "card"      // 卡片消息
	MessageTypeCallback = "callback"  // 卡片回调
	MessageTypeSystem   = "system"    // 系统消息
)

// EventType 事件类型常量
const (
	EventUserMessage     = "user_message"      // 用户消息
	EventCardCallback    = "card_callback"     // 卡片回调
	EventApprovalChange  = "approval_change"   // 审批状态变更
	EventSubscribe       = "subscribe"         // 订阅事件
	EventUnsubscribe     = "unsubscribe"       // 取消订阅
	EventPing            = "ping"              // 心跳
	EventPong            = "pong"              // 心跳响应
)

// ============ 消息结构体 ============

// UserMessage 用户消息
type UserMessage struct {
	MessageID   string          `json:"message_id"`   // 消息 ID
	UserID       string          `json:"user_id"`      // 用户 ID
	UserName     string          `json:"user_name"`    // 用户名称
	Content      string          `json:"content"`      // 消息内容
	MessageType  string          `json:"message_type"` // 消息类型
	ChannelID    string          `json:"channel_id"`   // 频道 ID
	GroupID      string          `json:"group_id"`     // 群组 ID (可选)
	Timestamp    int64           `json:"timestamp"`    // 时间戳
	Raw          json.RawMessage `json:"raw"`          // 原始消息
}

// CardCallback 卡片回调
type CardCallback struct {
	CallbackID   string          `json:"callback_id"`   // 回调 ID
	UserID        string          `json:"user_id"`       // 用户 ID
	CardID        string          `json:"card_id"`       // 卡片 ID
	Action        string          `json:"action"`        // 点击的按钮 ID
	Data          json.RawMessage `json:"data"`          // 自定义数据
	Timestamp     int64           `json:"timestamp"`     // 时间戳
	Raw           json.RawMessage `json:"raw"`           // 原始消息
}

// ApprovalChange 审批状态变更
type ApprovalChange struct {
	ApprovalID    string `json:"approval_id"`    // 审批单 ID
	Status        string `json:"status"`         // 审批状态: pending/approved/rejected
	ApplicantID   string `json:"applicant_id"`   // 申请人 ID
	ApplicantName string `json:"applicant_name"` // 申请人名称
	ApproverID    string `json:"approver_id"`   // 审批人 ID
	ApproverName  string `json:"approver_name"`  // 审批人名称
	Comment       string `json:"comment"`        // 审批意见
	Timestamp     int64  `json:"timestamp"`      // 时间戳
	Raw           json.RawMessage `json:"raw"`   // 原始消息
}

// ============ 消息处理器 ============

// MessageHandler 消息处理器
type MessageHandler struct {
	// 用户消息处理
	UserMessageHandlers []func(msg *UserMessage)

	// 卡片回调处理
	CardCallbackHandlers []func(callback *CardCallback)

	// 审批状态变更处理
	ApprovalChangeHandlers []func(change *ApprovalChange)

	// 系统消息处理
	SystemHandlers []func(event string, data json.RawMessage)
}

// NewMessageHandler 创建消息处理器
func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		UserMessageHandlers:    make([]func(msg *UserMessage), 0),
		CardCallbackHandlers:   make([]func(callback *CardCallback), 0),
		ApprovalChangeHandlers: make([]func(change *ApprovalChange), 0),
		SystemHandlers:         make([]func(event string, data json.RawMessage), 0),
	}
}

// OnUserMessage 注册用户消息处理函数
func (h *MessageHandler) OnUserMessage(fn func(msg *UserMessage)) {
	h.UserMessageHandlers = append(h.UserMessageHandlers, fn)
}

// OnCardCallback 注册卡片回调处理函数
func (h *MessageHandler) OnCardCallback(fn func(callback *CardCallback)) {
	h.CardCallbackHandlers = append(h.CardCallbackHandlers, fn)
}

// OnApprovalChange 注册审批状态变更处理函数
func (h *MessageHandler) OnApprovalChange(fn func(change *ApprovalChange)) {
	h.ApprovalChangeHandlers = append(h.ApprovalChangeHandlers, fn)
}

// OnSystem 注册系统消息处理函数
func (h *MessageHandler) OnSystem(fn func(event string, data json.RawMessage)) {
	h.SystemHandlers = append(h.SystemHandlers, fn)
}

// Handle 处理 WebSocket 消息
func (h *MessageHandler) Handle(wsMsg *WSMessage) error {
	switch wsMsg.Event {
	case EventUserMessage:
		return h.handleUserMessage(wsMsg.Payload)
	case EventCardCallback:
		return h.handleCardCallback(wsMsg.Payload)
	case EventApprovalChange:
		return h.handleApprovalChange(wsMsg.Payload)
	default:
		return h.handleSystem(wsMsg.Event, wsMsg.Payload)
	}
}

// handleUserMessage 处理用户消息
func (h *MessageHandler) handleUserMessage(payload json.RawMessage) error {
	var msg UserMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("解析用户消息失败: %w", err)
	}

	// 设置时间戳
	if msg.Timestamp == 0 {
		msg.Timestamp = time.Now().Unix()
	}

	// 调用所有处理函数
	for _, fn := range h.UserMessageHandlers {
		fn(&msg)
	}

	return nil
}

// handleCardCallback 处理卡片回调
func (h *MessageHandler) handleCardCallback(payload json.RawMessage) error {
	var callback CardCallback
	if err := json.Unmarshal(payload, &callback); err != nil {
		return fmt.Errorf("解析卡片回调失败: %w", err)
	}

	// 设置时间戳
	if callback.Timestamp == 0 {
		callback.Timestamp = time.Now().Unix()
	}

	// 调用所有处理函数
	for _, fn := range h.CardCallbackHandlers {
		fn(&callback)
	}

	return nil
}

// handleApprovalChange 处理审批状态变更
func (h *MessageHandler) handleApprovalChange(payload json.RawMessage) error {
	var change ApprovalChange
	if err := json.Unmarshal(payload, &change); err != nil {
		return fmt.Errorf("解析审批状态变更失败: %w", err)
	}

	// 设置时间戳
	if change.Timestamp == 0 {
		change.Timestamp = time.Now().Unix()
	}

	// 调用所有处理函数
	for _, fn := range h.ApprovalChangeHandlers {
		fn(&change)
	}

	return nil
}

// handleSystem 处理系统消息
func (h *MessageHandler) handleSystem(event string, data json.RawMessage) error {
	for _, fn := range h.SystemHandlers {
		fn(event, data)
	}
	return nil
}

// ============ 消息工具函数 ============

// ParseUserMessage 解析用户消息
func ParseUserMessage(data json.RawMessage) (*UserMessage, error) {
	var msg UserMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// ParseCardCallback 解析卡片回调
func ParseCardCallback(data json.RawMessage) (*CardCallback, error) {
	var callback CardCallback
	if err := json.Unmarshal(data, &callback); err != nil {
		return nil, err
	}
	return &callback, nil
}

// ParseApprovalChange 解析审批状态变更
func ParseApprovalChange(data json.RawMessage) (*ApprovalChange, error) {
	var change ApprovalChange
	if err := json.Unmarshal(data, &change); err != nil {
		return nil, err
	}
	return &change, nil
}

// ============ 便捷的回调包装器 ============

// UserMessageHandlerFunc 用户消息处理函数类型
type UserMessageHandlerFunc func(msg *UserMessage)

// CardCallbackHandlerFunc 卡片回调处理函数类型
type CardCallbackHandlerFunc func(callback *CardCallback)

// ApprovalChangeHandlerFunc 审批状态变更处理函数类型
type ApprovalChangeHandlerFunc func(change *ApprovalChange)

// ============ Client 包装器 (简化使用) ============

// WSClient WebSocket 客户端包装器
type WSClient struct {
	*WebSocketClient
	*MessageHandler
}

// NewWSClient 创建 WebSocket 客户端包装器
func NewWSClient(config *WebSocketConfig) *WSClient {
	ws := NewWebSocketClient(config)
	handler := NewMessageHandler()

	// 自动处理消息
	ws.OnMessage = func(msg *WSMessage) {
		if err := handler.Handle(msg); err != nil {
			ws.OnError(err)
		}
	}

	return &WSClient{
		WebSocketClient: ws,
		MessageHandler:  handler,
	}
}
