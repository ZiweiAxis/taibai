"""
Taibai Adapter Framework

Extensible adapter architecture for integrating diverse agent types
into the Ziwei governance platform.
"""

from .base import AgentAdapter, Task, TaskResult, HealthStatus, AdapterConfig
from .registry import AdapterRegistry
from .cli_base import CLIAdapterBase
from .plugin_base import PluginAdapterBase
from .sdk_base import SDKAdapterBase

__all__ = [
    "AgentAdapter",
    "Task",
    "TaskResult",
    "HealthStatus",
    "AdapterConfig",
    "AdapterRegistry",
    "CLIAdapterBase",
    "PluginAdapterBase",
    "SDKAdapterBase",
]
