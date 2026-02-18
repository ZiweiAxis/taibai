"""
Plugin Adapter Base for platform integrations.
"""

from typing import Any, Dict
from .base import AgentAdapter, Task, TaskResult, HealthStatus, AdapterConfig


class PluginAdapterBase(AgentAdapter):
    """
    Base class for plugin-based integrations (e.g., OpenClaw, Dify).

    Subclasses should implement:
    - register_hooks: Register hooks in target platform
    - intercept_event: Handle platform events
    """

    async def register_hooks(self) -> None:
        """
        Register hooks in target platform.

        This method should set up event listeners, webhooks, or
        other mechanisms to intercept platform events.
        """
        raise NotImplementedError("Subclasses must implement register_hooks")

    async def intercept_event(self, event: Dict[str, Any]) -> None:
        """
        Intercept platform events and convert to Taibai protocol.

        Args:
            event: Platform-specific event
        """
        raise NotImplementedError("Subclasses must implement intercept_event")

    async def health_check(self) -> HealthStatus:
        """Default health check - override if needed"""
        return HealthStatus.HEALTHY if self.is_initialized else HealthStatus.UNHEALTHY

    async def shutdown(self) -> None:
        """Default shutdown - override if needed"""
        self._initialized = False
