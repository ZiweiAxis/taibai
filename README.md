# 太白（Taibai）

紫微智能体**接入规范与 SDK** 子模块：协议封装、多语言 SDK 雏形、**接入验证用智能体**。

## 定位

- **协议**：以根 `ziwei/docs/open/technical/紫微智能体治理基础设施-技术方案.md` §4（太白协议）为权威，本仓提供实现侧封装与示例。
- **目标**：降低智能体接入紫微（天枢 + 谛听）门槛；配合**一个验证用智能体**做**集成接入验证**（发现天枢、注册/心跳、上报审计）。

## 目录结构

```
taibai/
├── README.md
├── docs/
│   ├── 接入验证指南.md      # 如何用验证智能体做天枢+谛听集成验证
│   └── 协议与SDK说明.md     # 协议事件类型、SDK 用法（指向技术方案 §4）
├── sdk/
│   └── python/              # Python SDK 雏形
│       ├── ziwei_taibai/
│       │   ├── __init__.py
│       │   ├── protocol.py  # 事件类型常量、载荷约定
│       │   └── agent.py     # Agent 封装：discovery、register、heartbeat、trace
│       └── pyproject.toml
├── examples/
│   └── verification_agent/  # 接入验证用智能体（配合集成验证）
│       ├── README.md
│       ├── main.py          # 发现 → 注册/心跳 → 上报一条 action
│       └── requirements.txt
└── .env.example
```

## 快速开始：接入验证

1. 确保 **天枢（tianshu）**、**谛听（diting）** 已部署并配置（见各子项目文档）。
2. 配置 `examples/verification_agent/.env`：`TIANSHU_API_BASE` 或 `MATRIX_HOMESERVER`、`DITING_AUDIT_URL` 等（见 `docs/接入验证指南.md`）。
3. 运行验证智能体：
   ```bash
   cd taibai/examples/verification_agent && pip install -r requirements.txt && python main.py
   ```
4. 脚本将依次：发现天枢端点 →（可选）注册/心跳 → 上报一条操作至谛听；全部成功则退出码 0，用于 CI/集成验证。

## 依赖

- Python 3.8+
- 天枢暴露发现接口（`/.well-known/tianshu-matrix` 或 `/api/v1/discovery`）；可选：注册/心跳 HTTP API。
- 谛听暴露审计上报 URL（`DITING_AUDIT_URL`）。

## 参考

- 技术方案（太白 §4）：`ziwei/docs/open/technical/紫微智能体治理基础设施-技术方案.md`
- 根架构与太白边界：`ziwei/_bmad-output/planning-artifacts/architecture-ziwei.md`、`太白边界决策.md`（已更新为建立太白子模块）
