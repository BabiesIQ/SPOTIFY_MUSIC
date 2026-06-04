from pyrogram import filters, Client
from pyrogram.types import Message, InlineKeyboardMarkup, InlineKeyboardButton
from unidecode import unidecode

from SPOTIFY_MUSIC import app
from SPOTIFY_MUSIC.misc import SUDOERS
from SPOTIFY_MUSIC.utils.database import (
    get_active_chats,
    get_active_video_chats,
    remove_active_chat,
    remove_active_video_chat,
)
from SPOTIFY_MUSIC.utils.decorators.language import language

@app.on_message(filters.command(["activevc", "activevoice", "vc"]) & SUDOERS)
@language
async def activevc(_, message: Message, _str):
    mystic = await message.reply_text(_str["active_1"])
    served_chats = await get_active_chats()
    text = ""
    j = 0
    for x in served_chats:
        try:
            chat = await app.get_chat(x)
            title = unidecode(chat.title).upper()
            link = f"<a href=https://t.me/{chat.username}>{title}</a>" if chat.username else title
            text += f"<b>{j + 1}.</b> {link}\n"
            j += 1
        except:
            await remove_active_chat(x)
    if not text:
        await mystic.edit_text(_str["active_2"].format(app.mention))
    else:
        await mystic.edit_text(
            _str["active_3"].format(text),
            disable_web_page_preview=True,
        )

@app.on_message(filters.command(["activev", "activevideo", "avc"]) & SUDOERS)
@language
async def activevi_(_, message: Message, _str):
    mystic = await message.reply_text(_str["active_4"])
    served_chats = await get_active_video_chats()
    text = ""
    j = 0
    for x in served_chats:
        try:
            chat = await app.get_chat(x)
            title = unidecode(chat.title).upper()
            link = f"<a href=https://t.me/{chat.username}>{title}</a>" if chat.username else title
            text += f"<b>{j + 1}.</b> {link} [<code>{x}</code>]\n"
            j += 1
        except:
            await remove_active_video_chat(x)
    if not text:
        await mystic.edit_text(_str["active_5"].format(app.mention))
    else:
        await mystic.edit_text(
            _str["active_6"].format(text),
            disable_web_page_preview=True,
        )

@app.on_message(filters.command(["ac", "av"]) & SUDOERS)
@language
async def active_count(client: Client, message: Message, _str):
    ac_audio = str(len(await get_active_chats()))
    ac_video = str(len(await get_active_video_chats()))
    await message.reply_text(
        _str["active_7"].format(ac_audio, ac_video),
        reply_markup=InlineKeyboardMarkup(
            [[InlineKeyboardButton(_str["active_8"], callback_data="close")]]
        )
    )
