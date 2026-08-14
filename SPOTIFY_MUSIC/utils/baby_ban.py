
# ======================================================================
# ||                                                               ||
# ||   ██████╗  █████╗ ██████╗ ██╗   ██╗███████╗███████╗██╗ ██████╗  ||
# ||   ██╔══██╗██╔══██╗██╔══██╗██║   ██║██╔════╝██╔════╝██║██╔═══██╗ ||
# ||   ██████╔╝███████║██████╔╝██║   ██║█████╗  ███████╗██║██║   ██║ ||
# ||   ██╔══██╗██╔══██║██╔══██╗██║   ██║██╔══╝  ╚════██║██║██║▄▄ ██║ ||
# ||   ██████╔╝██║  ██║██████╔╝╚██████╔╝███████╗███████║██║╚██████╔╝ ||
# ||   ╚═════╝ ╚═╝  ╚═╝╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚═╝ ╚══▀▀═╝  ||
# ║    ▓▒░ ʙ ᴀ ʙ ɪ ᴇ sＩＱ ░▒▓  s ᴇ ᴄ ᴜ ʀ ᴇ  ▓▒░ ɴ ᴇ ᴛ ᴡ ᴏ ʀ ᴋ ░▒▓    ║
# ||                                                               ||
# ======================================================================
# || PROJECT  : SPOTIFY_MUSIC Public Music Repository                  ||
# || AUTHOR   : BabiesIQ Team                                      ||
# || REPO     : github.com/BABY-MUSIC/SPOTIFY_MUSIC                ||
# || API      : www.babyapi.pro                                    ||
# || TELEGRAM : t.me/BabiesIQ                                      ||
# ----------------------------------------------------------------------
# || LEGAL NOTICE                                                  ||
# || Use / upload / modify at your own risk.                       ||
# || Only config /.env edit allowed.                               ||
# || Do not modify core files.                                     ||
# || Keep this header if forked.                                   ||
# || Dev not responsible for ban / damage / api block.             ||
# ----------------------------------------------------------------------
# || SECURITY                                                      ||
# || Internal protection may exist.                                ||
# || Unauthorized change may stop system.                          ||
# || Use official API only -> www.babyapi.pro                      ||
# ======================================================================


from pyrogram import filters
from SPOTIFY_MUSIC.utils.admin_check import admin_check


USE_AS_BOT = True

def f_sudo_filter(filt, client, message):
    return bool(
        (
            (message.from_user and message.from_user.id in SUDO_USERS)
            or (message.sender_chat and message.sender_chat.id in SUDO_USERS)
        )
        and
       
        not message.edit_date
    )


sudo_filter = filters.create(func=f_sudo_filter, name="SudoFilter")


def onw_filter(filt, client, message):
    if USE_AS_BOT:
        return bool(
            True
            and 
           
            not message.edit_date
        )
    else:
        return bool(
            message.from_user
            and message.from_user.is_self
            and
           
            not message.edit_date
        )


f_onw_fliter = filters.create(func=onw_filter, name="OnwFilter")


async def admin_filter_f(filt, client, message):
    return (
       
        not message.edit_date
        and await admin_check(message)
    )


admin_filter = filters.create(func=admin_filter_f, name="AdminFilter")