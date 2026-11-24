"""Voice-assistant pool used only by the HTTP Bot API fallback.

Each active voice chat is bound to exactly one PyTgCalls instance. A failed
join or stream change tries the next configured, running assistant.
"""
import asyncio
import logging
import time
from collections import defaultdict
from SPOTIFY_MUSIC.safe_updates import patch_assistant_updates

from pyrogram import Client
from pytgcalls import PyTgCalls, StreamType

import config
from SPOTIFY_MUSIC.core.call import BABY

log = logging.getLogger("SPOTIFY_MUSIC.assistant_pool")


class AssistantPool:
    def __init__(self):
        self.clients = {}
        self.voices = {}
        self.ready = set()
        self.assigned = {}
        self.last_good = {}
        self.missing_until = {}
        self.join_limits = defaultdict(lambda: asyncio.Semaphore(12))
        self.names = {}
        self.end_callback = None
        self.left_callback = None
        self.sessions = {
            1: config.STRING1,
            2: getattr(config, "STRING2", None),
            3: getattr(config, "STRING3", None),
            4: getattr(config, "STRING4", None),
            5: getattr(config, "STRING5", None),
        }
        for number, session in self.sessions.items():
            if not session or not str(session).strip():
                continue
            if number == 1:
                client = BABY.userbot1
                voice = BABY.one
            else:
                client = Client(
                    name=f"SPOTIFY_HTTP_Assistant_{number}",
                    api_id=config.API_ID,
                    api_hash=config.API_HASH,
                    session_string=str(session),
                )
                voice = PyTgCalls(client, cache_duration=150)
            patch_assistant_updates(client, number)
            self.clients[number] = client
            self.voices[number] = voice
            self._register_events(number, voice)

    def _register_events(self, number, voice):
        @voice.on_stream_end()
        async def stream_end_handler(client, update):
            from pytgcalls.types.stream import StreamAudioEnded
            if not isinstance(update, StreamAudioEnded):
                return
            if self.assigned.get(update.chat_id) == number and self.end_callback:
                try:
                    await self.end_callback(update.chat_id)
                except Exception:
                    log.exception("STREAM_END_CALLBACK_FAILED chat=%s assistant=%s", update.chat_id, number)

        @voice.on_kicked()
        @voice.on_closed_voice_chat()
        @voice.on_left()
        async def voice_left_handler(client, chat_id):
            if self.assigned.get(chat_id) == number and self.left_callback:
                try:
                    await self.left_callback(chat_id, number)
                except Exception:
                    log.exception("VOICE_LEFT_CALLBACK_FAILED chat=%s assistant=%s", chat_id, number)

    async def start(self):
        async def start_one(number):
            try:
                await asyncio.wait_for(self.voices[number].start(), timeout=55)
                who = await asyncio.wait_for(self.clients[number].get_me(), timeout=12)
                self.ready.add(number)
                self.names[number] = who.username or str(who.id)
                log.info("VOICE_ASSISTANT_READY assistant=%d user_id=%d", number, who.id)
            except Exception as exc:
                log.warning(
                    "VOICE_ASSISTANT_UNAVAILABLE assistant=%d reason=%s",
                    number, type(exc).__name__,
                )
                try:
                    await asyncio.wait_for(self.voices[number].stop(), timeout=10)
                except Exception:
                    pass
        await asyncio.gather(*(start_one(n) for n in sorted(self.voices)))
        log.info("VOICE_POOL_READY ready=%s configured=%s",
                 len(self.ready), len(self.voices))
        return len(self.ready)

    def ordered_candidates(self, chat_id, *, change=False, exclude=None):
        exclude = set(exclude or ())
        candidates = [n for n in sorted(self.ready) if n not in exclude]
        if not candidates:
            return []
        offset = abs(chat_id) % len(candidates)
        candidates = candidates[offset:] + candidates[:offset]
        preferred = self.assigned.get(chat_id) if change else self.last_good.get(chat_id)
        if preferred in candidates:
            candidates.remove(preferred)
            candidates.insert(0, preferred)
        now = time.monotonic()
        valid = [n for n in candidates
                 if self.missing_until.get((chat_id, n), 0) <= now]
        return valid or candidates

    async def play(self, chat_id, stream_factory, *, change=False, exclude=None):
        previous = self.assigned.get(chat_id)
        last_error = None
        tried = self.ordered_candidates(chat_id, change=change, exclude=exclude)
        for number in tried:
            voice = self.voices[number]
            try:
                async with self.join_limits[number]:
                    stream = stream_factory()
                    if change and number == previous:
                        await asyncio.wait_for(
                            voice.change_stream(chat_id, stream), timeout=30
                        )
                    else:
                        await asyncio.wait_for(
                            voice.join_group_call(
                                chat_id, stream,
                                stream_type=StreamType().pulse_stream,
                            ), timeout=35,
                        )
                self.assigned[chat_id] = number
                self.last_good[chat_id] = number
                self.missing_until.pop((chat_id, number), None)
                log.info(
                    "VOICE_PLAY_STARTED chat=%s assistant=%d failover=%s",
                    chat_id, number, bool(previous and previous != number),
                )
                if previous is not None and previous != number:
                    try:
                        await asyncio.wait_for(
                            self.voices[previous].leave_group_call(chat_id),
                            timeout=7,
                        )
                    except Exception:
                        pass
                return number
            except Exception as exc:
                last_error = exc
                kind = type(exc).__name__
                log.warning(
                    "VOICE_ATTEMPT_FAILED chat=%s assistant=%d error=%s",
                    chat_id, number, kind,
                )
                if kind == "NoActiveGroupCall":
                    self.missing_until[(chat_id, number)] = time.monotonic() + 120
                if kind in ("NoAudioSourceFound", "FileNotFoundError"):
                    raise
                if number != previous and kind != "NoActiveGroupCall":
                    try:
                        await asyncio.wait_for(voice.leave_group_call(chat_id), timeout=3)
                    except Exception:
                        pass
        raise RuntimeError(
            f"All {len(tried)} available assistants failed: "
            f"{type(last_error).__name__ if last_error else 'none ready'}"
        )

    async def action(self, chat_id, method):
        number = self.assigned.get(chat_id)
        if number is None:
            raise RuntimeError("No voice assistant assigned")
        voice = self.voices[number]
        return await asyncio.wait_for(getattr(voice, method)(chat_id), timeout=20)

    async def leave(self, chat_id):
        number = self.assigned.pop(chat_id, None)
        if number is not None:
            try:
                await asyncio.wait_for(
                    self.voices[number].leave_group_call(chat_id), timeout=15
                )
            except Exception as exc:
                log.warning(
                    "VOICE_LEAVE_FAILED chat=%s assistant=%d error=%s",
                    chat_id, number, type(exc).__name__,
                )

    async def stop(self):
        for number in sorted(self.ready, reverse=True):
            try:
                await asyncio.wait_for(self.voices[number].stop(), timeout=15)
            except Exception:
                pass
        self.ready.clear()
        self.assigned.clear()
