"""
Adapter Registry and Factory.
"""

from typing import Dict, Type
from .base import AgentAdapter, AdapterConfig


class AdapterRegistry:
    """
    Registry for adapter types.

    Adapters can be registered and created by name.
    """

    _adapters: Dict[str, Type[AgentAdapter]] = {}

    @classmethod
    def register(cls, name: str, adapter_class: Type[AgentAdapter]) -> None:
        """
        Register an adapter type.

        Args:
            name: Adapter name (e.g., "claude-code-cli")
            adapter_class: Adapter class
        """
        cls._adapters[name] = adapter_class

    @classmethod
    def create(cls, name: str, config: AdapterConfig) -> AgentAdapter:
        """
        Create an adapter instance by name.

        Args:
            name: Adapter name
            config: Adapter configuration

        Returns:
            AgentAdapter instance

        Raises:
            ValueError: If adapter name is not registered
        """
        if name not in cls._adapters:
            raise ValueError(
                f"Unknown adapter: {name}. "
                f"Available adapters: {', '.join(cls._adapters.keys())}"
            )
        return cls._adapters[name](config)

    @classmethod
    def list_adapters(cls) -> list[str]:
        """
        List all registered adapter names.

        Returns:
            List of adapter names
        """
        return list(cls._adapters.keys())

    @classmethod
    def is_registered(cls, name: str) -> bool:
        """
        Check if an adapter is registered.

        Args:
            name: Adapter name

        Returns:
            True if registered, False otherwise
        """
        return name in cls._adapters
