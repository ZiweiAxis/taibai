"""HTTP 客户端模块

提供异步 HTTP 通信支持，支持 GET/POST/PUT/DELETE 方法，
JSON 序列化/反序列化，以及超时控制。
"""

import asyncio
import json
from typing import Any, Dict, Optional
from dataclasses import dataclass

import aiohttp


@dataclass
class HTTPResponse:
    """HTTP 响应封装"""
    status: int
    data: Any
    headers: Dict[str, str]
    
    @property
    def ok(self) -> bool:
        return 200 <= self.status < 300


class HTTPClient:
    """异步 HTTP 客户端"""
    
    def __init__(
        self,
        base_url: str,
        timeout: float = 30.0,
        headers: Optional[Dict[str, str]] = None
    ):
        """
        初始化 HTTP 客户端
        
        Args:
            base_url: 基础 URL
            timeout: 超时时间（秒）
            headers: 默认请求头
        """
        self.base_url = base_url.rstrip('/')
        self.timeout = aiohttp.ClientTimeout(total=timeout)
        self.default_headers = headers or {}
        self._session: Optional[aiohttp.ClientSession] = None
    
    async def _get_session(self) -> aiohttp.ClientSession:
        """获取或创建会话"""
        if self._session is None or self._session.closed:
            self._session = aiohttp.ClientSession(
                timeout=self.timeout,
                headers=self.default_headers
            )
        return self._session
    
    async def _build_url(self, path: str) -> str:
        """构建完整 URL"""
        path = path.lstrip('/')
        return f"{self.base_url}/{path}"
    
    async def request(
        self,
        method: str,
        path: str,
        data: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, str]] = None,
        timeout: Optional[float] = None
    ) -> HTTPResponse:
        """
        发送 HTTP 请求
        
        Args:
            method: HTTP 方法
            path: data: 请求数据（JSON）
 请求路径
                       params: URL 查询参数
            headers: 请求头
            timeout: 超时时间（覆盖默认）
            
        Returns:
            HTTPResponse 对象
        """
        session = await self._get_session()
        url = await self._build_url(path)
        
        request_headers = {**self.default_headers}
        if headers:
            request_headers.update(headers)
        
        request_kwargs: Dict[str, Any] = {
            'method': method,
            'url': url,
            'headers': request_headers,
        }
        
        if params:
            request_kwargs['params'] = params
        
        if data is not None:
            request_kwargs['json'] = data
        
        if timeout is not None:
            request_kwargs['timeout'] = aiohttp.ClientTimeout(total=timeout)
        
        async with session.request(**request_kwargs) as response:
            try:
                response_data = await response.json()
            except (aiohttp.ContentTypeError, json.JSONDecodeError):
                response_data = await response.text()
            
            return HTTPResponse(
                status=response.status,
                data=response_data,
                headers=dict(response.headers)
            )
    
    async def get(
        self,
        path: str,
        params: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, str]] = None,
        timeout: Optional[float] = None
    ) -> HTTPResponse:
        """发送 GET 请求"""
        return await self.request('GET', path, params=params, headers=headers, timeout=timeout)
    
    async def post(
        self,
        path: str,
        data: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, str]] = None,
        timeout: Optional[float] = None
    ) -> HTTPResponse:
        """发送 POST 请求"""
        return await self.request('POST', path, data=data, headers=headers, timeout=timeout)
    
    async def put(
        self,
        path: str,
        data: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, str]] = None,
        timeout: Optional[float] = None
    ) -> HTTPResponse:
        """发送 PUT 请求"""
        return await self.request('PUT', path, data=data, headers=headers, timeout=timeout)
    
    async def delete(
        self,
        path: str,
        data: Optional[Dict[str, Any]] = None,
        headers: Optional[Dict[str, str]] = None,
        timeout: Optional[float] = None
    ) -> HTTPResponse:
        """发送 DELETE 请求"""
        return await self.request('DELETE', path, data=data, headers=headers, timeout=timeout)
    
    async def close(self):
        """关闭会话"""
        if self._session and not self._session.closed:
            await self._session.close()
            self._session = None
    
    async def __aenter__(self):
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        await self.close()
