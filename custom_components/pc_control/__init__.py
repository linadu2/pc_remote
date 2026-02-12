"""The PC Control integration."""
import asyncio
import logging
from homeassistant.config_entries import ConfigEntry
from homeassistant.core import HomeAssistant

from .const import DOMAIN

_LOGGER = logging.getLogger(__name__)

PLATFORMS = ["media_player"]

async def async_setup(hass: HomeAssistant, config: dict):
    """Set up the PC Control component."""
    # We return True because we don't support YAML setup via async_setup 
    # (unless we implemented a bridge, but UI flow is preferred now).
    return True

async def async_setup_entry(hass: HomeAssistant, entry: ConfigEntry):
    """Set up PC Control from a config entry."""
    # Store the config entry ID if needed (optional)
    hass.data.setdefault(DOMAIN, {})
    hass.data[DOMAIN][entry.entry_id] = entry.data

    # Forward the setup to the media_player platform
    await hass.config_entries.async_forward_entry_setups(entry, PLATFORMS)

    return True

async def async_unload_entry(hass: HomeAssistant, entry: ConfigEntry):
    """Unload a config entry."""
    unload_ok = await hass.config_entries.async_unload_platforms(entry, PLATFORMS)
    if unload_ok:
        hass.data[DOMAIN].pop(entry.entry_id)

    return unload_ok
