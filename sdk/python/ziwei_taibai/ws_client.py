"""WebSocket 客户端模块

提供异步 WebSocket 连接支持，包括心跳保活、自动重连和消息接收。
"""

import asyncio
import json
import logging
from typing import Any, Callable, Dict, Optional
from dataclasses import dataclass
from enum import Enum

import aiohttp


logger = logging.getLogger(__name__)


class ConnectionState(Enum):
    """连接状态"""
    DISCONNECTED = "disconnected"
    CONNECTING = "connecting"
    CONNECTED = "connected"
    RECONNECTING = "reconnecting"


@dataclass
class WSMessage:
    """WebSocket 消息"""
    type: str
    data: Any
    raw: str


class WSClient:
    """异步 WebSocket 客户端"""
    
    def __init__(
        self,
        url: str,
        heartbeat_interval: float = 30.0,
        reconnect_delay: float = 1.0,
        max_reconnect_delay: float = 60.0,
        reconnect_attempts: int = 0,  # 0 = 无限重连
        headers: Optional[Dict[str, str]] = None
    ):
        """
        初始化 WebSocket 客户端
        
        Args:
            url: WebSocket 服务器 URL
            heartbeat_interval: 心跳间隔（秒）
            reconnect_delay: 初始重连延迟（秒）
            max_reconnect_delay: 最大重连延迟（秒）
            reconnect_attempts: 最大重连次数（0 = 无限）
            headers: 连接 headers
        """
        self.url = url
        self.heartbeat_interval = heartbeat_interval
        self.reconnect_delay = reconnect_delay
        self.max_reconnect_delay = max_reconnect_delay
        self.reconnect_attempts = reconnect_attempts
        self.headers = headers or {}
        
        self._state = ConnectionState.DISCONNECTED
        self._ws: Optional[aiohttp.ClientWebSocketResponse] = None
        self._session: Optional[aiohttp.ClientSession] = None
        self._heartbeat_task: Optional[asyncio.Task] = None
        self._receive_task: Optional[asyncio.Task] = None
        self._reconnect_count = 0
        self._running = False
        
        # 消息回调
        self._message_callbacks: list[Callable[[WSMessage], Any]] = []
        self._connect_callbacks: list[Callable[[], Any]] = []
        self._disconnect_callbacks: list[Callable[[], Any]] = []
        self._error_callbacks: list[Callable[[Exception], Any]] = []
    
    @property
    def state(self) -> ConnectionState:
        """获取连接状态"""
        return self._state
    
    @property
    def is_connected(self) -> bool:
        """是否已连接"""
        return self._state == ConnectionState.CONNECTED
    
    def on_message(self, callback: Callable[[WSMessage], Any]):
        """注册消息回调"""
        self._message_callbacks.append(callback)
    
    def on_connect(self, callback: Callable[[], Any]):
        """注册连接回调"""
        self._connect_callbacks.append(callback)
    
    def on_disconnect(self, callback: Callable[[], Any]):
        """注册断开连接回调"""
        self._disconnect_callbacks.append(callback)
    
    def on_error(self, callback: Callable[[Exception], Any]):
        """注册错误回调"""
        self._error_callbacks.append(callback)
    
    async def connect(self):
        """建立 WebSocket 连接"""
        if self._state == ConnectionState.CONNECTED:
            return
        
        self._running = True
        self._state = ConnectionState.CONNECTING
        
        try:
            self._session = aiohttp.ClientSession()
            self._ws = await self._session.ws_connect(
                self.url,
                headers=self.headers,
                autoclose=False
            )
            self._state = ConnectionState.CONNECTED
            self._reconnect_count = 0
            
            logger.info(f"WebSocket connected to {self.url}")
            
            # 启动心跳和接收任务
            self._heartbeat_task = asyncio.create_task(self._heartbeat_loop())
            self._receive_task = asyncio.create_task(self._receive_loop())
            
            # 触发连接回调
            for callback in self._connect_callbacks:
                try:
                    callback()
                except Exception as e:
                    logger.error(f"Connect callback error: {e}")
            
        except Exception as e:
            self._state = ConnectionState.DISCONNECTED
            logger.error(f"WebSocket connection failed: {e}")
            await self._handle_error(e)
            await self._maybe_reconnect()
    
    async def _receive_loop(self):
        """接收消息循环"""
        while self._running and self._ws:
            try:
                msg = await self._ws.receive()
                
                if msg.type == aiohttp.WSMsgType.TEXT:
                    try:
                        data = json.loads(msg.data)
                    except json.JSONDecodeError:
                        data = msg.data
                    
                    ws_msg = WSMessage(
                        type='text',
                        data=data,
                        raw=msg.data
                    )
                    
                    # 触发消息回调
                    for callback in self._message_callbacks:
                        try:
                            callback(ws_msg)
                        except Exception as e:
                            logger.error(f"Message callback error: {e}")
                            
                elif msg.type == aiohttp.WSMsgType.ERROR:
                    logger.error(f"WebSocket error: {msg.data}")
                    break
                elif msg.type in (aiohttp.WSMsgType.CLOSE, aiohttp.WSMsgType.CLOSED):
                    logger.info("WebSocket closed by server")
                    break
                elif msg.type == aiohttp.WSMsgType.PING:
                    await self._ws.pong()
                elif msg.type == aiohttp.WSMsgType.PONG:
                    pass  # 心跳响应
                    
            except asyncio.CancelledError:
                break
            except Exception as e:
                logger.error(f"Receive error: {e}")
                await self._handle_error(e)
                break
        
        await self._handle_disconnect()
    
    async def _heartbeat_loop(self):
        """心跳循环"""
        while self._running and self._ws:
            try:
                await asyncio.sleep(self.heartbeat_interval)
                
                if self._ws and not self._ws.closed:
                    await self._ws.send_str(json.dumps({
                        'type': 'heartbeat',
                        'timestamp': asyncio.get_event_loop().time()
                    }))
                    logger.debug("Heartbeat sent")
                    
            except asyncio.CancelledError:
                break
            except Exception as e:
                logger.error(f"Heartbeat error: {e}")
                break
    
    async def send(self, data: Dict[str, Any]):
        """
        发送消息
        
        Args:
            data: 要发送的数据（会自动序列化为 JSON）
        """
        if not self._ws or self._ws.closed:
            raise ConnectionError("WebSocket is not connected")
        
        await self._ws.send_str(json.dumps(data))
        logger.debug(f"Sent: {data}")
    
    async def send_text(self, text: str):
        """发送文本消息"""
        if not self._ws or self._ws.closed:
            raise ConnectionError("WebSocket is not connected")
        
        await self._ws.send_str(text)
    
    async def close(self, code: int = 1000, reason: str = ""):
        """
        关闭连接
        
        Args:
            code: 关闭代码
            reason: 关闭原因
        """
        self._running = False
        
        # 取消任务
        if self._heartbeat_task:
            self._heartbeat_task.cancel()
            try:
                await self._heartbeat_task
            except asyncio.CancelledError:
                pass
        
        if self._receive_task:
            self._receive_task.cancel()
            try:
                await self._receive_task
            except asyncio.CancelledError:
                pass
        
        # 关闭连接
        if self._ws and not self._ws.closed:
            await self._ws.close(code=code, message=reason.encode())
        
        if self._session:
            await self._session.close()
        
        self._state = ConnectionState.DISCONNECTED
        logger.info("WebSocket closed")
    
    async def _handle_disconnect(self):
        """处理断开连接"""
        if self._state == ConnectionState.DISCONNECTED:
            return
        
        self._state = ConnectionState.DISCONNECTED
        
        # 触发断开回调
        for callback in self._disconnect_callbacks:
            try:
                callback()
            except Exception as e:
                logger.error(f"Disconnect callback error: {e}")
        
        await self._maybe_reconnect()
    
    async def _maybe_reconnect(self):
        """尝试重连"""
        if not self._running:
            return
        
        # 检查重连次数
        if self.reconnect_attempts > 0 and self._reconnect_count >= self.reconnect_attempts:
            logger.error("Max reconnect attempts reached")
            return
        
        self._state = ConnectionState.RECONNECTING
        self._reconnect_count += 1
        
        delay = min(
            self.reconnect_delay * (2 ** (self._reconnect_count - 1)),
            self.max_reconnect_delay
        )
        
        logger.info(f"Reconnecting in {delay}s (attempt {self._reconnect_count})")
        
        await asyncio.sleep(delay)
        
        if self._running:
            await self.connect()
    
    async def _handle_error(self, error: Exception):
        """处理错误"""
        for callback in self._error_callbacks:
            try:
                callback(error)
            except Exception as e:
                logger.error(f"Error callback error: {e}")
    
    async def __aenter__(self):
        await self.connect()
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        await self.close()
