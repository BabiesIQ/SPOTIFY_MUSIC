"""Temporary Telegram HTTP Bot API frontend for SPOTIFY_MUSIC.

The bot never performs MTProto bot authorization in this mode.
Voice calls still require an already authorized user STRING_SESSION.
This is an intentionally limited, reversible compatibility frontend.
"""
import asyncio
import logging
import os
import re
import time
from collections import defaultdict, deque
from dataclasses import dataclass
from pathlib import Path
from urllib.parse import parse_qs, urlparse

import aiohttp
from telegram import InlineKeyboardButton, InlineKeyboardMarkup, Update
from telegram.ext import Application, CallbackQueryHandler, CommandHandler, ContextTypes
from pytgcalls import StreamType
from pytgcalls.types.input_stream import AudioPiped, AudioVideoPiped
from pytgcalls.types.input_stream.quality import LowQualityVideo, MediumQualityAudio
from pytgcalls.types.stream import StreamAudioEnded
from py_yt import VideosSearch

import config
from SPOTIFY_MUSIC.core.call import BABY

log = logging.getLogger("SPOTIFY_MUSIC.http_fallback")
logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(name)s: %(message)s")
logging.getLogger("httpx").setLevel(logging.WARNING)

voice = BABY.one
queues = defaultdict(deque)
locks = defaultdict(asyncio.Lock)
active = set()
voice_ready = False
assistant_name = None
bot_name = None
media_folder = Path("downloads/http_fallback")
media_folder.mkdir(parents=True, exist_ok=True)


@dataclass
class Track:
    title: str
    source: str
    video_id: str | None = None
    video: bool = False
    local_file: str | None = None
    requested_by: int = 0


def short_error(exc):
    return type(exc).__name__


def buttons():
    return InlineKeyboardMarkup([
        [
            InlineKeyboardButton("⏸ Pause", callback_data="fb:pause"),
            InlineKeyboardButton("▶️ Resume", callback_data="fb:resume"),
        ],
        [
            InlineKeyboardButton("⏭ Skip", callback_data="fb:skip"),
            InlineKeyboardButton("⏹ Stop", callback_data="fb:stop"),
        ],
    ])


def is_group(update):
    return update.effective_chat and update.effective_chat.type in ("group", "supergroup")


async def is_admin(update: Update, context: ContextTypes.DEFAULT_TYPE):
    user = update.effective_user
    if user is None:
        return False
    if user.id == config.OWNER_ID:
        return True
    if not is_group(update):
        return False
    if update.effective_message and update.effective_message.sender_chat:
        return update.effective_message.sender_chat.id == update.effective_chat.id
    try:
        member = await context.bot.get_chat_member(update.effective_chat.id, user.id)
        return member.status in ("creator", "administrator", "owner")
    except Exception as exc:
        log.warning("Cannot check admin status in chat %s: %s", update.effective_chat.id, short_error(exc))
        return False


async def answer(update: Update, text: str, *, controls=False):
    if update.effective_chat is None:
        return
    if update.callback_query:
        await update.callback_query.answer()
    kwargs = {"reply_markup": buttons()} if controls else {}
    await update.get_bot().send_message(
        update.effective_chat.id, text[:3900], disable_web_page_preview=True, **kwargs
    )


async def search_track(query: str, video: bool, sender_id: int):
    # Direct Telegram media and YouTube links are supported; search uses py_yt.
    if query.startswith(("http://", "https://")):
        parsed = urlparse(query)
        host = (parsed.hostname or "").lower()
        if not (host == "youtu.be" or host == "youtube.com" or host.endswith(".youtube.com")):
            raise ValueError("Only YouTube links or search terms are supported")
    results = await asyncio.wait_for(VideosSearch(query, limit=1).next(), timeout=22)
    items = results.get("result") or []
    if not items:
        raise ValueError("No matching video found")
    found = items[0]
    video_id = found.get("id")
    if not video_id:
        raise ValueError("YouTube result missing video ID")
    title = (found.get("title") or "Unknown track")[:110]
    return Track(
        title=title,
        source="https://www.youtube.com/watch?v=" + video_id,
        video_id=video_id,
        video=video,
        requested_by=sender_id,
    )


async def http_media_url(track: Track):
    kind = "video" if track.video else "song"
    base_url = (config.BASE_URL or "").rstrip("/")
    if not base_url or not track.video_id:
        raise ValueError("No media backend configured")
    params = {"query": track.video_id, "download": "true"}
    if getattr(config, "API_KEY", None):
        params["api"] = config.API_KEY
    timeout = aiohttp.ClientTimeout(total=30, connect=12)
    async with aiohttp.ClientSession(timeout=timeout) as session:
        async with session.get(base_url + "/api/" + kind, params=params) as response:
            if response.status != 200:
                raise ValueError("Media backend returned HTTP " + str(response.status))
            payload = await response.json(content_type=None)
    stream = payload.get("stream") if isinstance(payload, dict) else None
    if not stream or not stream.startswith(("https://", "http://")):
        raise ValueError("Media backend returned no playable stream")
    return stream


async def ytdlp_audio_url(track: Track):
    args = [
        "yt-dlp", "--no-playlist", "--no-warnings",
        "-f", "best[height<=?720]/best" if track.video else "bestaudio/best",
        "-g", track.source,
    ]
    proc = await asyncio.create_subprocess_exec(
        *args, stdout=asyncio.subprocess.PIPE, stderr=asyncio.subprocess.PIPE
    )
    try:
        out, _ = await asyncio.wait_for(proc.communicate(), timeout=35)
    except asyncio.TimeoutError:
        proc.kill()
        await proc.communicate()
        raise TimeoutError("yt-dlp resolution timed out")
    if proc.returncode or not out.strip():
        raise ValueError("yt-dlp could not resolve a playable stream")
    return out.decode("utf-8", "replace").splitlines()[0].strip()


async def stream_source(track: Track):
    if track.local_file:
        return track.local_file
    try:
        url = await http_media_url(track)
        # A backend may return a processing URL until media becomes available.
        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=12)) as session:
            for attempt in range(4):
                try:
                    async with session.get(url, headers={"Range": "bytes=0-0"}) as r:
                        if r.status in (200, 206):
                            return url
                        if r.status not in (202, 204, 404, 410, 423):
                            break
                except (aiohttp.ClientError, asyncio.TimeoutError):
                    pass
                await asyncio.sleep(2)
        # The returned URL may be an HLS stream or a signed URL; allow FFmpeg to try it.
        return url
    except Exception as exc:
        log.warning("Backend source unavailable for %s: %s; trying yt-dlp", track.video_id, short_error(exc))
        return await ytdlp_audio_url(track)


def build_stream(track: Track, path: str):
    if track.video:
        return AudioVideoPiped(
            path, audio_parameters=MediumQualityAudio(), video_parameters=LowQualityVideo()
        )
    return AudioPiped(path, audio_parameters=MediumQualityAudio())


async def start_track(chat_id: int, track: Track, *, change: bool):
    path = await stream_source(track)
    stream = build_stream(track, path)
    if change:
        await voice.change_stream(chat_id, stream)
    else:
        await voice.join_group_call(chat_id, stream, stream_type=StreamType().pulse_stream)
        active.add(chat_id)


async def next_track(chat_id: int, *, from_event=False):
    async with locks[chat_id]:
        q = queues[chat_id]
        if q:
            finished = q.popleft()
            if finished.local_file:
                try:
                    Path(finished.local_file).unlink(missing_ok=True)
                except OSError:
                    pass
        while q:
            try:
                await start_track(chat_id, q[0], change=chat_id in active)
                return q[0]
            except Exception as exc:
                log.error("Next track failed in %s: %s", chat_id, short_error(exc))
                q.popleft()
        active.discard(chat_id)
        try:
            await voice.leave_group_call(chat_id)
        except Exception:
            pass
        return None


@voice.on_stream_end()
async def on_stream_end(client, update):
    if not isinstance(update, StreamAudioEnded):
        return
    chat_id = update.chat_id
    if chat_id not in active:
        return
    await next_track(chat_id, from_event=True)


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    status = "ready" if voice_ready else "not connected"
    assistant = "@" + assistant_name if assistant_name else "the configured assistant account"
    await answer(
        update,
        "🎵 SPOTIFY_MUSIC (HTTP Bot API fallback)\n"
        "Bot commands are online without MTProto bot login.\n"
        "Voice assistant: " + status + " (" + assistant + ")\n\n"
        "In a group: /play query, /vplay query, /queue, /pause, /resume, "
        "/skip, /stop, /ping.\n"
        "Add the assistant account to the group and start its voice chat first.",
    )


async def ping(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await answer(update, "🏓 Bot API online. Voice assistant: " + ("ready" if voice_ready else "not connected"))


async def prepare_telegram_media(update: Update, video: bool):
    message = update.effective_message
    if not message or not message.reply_to_message:
        return None
    ref = message.reply_to_message
    media = ref.video if video else (ref.audio or ref.voice or ref.document)
    if media is None:
        return None
    file_size = getattr(media, "file_size", None) or 0
    if file_size and file_size > 19 * 1024 * 1024:
        raise ValueError("Telegram files above 19 MB are not supported by this temporary mode")
    path = media_folder / (str(update.effective_chat.id) + "-" + str(time.time_ns()) + (".mp4" if video else ".media"))
    f = await media.get_file()
    await f.download_to_drive(custom_path=str(path))
    return Track(
        title=getattr(media, "file_name", None) or ("Telegram video" if video else "Telegram audio"),
        source="telegram", video=video, local_file=str(path),
        requested_by=update.effective_user.id if update.effective_user else 0,
    )


async def play(update: Update, context: ContextTypes.DEFAULT_TYPE, *, video=False):
    if not is_group(update):
        return await answer(update, "Use /play or /vplay inside a Telegram group.")
    if not voice_ready:
        return await answer(update, "⚠️ HTTP bot is online, but the assistant MTProto session is not connected.")
    try:
        query = " ".join(context.args).strip()
        if query:
            track = await search_track(query, video, update.effective_user.id if update.effective_user else 0)
        else:
            track = await prepare_telegram_media(update, video)
            if track is None:
                return await answer(update, "Usage: /play song name or YouTube URL; or reply to an audio file with /play.")
    except Exception as exc:
        log.warning("Search/media failed: %s", short_error(exc))
        return await answer(update, "Could not find or load that track (" + short_error(exc) + ").")

    chat_id = update.effective_chat.id
    async with locks[chat_id]:
        q = queues[chat_id]
        q.append(track)
        if chat_id in active:
            return await answer(update, "✅ Queued #" + str(len(q)) + ": " + track.title, controls=True)
        try:
            await start_track(chat_id, track, change=False)
        except Exception as exc:
            q.pop()
            log.error("Join/play failed for %s: %s", chat_id, short_error(exc))
            return await answer(
                update, "⚠️ Could not join/play (" + short_error(exc) + "). "
                "Confirm assistant is in the group and a voice chat is active."
            )
    await answer(update, "▶️ Playing: " + track.title, controls=True)


async def play_audio(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await play(update, context, video=False)


async def play_video(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await play(update, context, video=True)


async def queue(update: Update, context: ContextTypes.DEFAULT_TYPE):
    q = queues[update.effective_chat.id]
    if not q:
        return await answer(update, "Queue is empty.")
    lines = ["🎶 Current queue:"]
    for i, item in enumerate(list(q)[:15]):
        lines.append(str(i + 1) + ". " + item.title + (" [video]" if item.video else ""))
    await answer(update, "\n".join(lines), controls=True)


async def control(update: Update, context: ContextTypes.DEFAULT_TYPE, action):
    if not is_group(update):
        return await answer(update, "Voice controls are only available in groups.")
    if not await is_admin(update, context):
        return await answer(update, "Only group admins can control playback.")
    chat_id = update.effective_chat.id
    if chat_id not in active:
        return await answer(update, "No active stream in this chat.")
    try:
        if action == "pause":
            await voice.pause_stream(chat_id)
            result = "⏸ Paused"
        elif action == "resume":
            await voice.resume_stream(chat_id)
            result = "▶️ Resumed"
        elif action == "skip":
            t = await next_track(chat_id)
            result = "⏭ Playing: " + t.title if t else "⏹ Queue completed"
        elif action == "stop":
            async with locks[chat_id]:
                for t in queues[chat_id]:
                    if t.local_file:
                        Path(t.local_file).unlink(missing_ok=True)
                queues[chat_id].clear()
                active.discard(chat_id)
                await voice.leave_group_call(chat_id)
            result = "⏹ Playback stopped"
        else:
            return
        await answer(update, result)
    except Exception as exc:
        log.error("Voice control %s failed: %s", action, short_error(exc))
        await answer(update, "⚠️ Playback operation failed (" + short_error(exc) + ").")


async def pause(update, context):
    await control(update, context, "pause")


async def resume(update, context):
    await control(update, context, "resume")


async def skip(update, context):
    await control(update, context, "skip")


async def stop(update, context):
    await control(update, context, "stop")


async def callback(update: Update, context: ContextTypes.DEFAULT_TYPE):
    data = update.callback_query.data or ""
    if not data.startswith("fb:"):
        return await update.callback_query.answer()
    await control(update, context, data.split(":", 1)[1])


async def run():
    global voice_ready, assistant_name, bot_name
    if not config.BOT_TOKEN:
        raise RuntimeError("BOT_TOKEN is required")
    bot_app = Application.builder().token(config.BOT_TOKEN).build()
    bot_app.add_handler(CommandHandler(["start", "help"], start))
    bot_app.add_handler(CommandHandler("ping", ping))
    bot_app.add_handler(CommandHandler(["play", "song"], play_audio))
    bot_app.add_handler(CommandHandler(["vplay", "video"], play_video))
    bot_app.add_handler(CommandHandler(["queue", "playlist"], queue))
    bot_app.add_handler(CommandHandler("pause", pause))
    bot_app.add_handler(CommandHandler("resume", resume))
    bot_app.add_handler(CommandHandler(["skip", "next"], skip))
    bot_app.add_handler(CommandHandler(["stop", "end"], stop))
    bot_app.add_handler(CallbackQueryHandler(callback, pattern=r"^fb:"))

    await bot_app.initialize()
    try:
        me = await bot_app.bot.get_me()
        bot_name = me.username
        log.info("HTTP_BOT_READY bot_id=%s username=@%s", me.id, me.username)
        try:
            if config.STRING1:
                await asyncio.wait_for(BABY.start(), timeout=65)
                voice_ready = True
                identity = await asyncio.wait_for(BABY.userbot1.get_me(), timeout=12)
                assistant_name = identity.username or str(identity.id)
                log.info("VOICE_ASSISTANT_READY user_id=%s", identity.id)
            else:
                log.error("VOICE_ASSISTANT_UNAVAILABLE reason=missing_STRING_SESSION")
        except Exception as exc:
            voice_ready = False
            log.error("VOICE_ASSISTANT_UNAVAILABLE reason=%s", short_error(exc))

        await bot_app.updater.start_polling(
            poll_interval=0.2, timeout=25,
            drop_pending_updates=True,
            allowed_updates=["message", "callback_query"],
        )
        await bot_app.start()
        log.info("HTTP_FALLBACK_POLLING_STARTED bot_id=%s voice_ready=%s", me.id, voice_ready)
        await asyncio.Event().wait()
    finally:
        if bot_app.updater and bot_app.updater.running:
            await bot_app.updater.stop()
        if bot_app.running:
            await bot_app.stop()
        await bot_app.shutdown()
        if voice_ready:
            try:
                await voice.stop()
            except Exception:
                pass


if __name__ == "__main__":
    asyncio.run(run())
