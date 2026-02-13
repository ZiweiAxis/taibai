# 太白协议常量（与紫微技术方案 §4 对齐）
# 权威定义：ziwei/docs/open/technical/紫微智能体治理基础设施-技术方案.md

# Matrix 事件类型（太白扩展）
EVENT_REGISTER_REQUEST = "m.agent.register_request"
EVENT_IDENTITY = "m.agent.identity"
EVENT_ACTION = "m.agent.action"
EVENT_AUDIT = "m.agent.audit"
EVENT_HEARTBEAT = "m.agent.heartbeat"
EVENT_REVOKE = "m.agent.revoke"

# 操作类型示例（上报 m.agent.action 时使用）
ACTION_FILE_WRITE = "file_write"
ACTION_FILE_READ = "file_read"
ACTION_API_CALL = "api_call"
ACTION_VERIFICATION_PING = "verification_ping"  # 接入验证用
