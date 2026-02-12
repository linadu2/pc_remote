"""Config flow for PC Control integration."""
import logging
import aiohttp
import async_timeout
import voluptuous as vol
import asyncio

from homeassistant import config_entries
from homeassistant.const import CONF_HOST, CONF_NAME, CONF_PORT
from .const import DOMAIN, DEFAULT_PORT, DEFAULT_NAME, CONF_ACCESS_TOKEN

_LOGGER = logging.getLogger(__name__)

DATA_SCHEMA = vol.Schema({
    vol.Required(CONF_HOST): str,
    vol.Required(CONF_PORT, default=DEFAULT_PORT): int,
    vol.Required(CONF_ACCESS_TOKEN): str,
    vol.Optional(CONF_NAME, default=DEFAULT_NAME): str,
})

async def validate_input(hass, data):
    """Validate the user input allows us to connect."""
    host = data[CONF_HOST]
    port = data[CONF_PORT]
    token = data[CONF_ACCESS_TOKEN]
    
    url = f"http://{host}:{port}/volume"
    headers = {"Authorization": f"Bearer {token}"}

    try:
        async with aiohttp.ClientSession() as session:
            async with async_timeout.timeout(5):
                async with session.get(url, headers=headers) as response:
                    if response.status == 401:
                        raise InvalidAuth
                    if response.status != 200:
                        raise CannotConnect
    except (aiohttp.ClientError, asyncio.TimeoutError):
        raise CannotConnect

    return {"title": data[CONF_NAME]}

class ConfigFlow(config_entries.ConfigFlow, domain=DOMAIN):
    """Handle a config flow for PC Control."""

    VERSION = 1

    async def async_step_user(self, user_input=None):
        """Handle the initial step."""
        errors = {}
        if user_input is not None:
            try:
                info = await validate_input(self.hass, user_input)
                
                # Check if already configured
                await self.async_set_unique_id(f"{user_input[CONF_HOST]}_{user_input[CONF_PORT]}")
                self._abort_if_unique_id_configured()

                return self.async_create_entry(title=info["title"], data=user_input)
            except CannotConnect:
                errors["base"] = "cannot_connect"
            except InvalidAuth:
                errors["base"] = "invalid_auth"
            except Exception:  # pylint: disable=broad-except
                _LOGGER.exception("Unexpected exception")
                errors["base"] = "unknown"

        return self.async_show_form(
            step_id="user", data_schema=DATA_SCHEMA, errors=errors
        )

class CannotConnect(Exception):
    """Error to indicate we cannot connect."""

class InvalidAuth(Exception):
    """Error to indicate there is invalid auth."""
