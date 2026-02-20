"""房间 API 模块

提供房间创建、加入和管理功能。
"""

from typing import Any, Dict, List, Optional
from dataclasses import dataclass

from .http_client import HTTPClient


@dataclass
class Room:
    """房间对象"""
    id: str
    name: str
    owner: str
    members: List[str]
    created_at: Optional[float] = None
    metadata: Optional[Dict[str, Any]] = None


class RoomAPI:
    """房间 API"""
    
    def __init__(self, http_client: HTTPClient):
        """
        初始化房间 API
        
        Args:
            http_client: HTTP 客户端实例
        """
        self._http = http_client
    
    async def create_room(
        self,
        name: str,
        members: Optional[List[str]] = None,
        metadata: Optional[Dict[str, Any]] = None
    ) -> Room:
        """
        创建房间
        
        Args:
            name: 房间名称
            members: 初始成员列表
            metadata: 附加元数据
            
        Returns:
            Room 对象
        """
        data = {
            "name": name,
        }
        
        if members:
            data["members"] = members
        
        if metadata:
            data["metadata"] = metadata
        
        response = await self._http.post("/api/rooms", data=data)
        
        if not response.ok:
            raise Exception(f"Failed to create room: {response.data}")
        
        room_data = response.data
        return Room(
            id=room_data.get("id", ""),
            name=room_data.get("name", ""),
            owner=room_data.get("owner", ""),
            members=room_data.get("members", []),
            created_at=room_data.get("created_at"),
            metadata=room_data.get("metadata")
        )
    
    async def get_room(self, room_id: str) -> Room:
        """
        获取房间详情
        
        Args:
            room_id: 房间 ID
            
        Returns:
            Room 对象
        """
        response = await self._http.get(f"/api/rooms/{room_id}")
        
        if not response.ok:
            raise Exception(f"Failed to get room: {response.data}")
        
        room_data = response.data
        return Room(
            id=room_data.get("id", ""),
            name=room_data.get("name", ""),
            owner=room_data.get("owner", ""),
            members=room_data.get("members", []),
            created_at=room_data.get("created_at"),
            metadata=room_data.get("metadata")
        )
    
    async def list_rooms(
        self,
        user_id: Optional[str] = None,
        limit: int = 50,
        offset: int = 0
    ) -> List[Room]:
        """
        获取房间列表
        
        Args:
            user_id: 用户 ID（可选，用于筛选用户加入的房间）
            limit: 返回数量限制
            offset: 偏移量
            
        Returns:
            Room 对象列表
        """
        params = {
            "limit": limit,
            "offset": offset,
        }
        
        if user_id:
            params["user_id"] = user_id
        
        response = await self._http.get("/api/rooms", params=params)
        
        if not response.ok:
            raise Exception(f"Failed to list rooms: {response.data}")
        
        rooms = []
        for room_data in response.data.get("rooms", []):
            rooms.append(Room(
                id=room_data.get("id", ""),
                name=room_data.get("name", ""),
                owner=room_data.get("owner", ""),
                members=room_data.get("members", []),
                created_at=room_data.get("created_at"),
                metadata=room_data.get("metadata")
            ))
        
        return rooms
    
    async def join_room(self, room_id: str, user_id: str) -> Room:
        """
        加入房间
        
        Args:
            room_id: 房间 ID
            user_id: 用户 ID
            
        Returns:
            Room 对象
        """
        response = await self._http.post(
            f"/api/rooms/{room_id}/join",
            data={"user_id": user_id}
        )
        
        if not response.ok:
            raise Exception(f"Failed to join room: {response.data}")
        
        room_data = response.data
        return Room(
            id=room_data.get("id", ""),
            name=room_data.get("name", ""),
            owner=room_data.get("owner", ""),
            members=room_data.get("members", []),
            created_at=room_data.get("created_at"),
            metadata=room_data.get("metadata")
        )
    
    async def leave_room(self, room_id: str, user_id: str) -> bool:
        """
        离开房间
        
        Args:
            room_id: 房间 ID
            user_id: 用户 ID
            
        Returns:
            是否离开成功
        """
        response = await self._http.post(
            f"/api/rooms/{room_id}/leave",
            data={"user_id": user_id}
        )
        
        if not response.ok:
            raise Exception(f"Failed to leave room: {response.data}")
        
        return True
    
    async def add_member(self, room_id: str, user_id: str) -> Room:
        """
        添加房间成员（需要房间所有者权限）
        
        Args:
            room_id: 房间 ID
            user_id: 用户 ID
            
        Returns:
            Room 对象
        """
        response = await self._http.post(
            f"/api/rooms/{room_id}/members",
            data={"user_id": user_id}
        )
        
        if not response.ok:
            raise Exception(f"Failed to add member: {response.data}")
        
        room_data = response.data
        return Room(
            id=room_data.get("id", ""),
            name=room_data.get("name", ""),
            owner=room_data.get("owner", ""),
            members=room_data.get("members", []),
            created_at=room_data.get("created_at"),
            metadata=room_data.get("metadata")
        )
    
    async def remove_member(self, room_id: str, user_id: str) -> bool:
        """
        移除房间成员（需要房间所有者权限）
        
        Args:
            room_id: 房间 ID
            user_id: 用户 ID
            
        Returns:
            是否移除成功
        """
        response = await self._http.delete(
            f"/api/rooms/{room_id}/members/{user_id}"
        )
        
        if not response.ok:
            raise Exception(f"Failed to remove member: {response.data}")
        
        return True
    
    async def update_room(
        self,
        room_id: str,
        name: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None
    ) -> Room:
        """
        更新房间信息
        
        Args:
            room_id: 房间 ID
            name: 新房间名称
            metadata: 新元数据
            
        Returns:
            Room 对象
        """
        data = {}
        
        if name:
            data["name"] = name
        
        if metadata:
            data["metadata"] = metadata
        
        response = await self._http.put(f"/api/rooms/{room_id}", data=data)
        
        if not response.ok:
            raise Exception(f"Failed to update room: {response.data}")
        
        room_data = response.data
        return Room(
            id=room_data.get("id", ""),
            name=room_data.get("name", ""),
            owner=room_data.get("owner", ""),
            members=room_data.get("members", []),
            created_at=room_data.get("created_at"),
            metadata=room_data.get("metadata")
        )
    
    async def delete_room(self, room_id: str) -> bool:
        """
        删除房间（需要房间所有者权限）
        
        Args:
            room_id: 房间 ID
            
        Returns:
            是否删除成功
        """
        response = await self._http.delete(f"/api/rooms/{room_id}")
        
        if not response.ok:
            raise Exception(f"Failed to delete room: {response.data}")
        
        return True
