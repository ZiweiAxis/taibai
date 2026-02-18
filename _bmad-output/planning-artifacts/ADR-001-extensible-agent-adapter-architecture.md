# ADR-001: Extensible Agent Adapter Architecture

**Status**: Proposed
**Date**: 2026-02-15
**Context**: Taibai (太白) subproject within Ziwei platform

## Context

Taibai needs to integrate diverse agent types into the Ziwei governance platform. Current implementation includes:
- Basic Python SDK with HTTP-based communication
- Verification agent for testing
- Planned Claude Code CLI adapter

The platform must support agents with different interaction models:
- **CLI tools**: stdin/stdout, subprocess invocation (Claude Code CLI)
- **Platform integrations**: Plugin systems (OpenClaw, Dify)
- **SDK-based**: Direct library integration (Co-Claw)
- **RPA agents**: Event monitoring and reporting
- **Future agents**: WebSocket, gRPC, custom protocols

## Decision

We adopt a **layered adapter architecture** with clear separation between protocol translation and agent-specific integration:

```
┌─────────────────────────────────────────────────────────────┐
│                    Tianshu (天枢)                            │
│              Task Distribution & Identity Hub                │
└────────────────────────┬────────────────────────────────────┘
                         │ Taibai Protocol (Matrix/HTTP)
                         │
┌────────────────────────▼────────────────────────────────────┐
│                  Taibai Core SDK                             │
│  - Discovery, Registration, Heartbeat                        │
│  - Protocol: m.agent.* events                                │
│  - Audit reporting to Diting                                 │
└────────────────────────┬────────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┬──────────────┐
         │               │               │              │
┌────────▼────────┐ ┌───▼────────┐ ┌───▼──────┐ ┌────▼─────────┐
│ CLI Adapter     │ │ Plugin     │ │ SDK      │ │ Custom       │
│ Base            │ │ Adapter    │ │ Adapter  │ │ Adapter      │
│                 │ │ Base       │ │ Base     │ │ Base         │
└────────┬────────┘ └───┬────────┘ └───┬──────┘ └────┬─────────┘
         │              │              │              │
    ┌────▼─────┐   ┌───▼──────┐  ┌───▼──────┐  ┌───▼──────┐
    │ Claude   │   │ OpenClaw │  │ Co-Claw  │  │ Future   │
    │ Code CLI │   │ Dify     │  │          │  │ Agents   │
    └──────────┘   └──────────┘  └──────────┘  └──────────┘
```

## Architecture Components

### 1. Core Adapter Interface

All adapters implement a common interface:

```python
class AgentAdapter(ABC):
    """Base adapter interface for all agent types"""

    @abstractmethod
    async def initialize(self) -> bool:
        """Initialize adapter and establish connection to agent"""
        pass

    @abstractmethod
    async def execute_task(self, task: Task) -> TaskResult:
        """Execute a task and return results"""
        pass

    @abstractmethod
    async def health_check(self) -> HealthStatus:
        """Check if agent is responsive"""
        pass

    @abstractmethod
    async def shutdown(self) -> None:
        """Gracefully shutdown adapter"""
        pass

    # Optional: Override for custom action reporting
    async def report_action(self, action_type: str, detail: dict) -> None:
        """Report action to Diting for audit"""
        await self.sdk.trace(action_type, **detail)
```

### 2. Adapter Base Classes

**CLI Adapter Base** (for subprocess-based agents):
```python
class CLIAdapterBase(AgentAdapter):
    """Base for CLI-based agents (stdin/stdout interaction)"""

    def __init__(self, cli_path: str, cli_args: List[str]):
        self.cli_path = cli_path
        self.cli_args = cli_args
        self.process = None

    async def start_process(self) -> subprocess.Popen:
        """Start CLI process"""
        pass

    async def send_command(self, command: str) -> str:
        """Send command to CLI and get response"""
        pass

    async def parse_output(self, output: str) -> TaskResult:
        """Parse CLI output into structured result"""
        pass
```

**Plugin Adapter Base** (for platform integrations):
```python
class PluginAdapterBase(AgentAdapter):
    """Base for plugin-based integrations"""

    async def register_hooks(self) -> None:
        """Register hooks in target platform"""
        pass

    async def intercept_event(self, event: PlatformEvent) -> None:
        """Intercept platform events and convert to Taibai protocol"""
        pass
```

**SDK Adapter Base** (for library integrations):
```python
class SDKAdapterBase(AgentAdapter):
    """Base for SDK-based integrations"""

    async def wrap_sdk_call(self, method: str, *args, **kwargs) -> Any:
        """Wrap SDK calls with governance"""
        pass
```

### 3. Claude Code CLI Adapter Implementation

```python
class ClaudeCodeCLIAdapter(CLIAdapterBase):
    """Adapter for Claude Code CLI"""

    def __init__(self, config: AdapterConfig):
        super().__init__(
            cli_path=config.get("CLAUDE_CODE_CLI_PATH", "claude"),
            cli_args=config.get("CLAUDE_CODE_CLI_ARGS", "").split()
        )
        self.sdk = Agent(
            owner=config.get("CLAUDE_CODE_AGENT_OWNER_ID"),
            tianshu_api_base=config.get("TIANSHU_API_BASE"),
            diting_audit_url=config.get("DITING_AUDIT_URL")
        )

    async def initialize(self) -> bool:
        # Discover Tianshu
        discovery = self.sdk.discover()

        # Register agent
        result = self.sdk.register(agent_display_id="claude-code-cli")

        # Start heartbeat loop
        asyncio.create_task(self._heartbeat_loop())

        return True

    async def execute_task(self, task: Task) -> TaskResult:
        # Convert task to CLI prompt
        prompt = self._task_to_prompt(task)

        # Execute via CLI
        output = await self.send_command(prompt)

        # Parse and report
        result = await self.parse_output(output)

        # Report action to Diting
        await self.report_action("task_execution", {
            "task_id": task.id,
            "status": result.status
        })

        return result

    def _task_to_prompt(self, task: Task) -> str:
        """Convert Tianshu task to Claude Code CLI prompt"""
        return task.description  # Simple pass-through for now
```

### 4. Adapter Registry & Factory

```python
class AdapterRegistry:
    """Registry for adapter types"""

    _adapters: Dict[str, Type[AgentAdapter]] = {}

    @classmethod
    def register(cls, name: str, adapter_class: Type[AgentAdapter]):
        cls._adapters[name] = adapter_class

    @classmethod
    def create(cls, name: str, config: AdapterConfig) -> AgentAdapter:
        if name not in cls._adapters:
            raise ValueError(f"Unknown adapter: {name}")
        return cls._adapters[name](config)

# Register built-in adapters
AdapterRegistry.register("claude-code-cli", ClaudeCodeCLIAdapter)
```

### 5. Configuration Schema

Convention-based configuration with minimal required fields:

```yaml
# adapter-config.yaml
adapter:
  type: "claude-code-cli"  # Adapter type from registry

  # Common fields (all adapters)
  owner_id: "user@company.com"
  tianshu_api_base: "http://localhost:8082"
  diting_audit_url: "http://localhost:8083/api/audit"

  # Adapter-specific fields
  cli_path: "claude"
  cli_args: "--model claude-3-5-sonnet"

  # Optional
  heartbeat_interval: 30
  task_timeout: 300
  auto_report_actions: true
```

## Implementation Plan

### Phase 1: Core Infrastructure (Week 1-2)
1. Define `AgentAdapter` interface
2. Implement base classes: `CLIAdapterBase`, `PluginAdapterBase`, `SDKAdapterBase`
3. Create `AdapterRegistry` and factory
4. Add configuration loading and validation

### Phase 2: Claude Code CLI Adapter (Week 2-3)
1. Implement `ClaudeCodeCLIAdapter` using `CLIAdapterBase`
2. Add task-to-prompt conversion
3. Implement output parsing
4. Add action reporting integration
5. Create example configuration

### Phase 3: Adapter Lifecycle Management (Week 3-4)
1. Implement adapter manager for lifecycle control
2. Add health monitoring and auto-restart
3. Implement graceful shutdown
4. Add metrics and logging

### Phase 4: Documentation & Examples (Week 4)
1. Create adapter development guide
2. Document configuration schema
3. Provide example adapters for other types
4. Update integration testing

## Extension Points

**Adding a New Adapter**:

1. Choose appropriate base class or implement `AgentAdapter` directly
2. Implement required methods
3. Register in `AdapterRegistry`
4. Add configuration schema
5. Document integration steps

**Example: WebSocket-based Agent**:

```python
class WebSocketAdapterBase(AgentAdapter):
    """Base for WebSocket-based agents"""

    async def connect(self, url: str) -> None:
        self.ws = await websockets.connect(url)

    async def execute_task(self, task: Task) -> TaskResult:
        await self.ws.send(json.dumps(task.to_dict()))
        response = await self.ws.recv()
        return TaskResult.from_json(response)
```

## Benefits

1. **Simple Integration**: One command to start any adapter type
2. **Extensible**: Clear extension points for new agent types
3. **Consistent**: All agents use same protocol with Tianshu/Diting
4. **Maintainable**: Shared code in base classes, agent-specific in implementations
5. **Testable**: Mock adapters for testing without real agents

## Consequences

**Positive**:
- Unified approach to agent integration
- Easy to add new agent types
- Clear separation of concerns
- Reusable components

**Negative**:
- Initial overhead to implement base classes
- May need refactoring as new patterns emerge
- Requires documentation for adapter developers

## Directory Structure

```
taibai/
├── sdk/
│   └── python/
│       └── ziwei_taibai/
│           ├── __init__.py
│           ├── protocol.py          # Existing
│           ├── agent.py             # Existing core SDK
│           ├── adapters/
│           │   ├── __init__.py
│           │   ├── base.py          # AgentAdapter interface
│           │   ├── cli_base.py      # CLIAdapterBase
│           │   ├── plugin_base.py   # PluginAdapterBase
│           │   ├── sdk_base.py      # SDKAdapterBase
│           │   └── registry.py      # AdapterRegistry
│           ├── config.py            # Configuration loading
│           └── manager.py           # Adapter lifecycle manager
├── adapters/
│   ├── claude_code_cli/
│   │   ├── __init__.py
│   │   ├── adapter.py               # ClaudeCodeCLIAdapter
│   │   ├── config.yaml.example
│   │   └── README.md
│   ├── openclaw/                    # Future
│   ├── dify/                        # Future
│   └── template/                    # Template for new adapters
└── examples/
    ├── verification_agent/          # Existing
    └── custom_adapter_example/      # Example of custom adapter
```

## One-Command Integration

Users can start any adapter with:

```bash
# Using CLI
taibai-adapter start --type claude-code-cli --config adapter-config.yaml

# Or using Python
python -m ziwei_taibai.adapters.claude_code_cli --config adapter-config.yaml

# Or with environment variables (no config file needed)
export ADAPTER_TYPE=claude-code-cli
export CLAUDE_CODE_CLI_PATH=claude
export TIANSHU_API_BASE=http://localhost:8082
taibai-adapter start
```

## Comparison with Alternatives

| Approach | Pros | Cons | Decision |
|----------|------|------|----------|
| **Monolithic SDK** | Simple, all-in-one | Hard to extend, bloated | ❌ Rejected |
| **Plugin System** | Very flexible | Complex, runtime loading issues | ❌ Too complex |
| **Base Classes** | Balance of structure and flexibility | Requires inheritance | ✅ **Chosen** |
| **Composition** | More flexible than inheritance | More boilerplate | 🔄 Consider for v2 |

## References

- Root Technical Spec: `ziwei/docs/open/technical/紫微智能体治理基础设施-技术方案.md`
- Claude Code CLI Integration: `ziwei/docs/open/technical/太白对接Claude-Code-CLI方案.md`
- Existing Taibai SDK: `ziwei/taibai/sdk/python/ziwei_taibai/`

---

**Version**: 1.0
**Author**: BMAD Architect Agent
**Review Status**: Pending
