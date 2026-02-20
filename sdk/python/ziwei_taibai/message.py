"""消息 API 模块

提供消息发送和订阅功能。
"""

from typing import Any, Callable, Dict, Optional
from dataclasses import dataclass

from .http_client import HTTPClient, HTTPResponse


@dataclass
class Message:
    """消息对象"""
    id: str
    sender: str
    recipient: str
    content: Any
    type: str = "text"
    timestamp: Optional[float] = None
    metadata: Optional[Dict[str, Any]] = None


class MessageAPI:
    """消息 API"""
    
    def __init__(self, http_client: HTTPClient):
        """
        初始化消息 API
        
        Args:
            http_client: HTTP 客户端实例
        """
        self._http = http_client
    
    async def send_message(
        self,
        recipient: str,
        content: Any,
        msg_type: str = "text",
        metadata: Optional[Dict[str, Any]] = None
    ) -> Message:
        """
        发送消息
        
        Args:
            recipient: 接收者 ID
            content: 消息内容
            msg_type: 消息类型 (text, image, file, etc.)
            metadata: 附加元数据
            
        Returns:
            Message 对象
        """
        data = {
            "recipient": recipient,
            "content": content,
            "type": msg_type,
        }
        
        if metadata:
            data["metadata"] = metadata
        
        response = await self._http.post("/api/messages", data=data)
        
        if not response.ok:
            raise Exception(f"Failed to send message: {response.data}")
        
        msg_data = response.data
        return Message(
            id=msg_data.get("id", ""),
            sender=msg_data.get("sender", ""),
            recipient=msg_data.get("recipient", ""),
            content=msg_data.get("content", ""),
            type=msg_data.get("type", msg_type),
            timestamp=msg_data.get("timestamp"),
            metadata=msg_data.get("metadata")
        )
    
    async def get_message(self, message_id: str) -> Message:
        """
        获取消息详情
        
        Args:
            message_id: 消息 ID
            
        Returns:
            Message 对象
        """
        response = await self._http.get(f"/api/messages/{message_id}")
        
        if not response.ok:
            raise Exception(f"Failed to get message: {response.data}")
        
        msg_data = response.data
        return Message(
            id=msg_data.get("id", ""),
            sender=msg_data.get("sender", ""),
            recipient=msg_data.get("recipient", ""),
            content=msg_data.get("content", ""),
            type=msg_data.get("type", "text"),
            timestamp=msg_data.get("timestamp"),
            metadata=msg_data.get("metadata")
        )
    
    async def list_messages(
        self,
        user_id: Optional[str] = None,
        limit: int = 50,
        offset: int = 0
    ) -> list[Message]:
        """
        获取消息列表
        
        Args:
            user_id: 用户 ID（可选，用于筛选）
            limit: 返回数量限制
            offset: 偏移量
            
        Returns:
            Message 对象列表
        """
        params = {
            "limit": limit,
            "offset": offset,
        }
        
        if user_id:
            params["user_id"] = user_id
        
        response = await self._http.get("/api/messages", params=params)
        
        if not response.ok:
            raise Exception(f"Failed to list messages: {response.data}")
        
        messages = []
        for msg_data in response.data.get("messages", []):
            messages.append(Message(
                id=msg_data.get("id", ""),
                sender=msg_data.get("sender", ""),
                recipient=msg_data.get("recipient", ""),
                content=msg_data.get("content", ""),
                type=msg_data.get("type", "text"),
                timestamp=msg_data.get("timestamp"),
                metadata=msg_data.get("metadata")
            ))
        
        return messages
    
    async def delete_message(self, message_id: str) -> bool:
        """
        删除消息
        
        Args:
            message_id: 消息 ID
            
        Returns:
            是否删除成功
        """
        response = await self._http.delete(f"/api/messages/{message_id}")
        
        if not response.ok:
            raise Exception(f"Failed to delete message: {response.data}")
        
        return True


class SubscriptionManager:
    """订阅管理器"""
    
    def __init__(self, ws_client):
        """
        初始化订阅管理器
        
        Args:
            ws_client: WebSocket 客户端实例
        """
        self._ws = ws_client
        self._subscriptions: Dict[str, Callable[[Any], Any]] = {}
    
    async def subscribe(
        self,
        event_type: str,
        callback: Callable[[Any], Any]
    ):
        """
        订阅消息
        
        Args:
            event_type: 事件类型
            callback: 回调函数
        """
        self._subscriptions[event_type] = callback
        
        # 发送订阅请求
        await self._ws.send({
            "type": "subscribe",
            "event_type": event_type
        })
    
    async def unsubscribe(self, event_type: str):
        """
        取消订阅
        
        Args:
            event_type: 事件类型
        """
        if event_type in self._subscriptions:
            del self._subscriptions[event_type]
        
        # 发送取消订阅请求
        await self._ws.send({
            "type": "unsubscribe",
            "event_type": event_type
        })
    
    def handle_message(self, msg):
        """处理接收到的消息"""
        event_type = msg.data.get("type", "")
        
        if event_type in self._subscriptions:
            callback = self._subscriptions[event_type]
            try:
                callback(msg.data)
            except Exception as e:
                print(f"Subscription callback error: {e}")
    
    @property
    def subscribed_events(self) -> list[str]:
        """获取已订阅的事件类型列表"""
        return list(self._subscriptions.keys())
