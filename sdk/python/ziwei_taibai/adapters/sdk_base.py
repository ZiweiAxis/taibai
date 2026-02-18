"""
SDK Adapter Base for library integrations.
"""

from typing import Any
from .base import AgentAdapter, Task, TaskResult, HealthStatus, AdapterConfig


class SDKAdapterBase(AgentAdapter):
    """
    Base class for SDK-based integrations (e.g., Co-Claw).

    Subclasses should implement:
    - wrap_sdk_call: Wrap SDK calls with governance
    """

    async def wrap_sdk_call(self, method: str, *args: Any, **kwargs: Any) -> Any:
        """
        Wrap SDK calls with governance.

        This method should:
        - Report action to Diting before execution
        - Execute the SDK call
        - Report result to Diting after execution

        Args:
            method: SDK method name
            *args: Positional arguments
            **kwargs: Keyword arguments

        Returns:
            SDK call result
        """
        raise NotImplementedError("Subclasses must implement wrap_sdk_call")

    async def health_check(self) -> HealthStatus:
        """Default health check - override if needed"""
        return HealthStatus.HEALTHY if self.is_initialized else HealthStatus.UNHEALTHY

    async def shutdown(self) -> None:
        """Default shutdown - override if needed"""
        self._initialized = False
