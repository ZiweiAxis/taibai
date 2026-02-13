# 协议与 SDK 说明

太白协议与 SDK 的**权威定义**见根技术方案；本文为子模块内实现与用法摘要。

## 1. 协议权威出处

- **文档**：`ziwei/docs/open/technical/紫微智能体治理基础设施-技术方案.md` §4（太白子系统）。
- **事件类型**（基于 Matrix）：

| 事件类型 | 方向 | 用途 |
|---------|------|------|
| `m.agent.register_request` | 智能体 → 天枢 | 发起注册，携带 owner、环境指纹 |
| `m.agent.identity` | 天枢 → 智能体 | 下发 DID 及凭证 |
| `m.agent.action` | 智能体 → 谛听 | 操作上报 |
| `m.agent.audit` | 谛听 → 智能体 | 存证回执 |
| `m.agent.heartbeat` | 智能体 → 天枢 | 保活、环境指纹校验 |
| `m.agent.revoke` | 天枢 → 智能体 | 吊销通知 |

负载需智能体私钥签名，接收方通过链上 DID 公钥验签。

## 2. 本仓 SDK 范围

- **sdk/python/ziwei_taibai/**：Python 雏形
  - `protocol.py`：事件类型常量、载荷结构（与 §4 一致）。
  - `agent.py`：`Agent` 类，封装发现、注册、心跳、操作上报（当前以 HTTP 调用天枢/谛听为主；Matrix 事件可后续扩展）。
- **examples/verification_agent**：接入验证用智能体，依赖 SDK 完成发现 → 注册/心跳 → 上报一条 action。

## 3. 使用示例（与技术方案 §4.3 对齐）

```python
from ziwei_taibai import Agent
from ziwei_taibai.protocol import EVENT_ACTION

agent = Agent(
    owner="user@company.com",
    tianshu_api_base="https://tianshu.example.com",
    diting_audit_url="https://diting.example.com/api/audit",
)
agent.register()           # 触发注册流程（若 API 支持）
agent.heartbeat()          # 保活
agent.trace(EVENT_ACTION, action_type="file_write", path="/data/example.txt")
```

验证智能体见 `examples/verification_agent/main.py`。
