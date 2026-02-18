# Claude Code CLI Adapter

Integrates Claude Code CLI with the Ziwei governance platform through Taibai.

## Overview

This adapter wraps the Claude Code CLI and provides:
- Automatic registration with Tianshu (communication hub)
- Periodic heartbeat reporting
- Task execution with audit trail to Diting
- Governance and compliance integration

## Quick Start

### 1. Install Dependencies

```bash
cd /path/to/taibai/sdk/python
pip install -e .
pip install pyyaml
```

### 2. Configure

Create a configuration file (or use environment variables):

```yaml
# config.yaml
adapter:
  type: "claude-code-cli"
  owner_id: "your-email@company.com"
  tianshu_api_base: "http://localhost:8082"
  diting_audit_url: "http://localhost:8080/api/audit"
  CLAUDE_CODE_CLI_PATH: "claude"
```

### 3. Run

```bash
# Using config file
python -m ziwei_taibai.adapters.claude_code_cli --config config.yaml

# Or using environment variables
export ADAPTER_TYPE=claude-code-cli
export ADAPTER_OWNER_ID=your-email@company.com
export TIANSHU_API_BASE=http://localhost:8082
export DITING_AUDIT_URL=http://localhost:8080/api/audit
export ADAPTER_CLAUDE_CODE_CLI_PATH=claude

python -m ziwei_taibai.adapters.claude_code_cli
```

## Configuration

### Required Fields

- `type`: Adapter type (must be "claude-code-cli")
- `owner_id`: Owner identifier (email or user ID)
- `tianshu_api_base`: Tianshu API base URL

### Optional Fields

- `diting_audit_url`: Diting audit API URL (for action reporting)
- `CLAUDE_CODE_CLI_PATH`: Path to claude CLI binary (default: "claude")
- `CLAUDE_CODE_CLI_ARGS`: Additional CLI arguments
- `heartbeat_interval`: Heartbeat interval in seconds (default: 30)
- `task_timeout`: Task timeout in seconds (default: 300)
- `auto_report_actions`: Auto report actions to Diting (default: true)

### Environment Variables

All configuration can be provided via environment variables:

- `ADAPTER_TYPE`: Adapter type
- `ADAPTER_OWNER_ID`: Owner ID
- `TIANSHU_API_BASE`: Tianshu API base URL
- `DITING_AUDIT_URL`: Diting audit URL
- `ADAPTER_CLAUDE_CODE_CLI_PATH`: Claude CLI path
- `ADAPTER_CLAUDE_CODE_CLI_ARGS`: Claude CLI arguments
- `ADAPTER_HEARTBEAT_INTERVAL`: Heartbeat interval
- `ADAPTER_TASK_TIMEOUT`: Task timeout
- `ADAPTER_AUTO_REPORT_ACTIONS`: Auto report actions

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  Tianshu (天枢)                          │
│            Task Distribution & Identity Hub              │
└────────────────────────┬────────────────────────────────┘
                         │ HTTP API
                         │
┌────────────────────────▼────────────────────────────────┐
│          ClaudeCodeCLIAdapter (Taibai)                   │
│  - Registration & Heartbeat                              │
│  - Task Execution                                        │
│  - Audit Reporting                                       │
└────────────────────────┬────────────────────────────────┘
                         │ subprocess
                         │
┌────────────────────────▼────────────────────────────────┐
│              Claude Code CLI                             │
│          (claude command-line tool)                      │
└──────────────────────────────────────────────────────────┘
```

## Task Execution Flow

1. **Task Received**: Adapter receives task from Tianshu
2. **Report Start**: Report task start to Diting (if auto_report_actions enabled)
3. **Execute**: Convert task to CLI command and execute
4. **Parse Output**: Parse CLI output into structured result
5. **Report Complete**: Report task completion to Diting
6. **Return Result**: Return TaskResult to caller

## Audit Trail

All actions are reported to Diting for audit:

- `task_start`: Task execution started
- `task_complete`: Task execution completed successfully
- `task_failed`: Task execution failed

Each audit record includes:
- `agent_id`: Agent identifier
- `action_type`: Action type
- `timestamp`: Unix timestamp
- `detail`: Action-specific details (task_id, status, error, etc.)

## Development

### Running Tests

```bash
cd /path/to/taibai
pytest adapters/claude_code_cli/tests/
```

### Adding Features

To extend the adapter:

1. Override `_task_to_command()` for custom task-to-command conversion
2. Override `_parse_output()` for custom output parsing
3. Override `report_action()` for custom audit reporting

## Troubleshooting

### "TIANSHU_API_BASE not set"

Ensure Tianshu API base URL is configured via config file or environment variable.

### "Registration failed"

Check that Tianshu is running and accessible at the configured URL.

### "Command timed out"

Increase `task_timeout` in configuration if tasks take longer to execute.

### "Failed to report action"

Check that Diting is running and accessible. Audit reporting failures are logged but don't fail task execution.

## See Also

- [Taibai SDK Documentation](../../sdk/python/README.md)
- [Adapter Development Guide](../../docs/adapter-development-guide.md)
- [Ziwei Technical Specification](ziwei/docs/open/technical/紫微智能体治理基础设施-技术方案.md)
