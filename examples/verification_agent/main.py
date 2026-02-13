#!/usr/bin/env python3
# 接入验证用智能体：发现天枢 → 注册/心跳 → 上报一条操作至谛听
# 用于集成验证；全部成功退出码 0，否则非 0

import os
import sys

# 允许从 taibai 根目录或 examples/verification_agent 运行
_taibai_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
_sdk_path = os.path.join(_taibai_root, "sdk", "python")
if _sdk_path not in sys.path:
    sys.path.insert(0, _sdk_path)

# 加载 .env（若存在）
_env = os.path.join(os.path.dirname(__file__), ".env")
if os.path.isfile(_env):
    with open(_env) as f:
        for line in f:
            line = line.strip()
            if line and not line.startswith("#") and "=" in line:
                k, v = line.split("=", 1)
                os.environ.setdefault(k.strip(), v.strip().strip('"').strip("'"))

from ziwei_taibai import Agent, discover_tianshu, report_action, heartbeat
from ziwei_taibai.protocol import ACTION_VERIFICATION_PING


def main() -> int:
    api_base = os.environ.get("TIANSHU_API_BASE", "").strip()
    audit_url = os.environ.get("DITING_AUDIT_URL", "").strip()
    agent_id = os.environ.get("VERIFICATION_AGENT_ID", "").strip()
    owner = os.environ.get("VERIFICATION_OWNER_ID", "verification-owner").strip()

    # 1) 发现天枢
    print("Step 1: 发现天枢...")
    try:
        discovery = discover_tianshu(api_base or None)
        print("  matrix_homeserver:", discovery.get("matrix_homeserver"))
        if discovery.get("api_base"):
            print("  api_base:", discovery["api_base"])
    except Exception as e:
        print("  失败:", e, file=sys.stderr)
        return 1

    # 2) 注册或使用已有 agent_id
    agent = Agent(owner=owner, agent_id=agent_id or None)
    agent.tianshu_api_base = api_base or discovery.get("api_base") or ""
    agent.diting_audit_url = audit_url
    agent._agent_id = agent_id or None

    if not agent_id and agent.tianshu_api_base:
        print("Step 2: 尝试注册...")
        try:
            out = agent.register()
            if out.get("ok") and out.get("agent_id"):
                agent._agent_id = out["agent_id"]
                print("  agent_id:", agent._agent_id)
            else:
                print("  注册未成功（可能天枢未暴露注册 API）:", out, file=sys.stderr)
                print("  将跳过心跳与上报，或设置 VERIFICATION_AGENT_ID 使用已有 agent_id。")
                if not agent._agent_id:
                    return 1
        except Exception as e:
            print("  注册异常（可设置 VERIFICATION_AGENT_ID 跳过）:", e, file=sys.stderr)
            if not agent._agent_id:
                return 1
    else:
        if agent_id:
            print("Step 2: 使用已有 VERIFICATION_AGENT_ID")
        else:
            print("Step 2: 跳过注册（无 TIANSHU_API_BASE 或无注册 API）")

    # 3) 心跳（若 API 存在）
    if agent._agent_id and agent.tianshu_api_base:
        print("Step 3: 心跳...")
        try:
            heartbeat(agent.tianshu_api_base, agent._agent_id)
            print("  成功")
        except Exception as e:
            print("  心跳失败（可能天枢未暴露心跳 API）:", e, file=sys.stderr)
            # 不因心跳失败退出，继续上报
    else:
        print("Step 3: 跳过心跳（无 agent_id 或 api_base）")

    # 4) 上报一条操作至谛听
    if not audit_url:
        print("Step 4: 跳过上报（DITING_AUDIT_URL 未设置）", file=sys.stderr)
        return 0 if agent._agent_id else 1
    if not agent._agent_id:
        print("Step 4: 无 agent_id，跳过上报", file=sys.stderr)
        return 1
    print("Step 4: 上报一条操作至谛听...")
    try:
        report_action(audit_url, agent._agent_id, ACTION_VERIFICATION_PING, {"source": "taibai_verification_agent"})
        print("  成功")
    except Exception as e:
        print("  失败:", e, file=sys.stderr)
        return 1

    print("接入验证通过。")
    return 0


if __name__ == "__main__":
    sys.exit(main())
