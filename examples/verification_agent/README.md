# 接入验证用智能体

配合**太白子模块**与紫微（天枢 + 谛听）做**集成接入验证**的示例智能体。用于常见集成场景的自动化或人工验证。

## 作用

1. **发现天枢**：GET 天枢发现端点，确认 `matrix_homeserver` / `api_base` 可访问。
2. **注册或使用已有 agent_id**：若配置了 `VERIFICATION_AGENT_ID` 则跳过注册；否则尝试调用天枢注册 API（若存在）。
3. **心跳**：向天枢上报在线（若 API 存在）。
4. **上报一条操作**：向谛听审计 URL 发送一条 `verification_ping` 操作，验证端到端链路。

全部成功退出码 0，任一步失败则非 0。

## 配置

复制 `.env.example` 为 `.env`，填写：

- `TIANSHU_API_BASE`：天枢 API 根（如 `http://localhost:8080`）。
- `DITING_AUDIT_URL`：谛听审计上报 URL。
- 可选 `VERIFICATION_AGENT_ID`：已在天枢侧注册的 agent_id，用于仅验证心跳与上报。

## 运行

```bash
# 安装依赖（含本地太白 SDK）
pip install -r requirements.txt
# 或从 taibai 根目录：pip install -e sdk/python
python main.py
echo "Exit: $?"
```

详见 `taibai/docs/接入验证指南.md`。
