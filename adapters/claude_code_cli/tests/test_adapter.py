"""
Tests for Claude Code CLI Adapter
"""

import pytest
from unittest.mock import Mock, patch, AsyncMock
import sys
from pathlib import Path

# Add SDK and adapters to path
sys.path.insert(0, str(Path(__file__).parent.parent.parent.parent / "sdk" / "python"))
sys.path.insert(0, str(Path(__file__).parent.parent.parent))

from ziwei_taibai.adapters.base import Task, AdapterConfig, HealthStatus
from claude_code_cli.adapter import ClaudeCodeCLIAdapter


@pytest.fixture
def config():
    return AdapterConfig(
        adapter_type="claude-code-cli",
        owner_id="test@example.com",
        tianshu_api_base="http://localhost:8082",
        diting_audit_url="http://localhost:8080/api/audit",
        extra={
            "CLAUDE_CODE_CLI_PATH": "claude",
            "CLAUDE_CODE_CLI_ARGS": "",
        }
    )


@pytest.fixture
def adapter(config):
    return ClaudeCodeCLIAdapter(config)


@pytest.mark.asyncio
async def test_adapter_creation(adapter):
    """Test adapter can be created"""
    assert adapter is not None
    assert adapter.cli_path == "claude"
    assert not adapter.is_initialized


@pytest.mark.asyncio
async def test_task_to_command(adapter):
    """Test task to command conversion"""
    task = Task(
        id="test-task-1",
        description="Write a hello world program",
        owner_id="test@example.com",
    )

    command = adapter._task_to_command(task)
    assert command == "Write a hello world program"


@pytest.mark.asyncio
async def test_parse_output(adapter):
    """Test output parsing"""
    task = Task(
        id="test-task-1",
        description="Test task",
        owner_id="test@example.com",
    )

    output = "Task completed successfully"
    result = adapter._parse_output(output, task)

    assert result.task_id == "test-task-1"
    assert result.status == "success"
    assert result.output == output


@pytest.mark.asyncio
@patch('ziwei_taibai.agent.Agent.discover')
@patch('ziwei_taibai.agent.Agent.register')
async def test_initialization(mock_register, mock_discover, adapter):
    """Test adapter initialization"""
    mock_discover.return_value = {"ok": True}
    mock_register.return_value = {"ok": True, "agent_id": "test-agent-123"}

    success = await adapter.initialize()

    assert success
    assert adapter.is_initialized
    mock_discover.assert_called_once()
    mock_register.assert_called_once()


@pytest.mark.asyncio
async def test_health_check_not_initialized(adapter):
    """Test health check when not initialized"""
    health = await adapter.health_check()
    assert health == HealthStatus.UNHEALTHY


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
