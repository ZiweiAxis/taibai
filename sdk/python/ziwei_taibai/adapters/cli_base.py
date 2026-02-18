"""
CLI Adapter Base for subprocess-based agents.
"""

import asyncio
import subprocess
from typing import List, Optional
from .base import AgentAdapter, Task, TaskResult, HealthStatus, AdapterConfig


class CLIAdapterBase(AgentAdapter):
    """
    Base class for CLI-based agents that interact via stdin/stdout.

    Subclasses should implement:
    - _task_to_command: Convert task to CLI command
    - _parse_output: Parse CLI output to TaskResult
    """

    def __init__(self, config: AdapterConfig, cli_path: str, cli_args: Optional[List[str]] = None):
        super().__init__(config)
        self.cli_path = cli_path
        self.cli_args = cli_args or []
        self.process: Optional[subprocess.Popen] = None
        self._heartbeat_task: Optional[asyncio.Task] = None

    async def start_process(self) -> subprocess.Popen:
        """
        Start CLI process.

        Returns:
            subprocess.Popen instance
        """
        cmd = [self.cli_path] + self.cli_args
        self.process = subprocess.Popen(
            cmd,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=1,
        )
        return self.process

    async def send_command(self, command: str) -> str:
        """
        Send command to CLI and get response.

        Args:
            command: Command to send

        Returns:
            CLI output as string
        """
        if not self.process or self.process.poll() is not None:
            raise RuntimeError("CLI process not running")

        # Write command
        self.process.stdin.write(command + "\n")
        self.process.stdin.flush()

        # Read output (simple implementation - may need buffering for complex CLIs)
        output_lines = []
        # TODO: Implement proper output reading with timeout
        # For now, read until EOF or timeout
        try:
            output = self.process.stdout.read()
            return output
        except Exception as e:
            raise RuntimeError(f"Failed to read CLI output: {e}")

    async def health_check(self) -> HealthStatus:
        """Check if CLI process is running"""
        if not self.process:
            return HealthStatus.UNHEALTHY

        if self.process.poll() is None:
            return HealthStatus.HEALTHY
        else:
            return HealthStatus.UNHEALTHY

    async def shutdown(self) -> None:
        """Shutdown CLI process"""
        if self._heartbeat_task:
            self._heartbeat_task.cancel()
            try:
                await self._heartbeat_task
            except asyncio.CancelledError:
                pass

        if self.process and self.process.poll() is None:
            self.process.terminate()
            try:
                self.process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                self.process.kill()
                self.process.wait()

    def _task_to_command(self, task: Task) -> str:
        """
        Convert task to CLI command.

        Subclasses must implement this method.

        Args:
            task: Task to convert

        Returns:
            CLI command string
        """
        raise NotImplementedError("Subclasses must implement _task_to_command")

    def _parse_output(self, output: str, task: Task) -> TaskResult:
        """
        Parse CLI output into TaskResult.

        Subclasses must implement this method.

        Args:
            output: CLI output
            task: Original task

        Returns:
            TaskResult
        """
        raise NotImplementedError("Subclasses must implement _parse_output")
