# 太白（Taibai）子项目 Agent 约束

本文件为本子项目**自认的 Agent 行为约束**。在 Cursor、Claude Code CLI、Codex 中打开本子项目或在本目录下工作时，应自动加载并遵守本文件。若在紫微主仓（ziwei）下工作，须同时遵守根目录 **AGENTS.md**。

---

## 1. 身份与职责

- **太白**：联络 Agent 与天枢的协议与 SDK；接收天枢信息、转发 Agent，Agent 按天枢规范回传。不执行业务逻辑，不实现监听；执行在 Agent 侧，监听由谛听机制达到。
- 任务由天枢下发；太白/适配器负责将任务以约定格式交给 Agent 消费。

---

## 2. 本子项目接到的开发任务：去哪里查

- **规划侧待办**（repo/文件）：本子项目 **ISSUE_LIST**（若有）、**_bmad-output/**；若可访问紫微主仓（ziwei），则还有根仓 **docs/open/technical/**、**_bmad-output/planning-artifacts/**。约定见 `ziwei/docs/open/technical/子项目任务下发与查看约定.md`。
- 若在独立环境中执行（无主仓目录），以本子项目内 ISSUE_LIST、_bmad-output 为准；规划产出位置定义见该约定 §1。
- **运行态「要做的内容」**：当前由天枢通过 Matrix 等投递，见根仓 `子项目任务消费快速参考.md`、天枢 `tianshu/docs/任务与消费-实现状态.md`。
- 任务归属与调度见 `ziwei/docs/open/technical/子项目任务与Agent调度约定.md`。

---

## 3. 必守规约

- 代码与配置仅限 **本目录（taibai/）**；跨子项目改动由主仓或对应子项目 Agent 执行。
- 引用根技术文档时使用路径 **`ziwei/docs/open/technical/...`**；子项目 PRD/架构须引用根技术方案及子项目任务与 Agent 调度约定。
- BMAD（若启用）：产出放在 `taibai/_bmad/`、`taibai/_bmad-output/`。

---

## 4. 参考

- 子项目任务下发与查看：`ziwei/docs/open/technical/子项目任务下发与查看约定.md`
- 子项目任务与 Agent 调度：`ziwei/docs/open/technical/子项目任务与Agent调度约定.md`
- 规约与多 IDE 约定：`ziwei/docs/open/technical/规约与多IDE约定.md`
