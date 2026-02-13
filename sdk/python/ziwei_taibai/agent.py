# 太白 Agent 封装：发现、注册、心跳、操作上报
# 当前以 HTTP 调用天枢/谛听为主；与技术方案 §4 协议对齐

import os
import time
import urllib.request
import urllib.error
import json
from typing import Any, Dict, Optional

from .protocol import ACTION_VERIFICATION_PING


def _get_env(key: str, default: Optional[str] = None) -> str:
    v = os.environ.get(key, default)
    return (v or "").strip()


def discover_tianshu(api_base: Optional[str] = None) -> Dict[str, Any]:
    """发现天枢端点：GET api_base/.well-known/tianshu-matrix 或 /api/v1/discovery。"""
    base = (api_base or _get_env("TIANSHU_API_BASE", "")).rstrip("/")
    if not base:
        raise ValueError("TIANSHU_API_BASE 未设置")
    for path in ("/.well-known/tianshu-matrix", "/api/v1/discovery"):
        url = base + path
        try:
            req = urllib.request.Request(url, method="GET")
            with urllib.request.urlopen(req, timeout=10) as r:
                return json.loads(r.read().decode())
        except Exception as e:
            if path == "/api/v1/discovery":
                raise RuntimeError(f"发现天枢失败: {url}") from e
            continue
    raise RuntimeError("发现天枢失败: 未找到 discovery 端点")


def register_agent(
    api_base: str,
    owner_id: str,
    agent_display_id: Optional[str] = None,
) -> Dict[str, Any]:
    """向天枢注册 Agent（若天枢暴露 POST /api/v1/agents/register）。"""
    url = api_base.rstrip("/") + "/api/v1/agents/register"
    payload = {"owner_id": owner_id}
    if agent_display_id:
        payload["agent_display_id"] = agent_display_id
    req = urllib.request.Request(
        url,
        data=json.dumps(payload).encode(),
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=15) as r:
            return json.loads(r.read().decode())
    except urllib.error.HTTPError as e:
        body = e.read().decode() if e.fp else ""
        raise RuntimeError(f"注册失败 {e.code}: {body}") from e


def heartbeat(api_base: str, agent_id: str) -> Dict[str, Any]:
    """向天枢上报心跳（若天枢暴露 POST /api/v1/agents/heartbeat）。"""
    url = api_base.rstrip("/") + "/api/v1/agents/heartbeat"
    req = urllib.request.Request(
        url,
        data=json.dumps({"agent_id": agent_id, "status": "online"}).encode(),
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as r:
            return json.loads(r.read().decode())
    except urllib.error.HTTPError as e:
        body = e.read().decode() if e.fp else ""
        raise RuntimeError(f"心跳失败 {e.code}: {body}") from e


def report_action(
    diting_audit_url: str,
    agent_id: str,
    action_type: str,
    detail: Optional[Dict[str, Any]] = None,
) -> Dict[str, Any]:
    """向谛听上报一条操作（审计）。"""
    url = (diting_audit_url or _get_env("DITING_AUDIT_URL", "")).rstrip("/")
    if not url:
        raise ValueError("DITING_AUDIT_URL 未设置")
    payload = {
        "agent_id": agent_id,
        "action_type": action_type,
        "timestamp": int(time.time()),
        "detail": detail or {},
    }
    req = urllib.request.Request(
        url,
        data=json.dumps(payload).encode(),
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as r:
            return json.loads(r.read().decode())
    except urllib.error.HTTPError as e:
        body = e.read().decode() if e.fp else ""
        raise RuntimeError(f"审计上报失败 {e.code}: {body}") from e


class Agent:
    """太白 Agent：封装发现、注册、心跳、操作上报（与技术方案 §4 对齐）。"""

    def __init__(
        self,
        owner: str,
        tianshu_api_base: Optional[str] = None,
        diting_audit_url: Optional[str] = None,
        agent_id: Optional[str] = None,
    ):
        self.owner = owner
        self.tianshu_api_base = tianshu_api_base or _get_env("TIANSHU_API_BASE")
        self.diting_audit_url = diting_audit_url or _get_env("DITING_AUDIT_URL")
        self._agent_id = agent_id or _get_env("VERIFICATION_AGENT_ID")

    def discover(self) -> Dict[str, Any]:
        return discover_tianshu(self.tianshu_api_base)

    def register(self, agent_display_id: Optional[str] = None) -> Dict[str, Any]:
        if not self.tianshu_api_base:
            raise ValueError("tianshu_api_base 未设置")
        out = register_agent(self.tianshu_api_base, self.owner, agent_display_id)
        if out.get("ok") and out.get("agent_id"):
            self._agent_id = out["agent_id"]
        return out

    def heartbeat(self) -> Dict[str, Any]:
        if not self._agent_id:
            raise ValueError("无 agent_id，请先 register 或设置 VERIFICATION_AGENT_ID")
        return heartbeat(self.tianshu_api_base or "", self._agent_id)

    def trace(
        self,
        action_type: str,
        **detail: Any,
    ) -> Dict[str, Any]:
        if not self._agent_id:
            raise ValueError("无 agent_id")
        return report_action(
            self.diting_audit_url or "",
            self._agent_id,
            action_type,
            detail if detail else None,
        )

    @property
    def agent_id(self) -> Optional[str]:
        return self._agent_id
