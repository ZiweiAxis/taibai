# 太白 · Claude Code CLI Agent 对接说明

太白**默认提供**将 **Claude Code CLI**（以命令行形态运行的 Claude 编程助手）作为智能体对接天枢的能力。

- **方案与架构**：见根技术文档 **`ziwei/docs/open/technical/太白对接Claude-Code-CLI方案.md`**（目标、角色、数据流、实现形态、配置、验收）。
- **本仓实现**：适配器将落在 `taibai/adapters/claude_code_cli/` 或 `taibai/examples/claude_code_cli_agent/`（规划中），使用太白 SDK 完成发现、注册、心跳、指令收发与结果回传。
- **与 verification_agent 区别**：verification_agent 仅做接入验证（发现→注册→心跳→一条上报）；Claude Code CLI Agent 为**可执行任务的默认 Agent**，持续接收天枢任务、驱动 CLI 执行、按规范回传结果。
