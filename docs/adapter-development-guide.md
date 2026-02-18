# Adapter Development Guide

This guide explains how to create custom adapters for the Taibai adapter framework.

## Overview

The Taibai adapter framework provides a flexible architecture for integrating diverse agent types into the Ziwei governance platform. Adapters translate between agent-specific protocols and the Ziwei platform protocol.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Tianshu (天枢)                            │
│              Task Distribution & Identity Hub                │
└────────────────────────┬────────────────────────────────────┘
                         │ Taibai Protocol
                         │
┌────────────────────────▼────────────────────────────────────┐
│                  Taibai Core SDK                             │
│  - Discovery, Registration, Heartbeat                        │
│  - Audit reporting to Diting                                 │
└────────────────────────┬────────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┬──────────────┐
         │               │               │              │
┌────────▼────────┐ ┌───▼────────┐ ┌───▼──────┐ ┌────▼─────────┐
│ Your Custom     │ │ CLI        │ │ Plugin   │ │ SDK          │
│ Adapter         │ │ Adapter    │ │ Adapter  │ │ Adapter      │
└────────┬────────┘ └───┬────────┘ └───┬──────┘ └────┬─────────┘
         │              │              │              │
    ┌────▼─────┐   ┌───▼──────┐  ┌───▼──────┐  ┌───▼──────┐
    │ Your     │   │ Claude   │  │ OpenClaw │  │ Co-Claw  │
    │ Agent    │   │ Code CLI │  │ Dify     │  │          │
    └──────────┘   └──────────┘  └──────────┘  └──────────┘
```

## Quick Start

### 1. Choose Base Class

Select the appropriate base class for your agent type:

- **CLIAdapterBase**: For command-line tools (stdin/stdout, subprocess)
- **PluginAdapterBase**: For platform integrations (hooks, webhooks)
- **SDKAdapterBase**: For library integrations (direct API calls)
- **AgentAdapter**: For custom protocols (implement from scratch)

### 2. Implement Required Methods

All adapters must implement:

```python
from ziwei_taibai.adapters.base import AgentAdapter, Task, TaskResult, HealthStatus, AdapterConfig

class MyAdapter(AgentAdapter):
    async def initialize(self) -> bool:
        """Initialize adapter and register with Tianshu"""
        pass

    async def execute_task(self, task: Task) -> TaskResult:
        """Execute a task"""
        pass

    async def health_check(self) -> HealthStatus:
        """Check adapter health"""
        pass

    async def shutdown(self) -> None:
        """Gracefully shutdown"""
        pass
```

### 3. Register Adapter

```python
from ziwei_taibai.adapters.registry import AdapterRegistry

AdapterRegistry.register("my-adapter", MyAdapter)
```

### 4. Create Configuration

```yaml
# config.yaml
adapter:
  type: "my-adapter"
  owner_id: "user@company.com"
  tianshu_api_base: "http://localhost:8082"
  diting_audit_url: "http://localhost:8080/api/audit"

  # Your adapter-specific fields
  my_custom_field: "value"
```

## Example: CLI Adapter

Here's a complete example of a CLI-based adapter:

```python
from ziwei_taibai.adapters.cli_base import CLIAdapterBase
from ziwei_taibai.adapters.base import Task, TaskResult, AdapterConfig
from ziwei_taibai.agent import Agent
import asyncio

class MyCliAdapter(CLIAdapterBase):
    def __init__(self, config: AdapterConfig):
        cli_path = config.get("CLI_PATH", "my-cli")
        cli_args = config.get("CLI_ARGS", "").split()

        super().__init__(config, cli_path, cli_args)

        # Initialize Taibai SDK
        self.sdk = Agent(
            owner=config.owner_id,
            tianshu_api_base=config.tianshu_api_base,
            diting_audit_url=config.diting_audit_url,
        )

    async def initialize(self) -> bool:
        # Discover and register
        self.sdk.discover()
        result = self.sdk.register(agent_display_id="my-cli-agent")

        if not result.get("ok"):
            return False

        # Start heartbeat
        self._heartbeat_task = asyncio.create_task(self._heartbeat_loop())
        self._initialized = True
        return True

    async def _heartbeat_loop(self):
        while True:
            await asyncio.sleep(self.config.heartbeat_interval)
            self.sdk.heartbeat()

    async def execute_task(self, task: Task) -> TaskResult:
        # Report start
        await self.report_action("task_start", {"task_id": task.id})

        # Execute
        command = self._task_to_command(task)
        output = await self._execute_command(command, task.timeout)
        result = self._parse_output(output, task)

        # Report complete
        await self.report_action("task_complete", {
            "task_id": task.id,
            "status": result.status
        })

        return result

    def _task_to_command(self, task: Task) -> str:
        # Convert task to CLI command
        return task.description

    def _parse_output(self, output: str, task: Task) -> TaskResult:
        # Parse CLI output
        return TaskResult(
            task_id=task.id,
            status="success",
            output=output,
        )

    async def report_action(self, action_type: str, detail: dict):
        self.sdk.trace(action_type, **detail)

# Register
from ziwei_taibai.adapters.registry import AdapterRegistry
AdapterRegistry.register("my-cli", MyCliAdapter)
```

## Example: Plugin Adapter

For platform integrations:

```python
from ziwei_taibai.adapters.plugin_base import PluginAdapterBase
from ziwei_taibai.adapters.base import Task, TaskResult, AdapterConfig
from ziwei_taibai.agent import Agent

class MyPluginAdapter(PluginAdapterBase):
    def __init__(self, config: AdapterConfig):
        super().__init__(config)
        self.sdk = Agent(
            owner=config.owner_id,
            tianshu_api_base=config.tianshu_api_base,
            diting_audit_url=config.diting_audit_url,
        )

    async def initialize(self) -> bool:
        # Register with Tianshu
        self.sdk.discover()
        result = self.sdk.register(agent_display_id="my-plugin")

        if not result.get("ok"):
            return False

        # Register hooks in target platform
        await self.register_hooks()

        self._initialized = True
        return True

    async def register_hooks(self):
        # Register webhooks, event listeners, etc.
        pass

    async def intercept_event(self, event: dict):
        # Handle platform events
        # Convert to Taibai protocol and report
        await self.report_action("event_intercepted", event)

    async def execute_task(self, task: Task) -> TaskResult:
        # Execute task via platform API
        pass

    async def report_action(self, action_type: str, detail: dict):
        self.sdk.trace(action_type, **detail)
```

## Configuration

### Required Fields

All adapters require:
- `type`: Adapter type name
- `owner_id`: Owner identifier
- `tianshu_api_base`: Tianshu API URL

### Optional Fields

- `diting_audit_url`: Diting audit URL
- `heartbeat_interval`: Heartbeat interval (default: 30s)
- `task_timeout`: Task timeout (default: 300s)
- `auto_report_actions`: Auto report to Diting (default: true)

### Custom Fields

Add adapter-specific fields to the `extra` dict:

```python
config.get("MY_CUSTOM_FIELD", "default_value")
```

## Testing

Create tests for your adapter:

```python
import pytest
from ziwei_taibai.adapters.base import Task, AdapterConfig
from my_adapter import MyAdapter

@pytest.mark.asyncio
async def test_adapter_initialization():
    config = AdapterConfig(
        adapter_type="my-adapter",
        owner_id="test@example.com",
        tianshu_api_base="http://localhost:8082",
    )

    adapter = MyAdapter(config)
    success = await adapter.initialize()

    assert success
    assert adapter.is_initialized

@pytest.mark.asyncio
async def test_task_execution():
    # ... test task execution
    pass
```

## Best Practices

1. **Error Handling**: Always handle errors gracefully and report failures
2. **Timeouts**: Respect task timeouts to prevent hanging
3. **Audit Trail**: Report all significant actions to Diting
4. **Health Checks**: Implement meaningful health checks
5. **Graceful Shutdown**: Clean up resources properly
6. **Logging**: Use print() or logging for debugging
7. **Configuration**: Use convention over configuration

## Directory Structure

```
taibai/
├── adapters/
│   ├── my_adapter/
│   │   ├── __init__.py
│   │   ├── adapter.py           # Adapter implementation
│   │   ├── config.yaml.example  # Configuration example
│   │   ├── README.md            # Adapter documentation
│   │   └── tests/               # Tests
│   │       └── test_adapter.py
```

## See Also

- [Claude Code CLI Adapter](../adapters/claude_code_cli/README.md) - Complete example
- [Base Classes](../sdk/python/ziwei_taibai/adapters/) - API reference
- [Ziwei Technical Specification](ziwei/docs/open/technical/紫微智能体治理基础设施-技术方案.md)
