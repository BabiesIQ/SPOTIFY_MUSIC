"""Protect assistant session update handling from unknown-peer traceback storms.

VC and other Pyrogram updates remain enabled. Unknown peers in partial
channel updates were already undeliverable by Pyrogram 2.0.106.
"""
import logging
import time
import types

from pyrogram import Client

log = logging.getLogger("SPOTIFY_MUSIC.safe_updates")


async def _safe_handle_updates(self, updates):
    try:
        return await Client.handle_updates(self, updates)
    except (ValueError, KeyError) as exc:
        detail = str(exc)
        if "Peer id invalid:" not in detail and "ID not found:" not in detail:
            raise
        now = time.monotonic()
        last = getattr(self, "_spotify_peer_warning_at", 0.0)
        if now - last > 90:
            self._spotify_peer_warning_at = now
            log.warning(
                "IGNORED_UNKNOWN_PEER_UPDATE assistant=%s (further warnings rate-limited)",
                getattr(self, "_spotify_assistant_number", "?"),
            )


def patch_assistant_updates(client, number):
    if getattr(client, "_spotify_update_guard", False):
        return
    client._spotify_assistant_number = number
    client._spotify_update_guard = True
    client.handle_updates = types.MethodType(_safe_handle_updates, client)
