"""Telegram HTTP Bot API presentation and play logging."""
import logging
import re
from html import escape

from telegram import InlineKeyboardButton, InlineKeyboardMarkup

import config

log = logging.getLogger("SPOTIFY_MUSIC.http_ui")


async def delete_safely(message):
    if message is None:
        return
    try:
        await message.delete()
    except Exception:
        pass


def playback_buttons(video_id=None):
    rows = [
        [
            InlineKeyboardButton("⏸ Pause", callback_data="fb:pause"),
            InlineKeyboardButton("▶️ Resume", callback_data="fb:resume"),
        ],
        [
            InlineKeyboardButton("⏭ Skip", callback_data="fb:skip"),
            InlineKeyboardButton("⏹ Stop", callback_data="fb:stop"),
        ],
    ]
    if video_id and re.fullmatch(r"[A-Za-z0-9_-]{11}", video_id):
        rows.append([
            InlineKeyboardButton(
                "▶️ Watch on YouTube",
                url=f"https://www.youtube.com/watch?v={video_id}",
            ),
        ])
    return InlineKeyboardMarkup(rows)


def thumbnail_url(video_id):
    if video_id and re.fullmatch(r"[A-Za-z0-9_-]{11}", video_id):
        return f"https://i.ytimg.com/vi/{video_id}/mqdefault.jpg"
    return None


async def send_music_card(update, context, track, assistant, *, queued=False, position=1):
    if queued:
        headline = f"📥 Added to Queue #{position}"
    else:
        headline = "▶️ Now Playing"
    if track.requested_by:
        requested_by = f'<a href="tg://user?id={track.requested_by}">Requester</a>'
    else:
        requested_by = "Unknown"
    caption = (
        f"<b>{headline}</b>\n\n"
        f"🎵 <b>{escape(track.title)}</b>\n\n"
        f"🎬 <b>Type:</b> {'Video' if track.video else 'Audio'}\n"
        f"👤 <b>Requested by:</b> {requested_by}\n"
        f"🎙 <b>Assistant:</b> #{assistant or '?'}\n\n"
        "✨ <b>BABYMUSIC</b>"
    )
    chat_id = update.effective_chat.id
    image = thumbnail_url(track.video_id)
    markup = playback_buttons(track.video_id)
    if image:
        try:
            await context.bot.send_photo(
                chat_id=chat_id,
                photo=image,
                caption=caption,
                parse_mode="HTML",
                reply_markup=markup,
            )
            return
        except Exception as exc:
            log.warning(
                "THUMBNAIL_FAILED chat=%s error=%s",
                chat_id, type(exc).__name__,
            )
    await context.bot.send_message(
        chat_id=chat_id,
        text=caption,
        parse_mode="HTML",
        reply_markup=markup,
        disable_web_page_preview=True,
    )


async def send_play_log(update, context, track, assistant):
    if not config.LOGGER_ID:
        return
    user = update.effective_user
    username = "@" + escape(user.username) if user and user.username else "No username"
    source = track.source if track.video_id else "Telegram upload"
    message = (
        "🎵 <b>Music Play Log</b>\n\n"
        f"🎶 <b>Song:</b> {escape(track.title)}\n"
        f"💬 <b>Chat:</b> {escape(update.effective_chat.title or 'Unknown')}\n"
        f"🆔 <b>Chat ID:</b> <code>{update.effective_chat.id}</code>\n"
        f"👤 <b>User:</b> {username} "
        f"(<code>{user.id if user else 'unknown'}</code>)\n"
        f"🎙 <b>Assistant:</b> #{assistant}\n"
        f"🔗 <b>Source:</b> {escape(source)}"
    )
    try:
        await context.bot.send_message(
            config.LOGGER_ID,
            message,
            parse_mode="HTML",
            disable_web_page_preview=True,
        )
    except Exception as exc:
        log.warning("PLAY_LOG_FAILED error=%s", type(exc).__name__)


async def show_active_chats(update, context, active, titles, queue, pool):
    if not update.effective_user or update.effective_user.id != config.OWNER_ID:
        return await context.bot.send_message(
            update.effective_chat.id, "⛔ Only the bot owner can use /ac."
        )
    lines = [
        f"🎙 Active Voice Chats: {len(active)}",
        f"🟢 Assistant sessions online: {len(pool.ready)}/5",
    ]
    if not active:
        lines.append("\nNo active voice chats right now.")
    for idx, chat_id in enumerate(sorted(active), 1):
        line = (
            f"{idx}. {titles.get(chat_id, str(chat_id))[:55]}\n"
            f"   ID: {chat_id} | Assistant #{pool.assigned.get(chat_id, '?')} "
            f"| Queue: {len(queue.get(chat_id) or ())}"
        )
        if len("\n".join(lines)) + len(line) > 3800:
            lines.append("… Additional chats omitted.")
            break
        lines.append(line)
    await context.bot.send_message(update.effective_chat.id, "\n".join(lines))
