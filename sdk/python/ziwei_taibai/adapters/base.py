"""
Base adapter interface and data models for all agent types.
"""

from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from typing import Any, Dict, Optional
from enum import Enum


class HealthStatus(Enum):
    """Health status of an adapter"""
    HEALTHY = "healthy"
    DEGRADED = "degraded"
    UNHEALTHY = "unhealthy"
    UNKNOWN = "unknown"


@dataclass
class Task:
    """Task to be executed by an agent"""
    id: str
    description: str
    owner_id: str
    metadata: Dict[str, Any] = field(default_factory=dict)
    timeout: Optional[int] = None  # seconds


@dataclass
class TaskResult:
    """Result of task execution"""
    task_id: str
    status: str  # "success", "failed", "timeout"
    output: str = ""
    error: Optional[str] = None
    metadata: Dict[str, Any] = field(default_factory=dict)


@dataclass
class AdapterConfig:
    """Configuration for an adapter"""
    adapter_type: str
    owner_id: str
    tianshu_api_base: str
    diting_audit_url: Optional[str] = None
    heartbeat_interval: int = 30
    task_timeout: int = 300
    auto_report_actions: bool = True
    extra: Dict[str, Any] = field(default_factory=dict)

    def get(self, key: str, default: Any = None) -> Any:
        """Get configuration value, checking extra dict first"""
        if key in self.extra:
            return self.extra[key]
        return getattr(self, key, default)


class AgentAdapter(ABC):
    """
    Base adapter interface for all agent types.

    All adapters must implement this interface to integrate with
    the Ziwei platform through Taibai.
    """

    def __init__(self, config: AdapterConfig):
        self.config = config
        self._initialized = False

    @abstractmethod
    async def initialize(self) -> bool:
        """
        Initialize adapter and establish connection to agent.

        This should:
        - Discover Tianshu
        - Register agent
        - Start heartbeat loop
        - Perform any adapter-specific initialization

        Returns:
            True if initialization successful, False otherwise
        """
        pass

    @abstractmethod
    async def execute_task(self, task: Task) -> TaskResult:
        """
        Execute a task and return results.

        Args:
            task: Task to execute

        Returns:
            TaskResult with execution outcome
        """
        pass

    @abstractmethod
    async def health_check(self) -> HealthStatus:
        """
        Check if agent is responsive.

        Returns:
            HealthStatus indicating agent health
        """
        pass

    @abstractmethod
    async def shutdown(self) -> None:
        """
        Gracefully shutdown adapter.

        This should:
        - Stop heartbeat loop
        - Close connections
        - Clean up resources
        """
        pass

    async def report_action(self, action_type: str, detail: Dict[str, Any]) -> None:
        """
        Report action to Diting for audit.

        Override this method for custom reporting logic.

        Args:
            action_type: Type of action (e.g., "task_execution")
            detail: Action details
        """
        # Default implementation - subclasses should override if they have SDK access
        pass

    @property
    def is_initialized(self) -> bool:
        """Check if adapter has been initialized"""
        return self._initialized
