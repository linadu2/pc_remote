"""Platform for PC Control Media Player."""
import logging
import asyncio
import aiohttp
import async_timeout

from homeassistant.components.media_player import (
    MediaPlayerEntity,
    MediaPlayerEntityFeature,
    MediaPlayerState,
)
from homeassistant.const import (
    CONF_HOST,
    CONF_NAME,
    CONF_PORT,
)
from .const import CONF_ACCESS_TOKEN, CONF_MAC_ADDRESS

_LOGGER = logging.getLogger(__name__)

# Note: We removed PLATFORM_SCHEMA and async_setup_platform 
# because we are now using Config Flow (async_setup_entry).

async def async_setup_entry(hass, config_entry, async_add_entities):
    """Set up the PC Control media player from a config entry."""
    config = config_entry.data
    
    host = config[CONF_HOST]
    port = config[CONF_PORT]
    name = config[CONF_NAME]
    token = config[CONF_ACCESS_TOKEN]
    mac = config.get(CONF_MAC_ADDRESS)

    entity = PCMediaPlayer(name, host, port, token, mac)
    async_add_entities([entity], update_before_add=False)


class PCMediaPlayer(MediaPlayerEntity):
    """Representation of the PC Media Player."""

    def __init__(self, name, host, port, token, mac):
        self._name = name
        self._host = host
        self._port = port
        self._token = token
        self._mac = mac
        self._volume = 0.0
        self._state = MediaPlayerState.OFF
        self._available = True
        self._session = None
        self._base_url = f"http://{self._host}:{self._port}"
        self._headers = {
            "Authorization": f"Bearer {self._token}",
            "Content-Type": "application/json"
        }
        self._attr_unique_id = f"pc_control_{host}_{port}"

    async def async_added_to_hass(self):
        self._session = aiohttp.ClientSession()

    async def async_will_remove_from_hass(self):
        if self._session:
            await self._session.close()

    @property
    def name(self):
        return self._name

    @property
    def state(self):
        return self._state

    @property
    def available(self):
        return True

    @property
    def volume_level(self):
        return self._volume

    @property
    def supported_features(self):
        features = (
            MediaPlayerEntityFeature.PLAY
            | MediaPlayerEntityFeature.PAUSE
            | MediaPlayerEntityFeature.STOP
            | MediaPlayerEntityFeature.NEXT_TRACK
            | MediaPlayerEntityFeature.PREVIOUS_TRACK
            | MediaPlayerEntityFeature.VOLUME_SET
            | MediaPlayerEntityFeature.TURN_OFF
        )
        # Only support TURN_ON if we have a MAC address
        if self._mac:
            features |= MediaPlayerEntityFeature.TURN_ON
        return features

    async def async_update(self):
        if self._session is None or self._session.closed:
            return

        try:
            url = f"{self._base_url}/volume"
            async with async_timeout.timeout(2):
                response = await self._session.get(url, headers=self._headers)
                
                if response.status == 200:
                    data = await response.json()
                    self._volume = data.get("volume", 0) / 100.0
                    self._state = MediaPlayerState.PLAYING
                else:
                    self._state = MediaPlayerState.OFF
                    
        except (aiohttp.ClientError, asyncio.TimeoutError):
            self._state = MediaPlayerState.OFF
        except Exception as e:
            _LOGGER.error("Error updating: %s", e)
            self._state = MediaPlayerState.OFF

    async def async_turn_on(self):
            """Turn the media player on."""
            if self._mac:
                await self.hass.services.async_call(
                    "wake_on_lan", "send_magic_packet", {"mac": self._mac}
                )

    async def async_turn_off(self):
        await self._send_action("turn_off")

    async def async_media_next_track(self):

        await self._send_action("next")

    async def async_media_previous_track(self):
        await self._send_action("prev")

    async def async_media_stop(self):
        await self._send_action("stop")

    async def async_media_play(self):
        await self._send_action("play_pause")

    async def async_media_pause(self):
        await self._send_action("play_pause")

    async def _send_action(self, action_name):
        if self._state == MediaPlayerState.OFF: return
        if self._session is None: return
        
        try:
            url = f"{self._base_url}/action"
            await self._session.post(
                url, 
                json={"action": action_name}, 
                headers=self._headers
            )
        except Exception as e:
            _LOGGER.error("Error sending %s: %s", action_name, e)

    async def async_set_volume_level(self, volume):
        if self._state == MediaPlayerState.OFF: return
        if self._session is None: return
        
        try:
            url = f"{self._base_url}/volume"
            vol_int = int(volume * 100)
            await self._session.post(
                url, 
                json={"volume": vol_int}, 
                headers=self._headers
            )
            self._volume = volume
        except Exception as e:
            _LOGGER.error("Error setting volume: %s", e)
