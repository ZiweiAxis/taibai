package taibai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketConfig WebSocket 配置
type WebSocketConfig struct {
	URL            string        // WebSocket 服务器地址
	Token          string        // 认证 Token
	HeartbeatInterval time.Duration // 心跳间隔 (默认 30 秒)
	ReconnectDelay   time.Duration // 重连延迟 (默认 5 秒)
	MaxReconnectAttempts int       // 最大重连次数 (默认 0 表示无限)
}

// WebSocketClient WebSocket 客户端
type WebSocketClient struct {
	config       *WebSocketConfig
	conn         *websocket.Conn
	isConnected bool
	isReconnecting bool
_mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc

	// 回调函数
	OnConnect       func()       // 连接成功回调
	OnDisconnect    func(error)  // 断线回调
	OnMessage       func(msg *WSMessage) // 消息接收回调
	OnError         func(error)  // 错误回调

	// 订阅管理
	subscriptions map[string]bool
	subMu         sync.RWMutex

	// 内部消息
	readChan  chan *WSMessage
	writeChan chan []byte
	closeChan chan struct{}
}

// WSMessage WebSocket 消息结构
type WSMessage struct {
	Type    string          `json:"type"`    // 消息类型
	Event   string          `json:"event"`   // 事件类型
	Payload json.RawMessage `json:"payload"` // 消息内容
	Seq     int64           `json:"seq"`     // 序列号
}

// NewWebSocketClient 创建 WebSocket 客户端
func NewWebSocketClient(config *WebSocketConfig) *WebSocketClient {
	if config.HeartbeatInterval == 0 {
		config.HeartbeatInterval = 30 * time.Second
	}
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = 5 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketClient{
		config:        config,
		isConnected:   false,
		ctx:           ctx,
		cancel:        cancel,
		subscriptions: make(map[string]bool),
		readChan:      make(chan *WSMessage, 100),
		writeChan:     make(chan []byte, 100),
		closeChan:     make(chan struct{}),
	}
}

// Connect 连接到 WebSocket 服务器
func (c *WebSocketClient) Connect() error {
	c._mu.Lock()
	if c.isConnected {
		c._mu.Unlock()
		return nil
	}
	c._mu.Unlock()

	// 构建认证 URL
	url := fmt.Sprintf("%s?token=%s", c.config.Token, c.config.Token)

	// 设置 WebSocket 握手超时
	dialer := &websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	// 添加认证头
	header := http.Header{}
	header.Set("Authorization", "Bearer "+c.config.Token)

	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		if c.OnError != nil {
			c.OnError(fmt.Errorf("连接失败: %w", err))
		}
		return err
	}

	c._mu.Lock()
	c.conn = conn
	c.isConnected = true
	c.isReconnecting = false
	c._mu.Unlock()

	// 启动读写协程
	go c.readLoop()
	go c.writeLoop()
	go c.heartbeatLoop()

	// 触发连接成功回调
	if c.OnConnect != nil {
		c.OnConnect()
	}

	// 重新订阅
	c.resubscribe()

	return nil
}

// Disconnect 断开连接
func (c *WebSocketClient) Disconnect() {
	c._mu.Lock()
	defer c._mu.Unlock()

	if c.cancel != nil {
		c.cancel()
	}

	if c.conn != nil {
		c.conn.Close()
		c.isConnected = false
	}

	close(c.closeChan)
}

// Reconnect 重新连接
func (c *WebSocketClient) Reconnect() {
	if c.isReconnecting {
		return
	}

	c.isReconnecting = true
	defer func() {
		c.isReconnecting = false
	}()

	attempts := 0
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		attempts++
		if c.config.MaxReconnectAttempts > 0 && attempts > c.config.MaxReconnectAttempts {
			if c.OnError != nil {
				c.OnError(fmt.Errorf("达到最大重连次数: %d", c.config.MaxReconnectAttempts))
			}
			return
		}

		time.Sleep(c.config.ReconnectDelay)

		if err := c.Connect(); err == nil {
			return
		}
	}
}

// readLoop 读取消息循环
func (c *WebSocketClient) readLoop() {
	defer func() {
		c.handleDisconnect()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if c.OnError != nil {
					c.OnError(fmt.Errorf("读取消息失败: %w", err))
				}
			}
			return
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			if c.OnError != nil {
				c.OnError(fmt.Errorf("解析消息失败: %w", err))
			}
			continue
		}

		// 处理心跳响应
		if wsMsg.Type == "pong" || wsMsg.Event == "pong" {
			continue
		}

		// 发送到消息通道
		select {
		case c.readChan <- &wsMsg:
		default:
		}

		// 触发消息回调
		if c.OnMessage != nil {
			c.OnMessage(&wsMsg)
		}
	}
}

// writeLoop 写入消息循环
func (c *WebSocketClient) writeLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case message := <-c.writeChan:
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				if c.OnError != nil {
					c.OnError(fmt.Errorf("发送消息失败: %w", err))
				}
			}
		case <-ticker.C:
			// 保持连接活跃
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				if c.OnError != nil {
					c.OnError(fmt.Errorf("发送 ping 失败: %w", err))
				}
			}
		}
	}
}

// heartbeatLoop 心跳循环
func (c *WebSocketClient) heartbeatLoop() {
	ticker := time.NewTicker(c.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.sendPing()
		}
	}
}

// sendPing 发送心跳
func (c *WebSocketClient) sendPing() {
	pingMsg := WSMessage{
		Type:  "ping",
		Event: "ping",
		Seq:   time.Now().UnixNano(),
	}
	data, err := json.Marshal(pingMsg)
	if err != nil {
		return
	}

	select {
	case c.writeChan <- data:
	default:
	}
}

// handleDisconnect 处理断线
func (c *WebSocketClient) handleDisconnect() {
	c._mu.Lock()
	wasConnected := c.isConnected
	c.isConnected = false
	c._mu.Unlock()

	if wasConnected && c.OnDisconnect != nil {
		c.OnDisconnect(fmt.Errorf("连接已断开"))
	}

	// 自动重连
	go c.Reconnect()
}

// Subscribe 订阅消息
func (c *WebSocketClient) Subscribe(event string) error {
	c.subMu.Lock()
	defer c.subMu.Unlock()

	if c.subscriptions[event] {
		return nil // 已经订阅
	}

	subscribeMsg := WSSubscribeRequest{
		Type:  "subscribe",
		Event: event,
	}

	data, err := json.Marshal(subscribeMsg)
	if err != nil {
		return err
	}

	select {
	case c.writeChan <- data:
		c.subscriptions[event] = true
		return nil
	default:
		return fmt.Errorf("发送通道已满")
	}
}

// Unsubscribe 取消订阅
func (c *WebSocketClient) Unsubscribe(event string) error {
	c.subMu.Lock()
	defer c.subMu.Unlock()

	if !c.subscriptions[event] {
		return nil // 未订阅
	}

	unsubscribeMsg := WSSubscribeRequest{
		Type:  "unsubscribe",
		Event: event,
	}

	data, err := json.Marshal(unsubscribeMsg)
	if err != nil {
		return err
	}

	select {
	case c.writeChan <- data:
		delete(c.subscriptions, event)
		return nil
	default:
		return fmt.Errorf("发送通道已满")
	}
}

// resubscribe 重新订阅
func (c *WebSocketClient) resubscribe() {
	c.subMu.RLock()
	events := make([]string, 0, len(c.subscriptions))
	for event := range c.subscriptions {
		events = append(events, event)
	}
	c.subMu.RUnlock()

	for _, event := range events {
		subscribeMsg := WSSubscribeRequest{
			Type:  "subscribe",
			Event: event,
		}
		data, _ := json.Marshal(subscribeMsg)
		select {
		case c.writeChan <- data:
		default:
		}
	}
}

// IsConnected 检查是否已连接
func (c *WebSocketClient) IsConnected() bool {
	c._mu.RLock()
	defer c._mu.RUnlock()
	return c.isConnected
}

// WSSubscribeRequest 订阅请求
type WSSubscribeRequest struct {
	Type  string `json:"type"`
	Event string `json:"event"`
}
