"""
Adapter Manager

Manages adapter lifecycle: initialization, health monitoring, shutdown.
"""

import asyncio
from typing import Optional
from .adapters.base import AgentAdapter, HealthStatus


class AdapterManager:
    """
    Manages adapter lifecycle.

    Provides:
    - Adapter initialization
    - Health monitoring
    - Auto-restart on failure
    - Graceful shutdown
    """

    def __init__(self, adapter: AgentAdapter, auto_restart: bool = True):
        self.adapter = adapter
        self.auto_restart = auto_restart
        self._running = False
        self._health_check_task: Optional[asyncio.Task] = None

    async def start(self) -> None:
        """Start adapter and health monitoring"""
        print(f"[AdapterManager] Starting adapter: {self.adapter.__class__.__name__}")

        # Initialize adapter
        success = await self.adapter.initialize()
        if not success:
            raise RuntimeError("Adapter initialization failed")

        self._running = True

        # Start health monitoring
        self._health_check_task = asyncio.create_task(self._health_monitor())

        print(f"[AdapterManager] Adapter started successfully")

    async def stop(self) -> None:
        """Stop adapter and health monitoring"""
        print(f"[AdapterManager] Stopping adapter")

        self._running = False

        # Stop health monitoring
        if self._health_check_task:
            self._health_check_task.cancel()
            try:
                await self._health_check_task
            except asyncio.CancelledError:
                pass

        # Shutdown adapter
        await self.adapter.shutdown()

        print(f"[AdapterManager] Adapter stopped")

    async def _health_monitor(self) -> None:
        """Monitor adapter health and restart if needed"""
        while self._running:
            try:
                await asyncio.sleep(30)  # Check every 30 seconds

                health = await self.adapter.health_check()
                print(f"[AdapterManager] Health check: {health.value}")

                if health == HealthStatus.UNHEALTHY and self.auto_restart:
                    print(f"[AdapterManager] Adapter unhealthy, restarting...")
                    await self._restart()

            except asyncio.CancelledError:
                break
            except Exception as e:
                print(f"[AdapterManager] Health check failed: {e}")

    async def _restart(self) -> None:
        """Restart adapter"""
        try:
            # Shutdown
            await self.adapter.shutdown()

            # Wait a bit
            await asyncio.sleep(5)

            # Reinitialize
            success = await self.adapter.initialize()
            if not success:
                print(f"[AdapterManager] Restart failed")
            else:
                print(f"[AdapterManager] Restart successful")

        except Exception as e:
            print(f"[AdapterManager] Restart failed: {e}")
