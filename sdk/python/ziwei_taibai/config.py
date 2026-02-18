"""
Configuration loading and validation for adapters.
"""

import os
import yaml
from typing import Any, Dict, Optional
from .adapters.base import AdapterConfig


def load_config_from_file(config_path: str) -> AdapterConfig:
    """
    Load adapter configuration from YAML file.

    Args:
        config_path: Path to configuration file

    Returns:
        AdapterConfig instance

    Raises:
        FileNotFoundError: If config file doesn't exist
        ValueError: If config is invalid
    """
    if not os.path.exists(config_path):
        raise FileNotFoundError(f"Config file not found: {config_path}")

    with open(config_path, "r") as f:
        data = yaml.safe_load(f)

    if not data or "adapter" not in data:
        raise ValueError("Invalid config: missing 'adapter' section")

    adapter_config = data["adapter"]
    return _dict_to_config(adapter_config)


def load_config_from_env() -> AdapterConfig:
    """
    Load adapter configuration from environment variables.

    Environment variables:
    - ADAPTER_TYPE: Adapter type (required)
    - ADAPTER_OWNER_ID: Owner ID (required)
    - TIANSHU_API_BASE: Tianshu API base URL (required)
    - DITING_AUDIT_URL: Diting audit URL (optional)
    - ADAPTER_HEARTBEAT_INTERVAL: Heartbeat interval in seconds (default: 30)
    - ADAPTER_TASK_TIMEOUT: Task timeout in seconds (default: 300)
    - ADAPTER_AUTO_REPORT_ACTIONS: Auto report actions (default: true)

    Additional adapter-specific variables can be prefixed with ADAPTER_

    Returns:
        AdapterConfig instance

    Raises:
        ValueError: If required variables are missing
    """
    adapter_type = os.environ.get("ADAPTER_TYPE")
    owner_id = os.environ.get("ADAPTER_OWNER_ID")
    tianshu_api_base = os.environ.get("TIANSHU_API_BASE")

    if not adapter_type:
        raise ValueError("ADAPTER_TYPE environment variable is required")
    if not owner_id:
        raise ValueError("ADAPTER_OWNER_ID environment variable is required")
    if not tianshu_api_base:
        raise ValueError("TIANSHU_API_BASE environment variable is required")

    # Collect extra config from ADAPTER_* variables
    extra = {}
    for key, value in os.environ.items():
        if key.startswith("ADAPTER_") and key not in [
            "ADAPTER_TYPE",
            "ADAPTER_OWNER_ID",
            "ADAPTER_HEARTBEAT_INTERVAL",
            "ADAPTER_TASK_TIMEOUT",
            "ADAPTER_AUTO_REPORT_ACTIONS",
        ]:
            # Remove ADAPTER_ prefix and convert to lowercase
            config_key = key[8:]  # Remove "ADAPTER_"
            extra[config_key] = value

    return AdapterConfig(
        adapter_type=adapter_type,
        owner_id=owner_id,
        tianshu_api_base=tianshu_api_base,
        diting_audit_url=os.environ.get("DITING_AUDIT_URL"),
        heartbeat_interval=int(os.environ.get("ADAPTER_HEARTBEAT_INTERVAL", "30")),
        task_timeout=int(os.environ.get("ADAPTER_TASK_TIMEOUT", "300")),
        auto_report_actions=os.environ.get("ADAPTER_AUTO_REPORT_ACTIONS", "true").lower() == "true",
        extra=extra,
    )


def _dict_to_config(data: Dict[str, Any]) -> AdapterConfig:
    """
    Convert dictionary to AdapterConfig.

    Args:
        data: Configuration dictionary

    Returns:
        AdapterConfig instance

    Raises:
        ValueError: If required fields are missing
    """
    adapter_type = data.get("type")
    owner_id = data.get("owner_id")
    tianshu_api_base = data.get("tianshu_api_base")

    if not adapter_type:
        raise ValueError("Missing required field: type")
    if not owner_id:
        raise ValueError("Missing required field: owner_id")
    if not tianshu_api_base:
        raise ValueError("Missing required field: tianshu_api_base")

    # Extract known fields
    known_fields = {
        "type",
        "owner_id",
        "tianshu_api_base",
        "diting_audit_url",
        "heartbeat_interval",
        "task_timeout",
        "auto_report_actions",
    }

    # Everything else goes into extra
    extra = {k: v for k, v in data.items() if k not in known_fields}

    return AdapterConfig(
        adapter_type=adapter_type,
        owner_id=owner_id,
        tianshu_api_base=tianshu_api_base,
        diting_audit_url=data.get("diting_audit_url"),
        heartbeat_interval=data.get("heartbeat_interval", 30),
        task_timeout=data.get("task_timeout", 300),
        auto_report_actions=data.get("auto_report_actions", True),
        extra=extra,
    )
