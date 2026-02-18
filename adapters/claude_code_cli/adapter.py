"""
Claude Code CLI Adapter

Integrates Claude Code CLI with Ziwei platform through Taibai.
"""

import asyncio
import json
import subprocess
import sys
from typing import Optional
from pathlib import Path

# Add parent directory to path for imports
sys.path.insert(0, str(Path(__file__).parent.parent.parent / "sdk" / "python"))

from ziwei_taibai.adapters.cli_base import CLIAdapterBase
from ziwei_taibai.adapters.base import Task, TaskResult, HealthStatus, AdapterConfig
from ziwei_taibai.agent import Agent


class ClaudeCodeCLIAdapter(CLIAdapterBase):
    """
    Adapter for Claude Code CLI.

    This adapter wraps the Claude Code CLI and integrates it with
    the Ziwei platform for governance and audit.
    """

    def __init__(self, config: AdapterConfig):
        cli_path = config.get("CLAUDE_CODE_CLI_PATH", "claude")
        cli_args_str = config.get("CLAUDE_CODE_CLI_ARGS", "")
        cli_args = cli_args_str.split() if cli_args_str else []

        super().__init__(config, cli_path, cli_args)

        # Initialize Taibai SDK
        self.sdk = Agent(
            owner=config.owner_id,
            tianshu_api_base=config.tianshu_api_base,
            diting_audit_url=config.diting_audit_url,
        )

    async def initialize(self) -> bool:
        """
        Initialize adapter:
        1. Discover Tianshu
        2. Register agent
        3. Start heartbeat loop
        """
        try:
            # Discover Tianshu
            discovery = self.sdk.discover()
            print(f"[ClaudeCodeCLI] Discovered Tianshu: {discovery}")

            # Register agent
            result = self.sdk.register(agent_display_id="claude-code-cli")
            print(f"[ClaudeCodeCLI] Registered: {result}")

            if not result.get("ok"):
                print(f"[ClaudeCodeCLI] Registration failed: {result}")
                return False

            # Start heartbeat loop
            self._heartbeat_task = asyncio.create_task(self._heartbeat_loop())

            self._initialized = True
            return True

        except Exception as e:
            print(f"[ClaudeCodeCLI] Initialization failed: {e}")
            return False

    async def _heartbeat_loop(self) -> None:
        """Send periodic heartbeats to Tianshu"""
        while True:
            try:
                await asyncio.sleep(self.config.heartbeat_interval)
                result = self.sdk.heartbeat()
                print(f"[ClaudeCodeCLI] Heartbeat: {result}")
            except asyncio.CancelledError:
                break
            except Exception as e:
                print(f"[ClaudeCodeCLI] Heartbeat failed: {e}")

    async def execute_task(self, task: Task) -> TaskResult:
        """
        Execute a task using Claude Code CLI.

        Args:
            task: Task to execute

        Returns:
            TaskResult with execution outcome
        """
        try:
            # Report task start to Diting
            if self.config.auto_report_actions:
                await self.report_action("task_start", {
                    "task_id": task.id,
                    "description": task.description,
                })

            # Convert task to CLI command
            command = self._task_to_command(task)

            # Execute via subprocess (simpler than interactive stdin/stdout)
            result = await self._execute_command(command, task.timeout or self.config.task_timeout)

            # Parse output
            task_result = self._parse_output(result, task)

            # Report task completion to Diting
            if self.config.auto_report_actions:
                await self.report_action("task_complete", {
                    "task_id": task.id,
                    "status": task_result.status,
                })

            return task_result

        except Exception as e:
            error_result = TaskResult(
                task_id=task.id,
                status="failed",
                error=str(e),
            )

            # Report failure
            if self.config.auto_report_actions:
                await self.report_action("task_failed", {
                    "task_id": task.id,
                    "error": str(e),
                })

            return error_result

    async def _execute_command(self, command: str, timeout: int) -> str:
        """
        Execute command using subprocess.

        Args:
            command: Command to execute
            timeout: Timeout in seconds

        Returns:
            Command output
        """
        cmd = [self.cli_path] + self.cli_args + [command]

        try:
            proc = await asyncio.create_subprocess_exec(
                *cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )

            stdout, stderr = await asyncio.wait_for(
                proc.communicate(),
                timeout=timeout,
            )

            output = stdout.decode() if stdout else ""
            error = stderr.decode() if stderr else ""

            if proc.returncode != 0:
                raise RuntimeError(f"Command failed with code {proc.returncode}: {error}")

            return output

        except asyncio.TimeoutError:
            proc.kill()
            await proc.wait()
            raise RuntimeError(f"Command timed out after {timeout} seconds")

    def _task_to_command(self, task: Task) -> str:
        """
        Convert task to Claude Code CLI command.

        For now, simple pass-through of task description.
        Future: Support structured task formats, file operations, etc.

        Args:
            task: Task to convert

        Returns:
            CLI command string
        """
        # Simple implementation: pass task description as prompt
        return task.description

    def _parse_output(self, output: str, task: Task) -> TaskResult:
        """
        Parse CLI output into TaskResult.

        Args:
            output: CLI output
            task: Original task

        Returns:
            TaskResult
        """
        # Simple implementation: treat any output as success
        return TaskResult(
            task_id=task.id,
            status="success",
            output=output,
        )

    async def report_action(self, action_type: str, detail: dict) -> None:
        """Report action to Diting for audit"""
        try:
            result = self.sdk.trace(action_type, **detail)
            print(f"[ClaudeCodeCLI] Reported action {action_type}: {result}")
        except Exception as e:
            print(f"[ClaudeCodeCLI] Failed to report action: {e}")


# Register adapter
from ziwei_taibai.adapters.registry import AdapterRegistry
AdapterRegistry.register("claude-code-cli", ClaudeCodeCLIAdapter)
