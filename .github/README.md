<img src="https://capsule-render.vercel.app/api?type=waving&color=0:6C0AC9,50:B721FF,100:FF006E&height=200&section=header&text=SPOTIFYMUSIC&fontSize=80&fontColor=ffffff&fontAlignY=38&desc=Powered%20By%20%3B-%20BabiesIQ&descAlignY=60&descSize=20&animation=fadeIn" width="100%"/>
<img src="https://capsule-render.vercel.app/api?type=waving&color=0:FF006E,50:B721FF,100:6C0AC9&height=120&section=footer&animation=fadeIn" width="100%"/>

### ──「 sᴘᴏᴛɪғʏ ダ ᴍᴜsɪᴄ 」── *The Ultimate Telegram Music Bot*

<br/>

[![Stars](https://img.shields.io/github/stars/BABY-MUSIC/SPOTIFY_MUSIC?color=1DB954&logo=github&logoColor=white&style=for-the-badge&label=STARS)](https://github.com/BABY-MUSIC/SPOTIFY_MUSIC)
[![Forks](https://img.shields.io/github/forks/BABY-MUSIC/SPOTIFY_MUSIC?color=1DB954&logo=github&logoColor=white&style=for-the-badge&label=FORKS)](https://github.com/BABY-MUSIC/SPOTIFY_MUSIC/network/members)
[![License](https://img.shields.io/badge/LICENSE-MIT-1DB954?style=for-the-badge&logo=opensourceinitiative&logoColor=white)](https://github.com/BABY-MUSIC/SPOTIFY_MUSIC/blob/master/LICENSE)
[![Python](https://img.shields.io/badge/PYTHON-3.10+-1DB954?style=for-the-badge&logo=python&logoColor=white)](https://www.python.org/)
[![Docker](https://img.shields.io/badge/DOCKER-READY-1DB954?style=for-the-badge&logo=docker&logoColor=white)](https://hub.docker.com/)
[![Telegram](https://img.shields.io/badge/TELEGRAM-BOT-1DB954?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/BabiesIQ)

<br/>

> **SPOTIFY_MUSIC** is a powerful, feature-rich Telegram music bot built by **BabiesIQ Team**.  
> Stream high-quality audio & video directly in Telegram voice chats — from YouTube, Spotify, Apple Music, SoundCloud, Resso, and more.

<br/>

[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/deploy?name=spotify-music&type=git&repository=BABY-MUSIC%2FSPOTIFY_MUSIC&branch=main&builder=dockerfile)
&nbsp;&nbsp;
[![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://dashboard.heroku.com/new?template=https://github.com/BABY-MUSIC/SPOTIFY_MUSIC)

</div>

---

## ✦ Table of Contents

| # | Section |
|---|---------|
| 01 | [What is SPOTIFY_MUSIC?](#-what-is-spotify_music) |
| 02 | [Features](#-features) |
| 03 | [Supported Platforms](#-supported-platforms) |
| 04 | [Project Architecture](#-project-architecture) |
| 05 | [Environment Variables](#-environment-variables) |
| 06 | [Deployment Guide](#-deployment-guide) |
| 07 | [Commands & Usage](#-commands--usage) |
| 08 | [Plugin System](#-plugin-system) |
| 09 | [Core Modules](#-core-modules) |
| 10 | [Utilities & Decorators](#-utilities--decorators) |
| 11 | [CI/CD & Workflows](#-cicd--workflows) |
| 12 | [Security Policy](#-security-policy) |
| 13 | [Credits & Legal](#-credits--legal) |

---

## ✦ What is SPOTIFY_MUSIC?

**SPOTIFY_MUSIC** is a Telegram Group Voice Chat Music Bot built on top of:

- **[Pyrogram](https://github.com/pyrogram/pyrogram)** — Async Telegram MTProto client for bot and userbot
- **[PyTgCalls](https://github.com/pytgcalls/pytgcalls)** — Python bindings for Telegram Group Calls (audio/video streaming)
- **[MongoDB](https://www.mongodb.com/)** — NoSQL database for persistent settings, user data, and bans
- **[BabyAPI](https://www.babyapi.pro)** — Custom proprietary API powering stream resolution and metadata

The bot supports streaming directly via the **BabyAPI** (`www.babyapi.pro`) or through local yt-dlp downloads, and can handle multiple assistant accounts simultaneously via Pyrogram session strings (`STRING1`–`STRING5`).

---

## ✦ Features

<details>
<summary><b>🎵 Music & Streaming</b> — click to expand</summary>

<br/>

| Feature | Description |
|---------|-------------|
| **Audio Streaming** | Stream any song directly in Telegram voice chats in high quality |
| **Video Streaming** | Stream full video to Telegram video chats |
| **Live Streams** | Play YouTube live streams and radio stations in real time |
| **Channel Play** | Stream audio/video in Telegram channels |
| **Queue System** | Fully automatic queue — songs line up and play one after another |
| **Loop Mode** | Loop the current song or the entire queue |
| **Shuffle** | Shuffle the queue randomly for a fresh listening order |
| **Seek** | Jump to any position in the currently playing track |
| **Speed Control** | Adjust playback speed (0.5x – 2.0x) |
| **Pause / Resume** | Pause and resume playback anytime |
| **Skip** | Skip the current track and jump to the next in queue |
| **Stop** | End playback and clear the queue |

</details>

<details>
<summary><b>📥 Downloader</b> — click to expand</summary>

<br/>

| Feature | Description |
|---------|-------------|
| `/song <name>` | Search YouTube and download a song as an **MP3** file directly to Telegram |
| `/video <name>` | Search YouTube and download a video as an **MP4** file directly to Telegram |
| **Search Results** | Shows top 10 results with inline buttons — tap to pick the one you want |
| **Animated Progress** | Real-time animated progress bar while downloading |
| **Auto Cleanup** | Downloaded files are cleaned up automatically after sending |

</details>

<details>
<summary><b>🔐 Admin & Permission System</b> — click to expand</summary>

<br/>

| Feature | Description |
|---------|-------------|
| **Admin-Only Controls** | All playback controls restricted to group admins only |
| **Sudo Users** | Bot owner can grant sudo privileges to trusted users |
| **Global Ban** | Globally ban/unban users across all groups |
| **Group Ban** | Ban specific users from using the bot in a specific group |
| **Decorator-Based Auth** | All admin checks are enforced at the decorator level — no bypass possible |

</details>

<details>
<summary><b>⚙️ Bot Configuration & Settings</b> — click to expand</summary>

<br/>

| Feature | Description |
|---------|-------------|
| **Language Support** | Multi-language support — configurable per group |
| **Play Mode Settings** | Admins can set group play modes (queue, direct, etc.) |
| **Auto-End** | Bot automatically leaves voice chat when queue is empty |
| **Heroku Auto-Update** | Optional auto-update via Heroku Git push |
| **Broadcast** | Sudo users can broadcast messages to all served groups |
| **Ping** | Check bot latency with `/ping` |

</details>

---

## ✦ Supported Platforms

| Platform | Link Type | Notes |
|----------|-----------|-------|
| 🎬 **YouTube** | URL or search query | Supports videos, playlists, and live streams |
| 🎵 **Spotify** | Spotify track/playlist/album URL | Requires `SPOTIFY_CLIENT_ID` + `SPOTIFY_CLIENT_SECRET` |
| 🍎 **Apple Music** | Apple Music URL | Auto-resolved to YouTube stream |
| 🔊 **SoundCloud** | SoundCloud URL | Direct stream support |
| 🎶 **Resso** | Resso URL | Auto-resolved |
| 📨 **Telegram** | Audio/video file in chat | Plays directly from Telegram CDN |

---

## ✦ Project Architecture

```
SPOTIFY_MUSIC/
│
├── SPOTIFY_MUSIC/                  ← Main Python package
│   │
│   ├── __init__.py                 ← Package init: exposes app, userbot, LOGGER, YouTube
│   ├── __main__.py                 ← Entry point: loads all modules, starts bot + userbots
│   ├── misc.py                     ← Misc helpers: SUDOERS filter, Heroku utils, DB init
│   ├── logging.py                  ← Custom logger factory
│   │
│   ├── core/                       ← Core infrastructure
│   │   ├── bot.py                  ← BABY(Client): main Pyrogram bot instance
│   │   ├── call.py                 ← PyTgCalls voice chat engine (BABY call handler)
│   │   ├── mongo.py                ← MongoDB connection singleton
│   │   ├── git.py                  ← GitHub/Heroku git update utilities
│   │   ├── dir.py                  ← Directory management (downloads, temp files)
│   │   ├── sections.py             ← Plugin discovery and section loading
│   │   └── userbot.py              ← Assistant userbot clients (STRING1–STRING5)
│   │
│   ├── platforms/                  ← Platform-specific stream resolvers
│   │   ├── Youtube.py              ← YouTube search + stream URL resolver (yt-dlp + BabyAPI)
│   │   ├── Spotify.py              ← Spotify → YouTube resolver via spotipy
│   │   ├── Apple.py                ← Apple Music → YouTube resolver
│   │   ├── Soundcloud.py           ← SoundCloud stream resolver
│   │   ├── Resso.py                ← Resso → YouTube resolver
│   │   ├── Telegram.py             ← Telegram file CDN stream handler
│   │   └── Carbon.py               ← Carbon/custom stream support
│   │
│   ├── plugins/                    ← All bot command plugins
│   │   ├── play/                   ← Playback commands
│   │   │   ├── play.py             ← /play command (audio)
│   │   │   ├── live.py             ← /live command (live streams)
│   │   │   ├── channel.py          ← /cplay command (channel playback)
│   │   │   ├── playmode.py         ← /playmode settings
│   │   │   └── playstream.py       ← /vplay, /vstream (video)
│   │   │
│   │   ├── admins/                 ← Admin-only playback controls
│   │   │   ├── pause.py            ← /pause
│   │   │   ├── resume.py           ← /resume
│   │   │   ├── skip.py             ← /skip
│   │   │   ├── stop.py             ← /stop
│   │   │   ├── seek.py             ← /seek <seconds>
│   │   │   ├── loop.py             ← /loop
│   │   │   ├── shuffle.py          ← /shuffle
│   │   │   ├── speed.py            ← /speed <value>
│   │   │   └── callback.py         ← Inline button callbacks for player UI
│   │   │
│   │   ├── bot/                    ← General bot commands
│   │   │   ├── start.py            ← /start
│   │   │   ├── help.py             ← /help
│   │   │   └── settings.py         ← /settings
│   │   │
│   │   ├── misc/                   ← Miscellaneous background tasks
│   │   │   ├── broadcast.py        ← /broadcast (sudo only)
│   │   │   ├── seeker.py           ← Auto-seek/progress watcher
│   │   │   └── watcher.py          ← Voice chat event watcher
│   │   │
│   │   ├── sudo/
│   │   │   └── sudoers.py          ← /addsudo, /rmsudo, /sudolist
│   │   │
│   │   └── tools/
│   │       ├── downloader.py       ← /song, /video (offline download)
│   │       ├── language.py         ← /setlang
│   │       ├── logger.py           ← Error logger utility
│   │       ├── ping.py             ← /ping
│   │       ├── reload.py           ← /reload
│   │       └── userid.py           ← /id
│   │
│   └── utils/                      ← Shared utility helpers
│       ├── database.py             ← All MongoDB CRUD operations
│       ├── admin_check.py          ← Admin verification helper
│       ├── baby_ban.py             ← Global/group ban management
│       ├── baby_font.py            ← Custom font rendering for thumbnails
│       ├── channelplay.py          ← Channel playback helpers
│       ├── thumbnails.py           ← Dynamic thumbnail generator
│       ├── formatters.py           ← Time/size formatters
│       └── decorators/
│           ├── admins.py           ← @admin_required decorator
│           ├── language.py         ← @language decorator for i18n
│           └── play.py             ← @check_blacklist, @require_voice decorator
│
├── Genstring.py                    ← Interactive session string generator
├── Dockerfile                      ← Docker container definition
├── Procfile                        ← Heroku worker process definition
└── .github/
    ├── README.md                   ← GitHub-facing README
    ├── SECURITY.md                 ← Security policy
    ├── FUNDING.yml                 ← Sponsorship config
    └── workflows/
        ├── docker-image.yml        ← Docker build CI on every push to main
        └── codeql.yml              ← GitHub CodeQL security scan (Python)
```

---

## ✦ Environment Variables

<details>
<summary><b>🔑 Required Variables</b> — click to expand</summary>

<br/>

| Variable | Description | How to Get |
|----------|-------------|------------|
| `API_ID` | Telegram API ID | [my.telegram.org](https://my.telegram.org/auth) → App configuration |
| `API_HASH` | Telegram API Hash | [my.telegram.org](https://my.telegram.org/auth) → App configuration |
| `BOT_TOKEN` | Your bot's token | [@BotFather](https://t.me/BotFather) → `/newbot` |
| `OWNER_ID` | Your Telegram User ID | [@userinfobot](https://t.me/userinfobot) |
| `MONGO_DB_URI` | MongoDB connection string | [MongoDB Atlas](https://cloud.mongodb.com/) (free tier works) |
| `LOGGER_ID` | Telegram Chat ID for bot logs | A private group or channel where the bot is admin |
| `STRING1` | Pyrogram session string (Assistant 1) | Run `Genstring.py` — explained below |
| `BASE_URL` | BabyAPI base URL | `https://api.babyapi.pro` (contact team) |
| `API_KEY` | BabyAPI key | Obtain from [babyapi.pro](https://www.babyapi.pro) |

</details>

<details>
<summary><b>🎵 Optional — Spotify Integration</b> — click to expand</summary>

<br/>

| Variable | Description | How to Get |
|----------|-------------|------------|
| `SPOTIFY_CLIENT_ID` | Spotify App Client ID | [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/) → Create App |
| `SPOTIFY_CLIENT_SECRET` | Spotify App Client Secret | Same as above |

> Without these, Spotify URLs will not resolve. The bot falls back gracefully to YouTube search when omitted.

</details>

<details>
<summary><b>🚀 Optional — Heroku Auto-Update</b> — click to expand</summary>

<br/>

| Variable | Description |
|----------|-------------|
| `HEROKU_API_KEY` | Your Heroku account API key (Account Settings → API Key) |
| `HEROKU_APP_NAME` | Your Heroku app name (e.g. `my-spotify-bot`) |

> These enable the `/update` command which automatically pulls the latest code and redeploys.

</details>

<details>
<summary><b>🤖 Optional — Multiple Assistant Accounts (STRING2–STRING5)</b> — click to expand</summary>

<br/>

The bot uses **Pyrogram userbot sessions** (assistant accounts) to join voice chats on behalf of users. You can configure up to **5 assistant accounts** simultaneously.

| Variable | Description |
|----------|-------------|
| `STRING1` | Session string for Assistant Account 1 (required) |
| `STRING2` | Session string for Assistant Account 2 (optional) |
| `STRING3` | Session string for Assistant Account 3 (optional) |
| `STRING4` | Session string for Assistant Account 4 (optional) |
| `STRING5` | Session string for Assistant Account 5 (optional) |

> More assistants = more groups can use the bot simultaneously without conflicts.

**How to generate a session string:**

```bash
python3 Genstring.py
```

Follow the prompts — enter your API ID, API Hash, and phone number. The script generates a session string to paste into your vars.

</details>

---

## ✦ Deployment Guide

<details>
<summary><b>🐳 Deploy with Docker</b> — click to expand</summary>

<br/>

The project ships with a production-ready `Dockerfile` using Python 3.10-slim.

```bash
# 1. Clone the repository
git clone https://github.com/BABY-MUSIC/SPOTIFY_MUSIC
cd SPOTIFY_MUSIC

# 2. Copy environment template and fill in your values
cp sample.env .env
nano .env

# 3. Build the Docker image
docker build -t spotify-music .

# 4. Run the container
docker run --env-file .env spotify-music
```

**What the Dockerfile does:**
- Starts from `python:3.10-slim`
- Installs system dependencies: `git`, `curl`, `ffmpeg`
- Installs Node.js 20 (required by some stream helpers)
- Installs all Python requirements via `pip`
- Runs `bash start` as the entrypoint

</details>

<details>
<summary><b>☁️ Deploy on Heroku</b> — click to expand</summary>

<br/>

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://dashboard.heroku.com/new?template=https://github.com/BABY-MUSIC/SPOTIFY_MUSIC)

**Manual steps:**

```bash
# 1. Login to Heroku CLI
heroku login

# 2. Create your app
heroku create your-app-name

# 3. Set environment variables
heroku config:set API_ID=your_api_id --app your-app-name
heroku config:set API_HASH=your_api_hash --app your-app-name
heroku config:set BOT_TOKEN=your_bot_token --app your-app-name
# ... set all required vars

# 4. Push to Heroku
git push heroku main

# 5. Scale the worker dyno
heroku ps:scale worker=1 --app your-app-name
```

> The `Procfile` runs: `worker: python app.py`

</details>

<details>
<summary><b>🌐 Deploy on Koyeb</b> — click to expand</summary>

<br/>

[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/deploy?name=spotify-music&type=git&repository=BABY-MUSIC%2FSPOTIFY_MUSIC&branch=main&builder=dockerfile)

1. Click the button above
2. Connect your GitHub account
3. Fill in all required environment variables in the Koyeb dashboard
4. Koyeb will build using the `Dockerfile` and deploy automatically

</details>

<details>
<summary><b>🖥️ Deploy on Linux VPS</b> — click to expand</summary>

<br/>

**Step 1 — Update system & install dependencies:**

```bash
sudo apt-get update && sudo apt-get upgrade -y
sudo apt-get install python3-pip ffmpeg git -y
sudo pip3 install -U pip
```

**Step 2 — Install Node.js:**

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.38.0/install.sh | bash
source ~/.bashrc
nvm install v18
```

**Step 3 — Clone the repository:**

```bash
git clone https://github.com/BABY-MUSIC/SPOTIFY_MUSIC
cd SPOTIFY_MUSIC
```

**Step 4 — Install Python requirements:**

```bash
pip3 install -U -r requirements.txt
```

**Step 5 — Configure environment:**

```bash
cp sample.env .env
vi .env
# Press I to enter insert mode
# Fill in all your variables
# Press Ctrl+C, then type :wq to save
```

**Step 6 — Run inside tmux (keeps running after disconnect):**

```bash
sudo apt install tmux -y
tmux new -s music
bash start
# Press Ctrl+B then D to detach safely
```

</details>

---

## ✦ Commands & Usage

<details>
<summary><b>▶️ Play Commands</b> — click to expand</summary>

<br/>

| Command | Description |
|---------|-------------|
| `/play <song name or URL>` | Play audio in the group voice chat — searches YouTube or resolves the URL |
| `/vplay <song name or URL>` | Play video in the group video chat |
| `/live <YouTube live URL>` | Stream a YouTube live broadcast in the voice chat |
| `/vstream <stream URL>` | Stream a video live broadcast in the video chat |
| `/cplay <song>` | Play music in a linked channel voice chat |
| `/song <name>` | Search and download a song as MP3 to Telegram |
| `/video <name>` | Search and download a video as MP4 to Telegram |

</details>

<details>
<summary><b>🎛️ Admin Playback Controls</b> — click to expand (admin only)</summary>

<br/>

| Command | Description |
|---------|-------------|
| `/pause` | Pause the currently playing track |
| `/resume` | Resume a paused track |
| `/skip` | Skip the current track and play the next in queue |
| `/stop` | Stop playback and clear the entire queue |
| `/seek <seconds>` | Jump to a specific position in the current track |
| `/loop` | Toggle loop mode (off → single track → full queue) |
| `/shuffle` | Shuffle all tracks remaining in the queue |
| `/speed <value>` | Change playback speed (e.g. `/speed 1.5`) |

</details>

<details>
<summary><b>🛠️ Tools & Utility Commands</b> — click to expand</summary>

<br/>

| Command | Description |
|---------|-------------|
| `/ping` | Check the bot's response latency |
| `/id` | Get your Telegram User ID or the chat ID |
| `/reload` | Reload all bot plugins without restarting |
| `/setlang <language_code>` | Set the bot language for your group |
| `/help` | Display a categorized help menu |
| `/settings` | Open the group settings panel |
| `/start` | Start the bot in DM and get the welcome message |

</details>

<details>
<summary><b>👑 Sudo Commands</b> — click to expand (owner/sudo only)</summary>

<br/>

| Command | Description |
|---------|-------------|
| `/addsudo <user_id>` | Grant sudo privileges to a user |
| `/rmsudo <user_id>` | Revoke sudo privileges from a user |
| `/sudolist` | List all current sudo users |
| `/gban <user_id>` | Globally ban a user from all groups |
| `/ungban <user_id>` | Remove a global ban |
| `/broadcast <message>` | Broadcast a message to all groups the bot serves |

</details>

---

## ✦ Plugin System

The bot uses a **modular plugin architecture**. All plugins live under `SPOTIFY_MUSIC/plugins/` and are discovered automatically at startup by `core/sections.py`.

<details>
<summary><b>📂 Plugin Categories Explained</b> — click to expand</summary>

<br/>

| Category | Folder | Purpose |
|----------|--------|---------|
| **Play** | `plugins/play/` | All playback-related commands: audio, video, live, channel, stream |
| **Admins** | `plugins/admins/` | Admin-restricted playback controls and inline button callbacks |
| **Bot** | `plugins/bot/` | General user-facing commands: start, help, settings |
| **Misc** | `plugins/misc/` | Background utilities: broadcast, auto-seeker, voice chat watcher |
| **Sudo** | `plugins/sudo/` | Owner-level management: sudo user list, add/remove sudo |
| **Tools** | `plugins/tools/` | Utility commands: downloader, ping, reload, language, user ID |

Each plugin file registers Pyrogram message handlers using `@app.on_message(filters.command(...))`. The decorator system (`utils/decorators/`) injects admin checks and language context automatically before the handler runs.

</details>

---

## ✦ Core Modules

<details>
<summary><b>🤖 core/bot.py — The Bot Client</b> — click to expand</summary>

<br/>

Defines the `BABY` class, which extends `pyrogram.Client`. This is the main bot account that handles all incoming commands and messages.

```python
class BABY(Client):
    def __init__(self):
        super().__init__(
            name="SPOTIFY_MUSIC",
            api_id=config.API_ID,
            api_hash=config.API_HASH,
            bot_token=config.BOT_TOKEN,
            in_memory=True,
            max_concurrent_transmissions=7,
        )
```

On startup, it sends a boot message to `LOGGER_ID` confirming the bot is live.

</details>

<details>
<summary><b>📞 core/call.py — The Voice Chat Engine</b> — click to expand</summary>

<br/>

The most critical module. Wraps **PyTgCalls** to handle:

- Joining and leaving Telegram group voice chats
- Streaming `AudioPiped` or `AudioVideoPiped` to the voice chat
- Handling stream end events → auto-plays next queue item
- Managing multiple assistant accounts (one per active group)
- Loop mode, seek, speed, and volume control
- Auto-end when the queue is empty (respects `is_autoend()` setting)

Key functions:
| Function | Description |
|----------|-------------|
| `stream_call()` | Join voice chat and start streaming |
| `force_stop_stream()` | Immediately stop and leave the voice chat |
| `skip_stream()` | Skip to the next item in queue |
| `pause_stream()` | Pause current playback |
| `resume_stream()` | Resume paused playback |
| `seek_stream()` | Seek to a timestamp |
| `speed_change()` | Adjust playback speed |

</details>

<details>
<summary><b>🍃 core/mongo.py — Database Connection</b> — click to expand</summary>

<br/>

Establishes a singleton MongoDB connection using `motor` (async MongoDB driver for Python). All collections are accessed via the `mongodb` object exported from this module.

The database stores:
- Per-group settings (language, play mode, auto-end toggle, loop state)
- Active chats and active video chats
- Sudo user list
- Global ban list and per-group ban list
- Queue data per group

</details>

<details>
<summary><b>👥 core/userbot.py — Assistant Accounts</b> — click to expand</summary>

<br/>

Creates up to **5 Pyrogram userbot clients** from `STRING1`–`STRING5`. These assistant accounts:

- Are the ones that actually **join** the Telegram voice chat (bots can't join VC directly)
- Handle voice chat joining via `group_call_participant` updates
- Are automatically rotated per group to distribute load

The `group_assistant()` utility in `utils/database.py` returns the correct assistant for any given group.

</details>

---

## ✦ Utilities & Decorators

<details>
<summary><b>🗄️ utils/database.py — All Database Operations</b> — click to expand</summary>

<br/>

Central module for all MongoDB read/write operations. Key functions grouped by purpose:

**Settings:**
```python
get_lang(chat_id)          # Get group language
set_lang(chat_id, lang)    # Set group language
is_autoend(chat_id)        # Check if auto-end is enabled
get_loop(chat_id)          # Get loop mode (0=off, 1=single, 2=queue)
set_loop(chat_id, mode)    # Set loop mode
```

**Active Chats:**
```python
add_active_chat(chat_id)           # Mark group as having active audio
add_active_video_chat(chat_id)     # Mark group as having active video
remove_active_chat(chat_id)
remove_active_video_chat(chat_id)
music_on(chat_id)                  # Check if music is playing
```

**Bans:**
```python
get_gbanned()              # Get list of globally banned users
is_gbanned(user_id)        # Check if user is globally banned
add_gban(user_id)
remove_gban(user_id)
get_banned_users()         # Per-group bans
```

**Assistants:**
```python
group_assistant(chat_id)   # Get assigned assistant for group
```

</details>

<details>
<summary><b>🔒 utils/decorators/ — Auth & Context Injection</b> — click to expand</summary>

<br/>

Three custom decorators are applied to plugin handlers:

**`@admin_required` (`decorators/admins.py`)**
Verifies the sender is a group admin before allowing admin commands like `/pause`, `/skip`, `/stop`, etc. Returns an error message if the user is not an admin.

**`@language` (`decorators/language.py`)**
Injects the group's configured language strings into the handler as a `_` argument, allowing all reply messages to be automatically localized.

**`@check_blacklist` / `@require_voice` (`decorators/play.py`)**
- `@check_blacklist` — rejects globally banned users at the handler level
- `@require_voice` — ensures the bot is in a voice chat before processing play commands

</details>

<details>
<summary><b>🖼️ utils/thumbnails.py — Dynamic Thumbnails</b> — click to expand</summary>

<br/>

Generates custom now-playing thumbnail images displayed when a new track starts. Uses:
- Custom fonts from `SPOTIFY_MUSIC/assets/` (multiple `.ttf` font files)
- Track metadata (title, artist, duration, thumbnail URL)
- `Pillow` for image manipulation and text rendering

</details>

---

## ✦ CI/CD & Workflows

<details>
<summary><b>🐳 Docker Image CI (.github/workflows/docker-image.yml)</b> — click to expand</summary>

<br/>

Runs on every push or pull request to `main`. Builds the Docker image to verify the `Dockerfile` is valid and all dependencies install correctly.

```yaml
on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]
```

</details>

<details>
<summary><b>🔍 CodeQL Security Scan (.github/workflows/codeql.yml)</b> — click to expand</summary>

<br/>

Runs GitHub's **CodeQL** analysis on the Python codebase:
- Scheduled every Friday at 01:25 UTC
- Also runs on every push and pull request to `main`
- Scans for common vulnerability patterns (injection, unsafe deserialization, etc.)
- Results appear under the **Security** tab of the repository

</details>

---

## ✦ Security Policy

<details>
<summary><b>🛡️ Reporting a Vulnerability</b> — click to expand</summary>

<br/>

| Version | Supported |
|---------|-----------|
| 5.1.x | ✅ Active support |
| 5.0.x | ❌ End of life |
| 4.0.x | ✅ Security patches only |
| < 4.0 | ❌ No support |

**To report a vulnerability**, contact the team directly via:
- Telegram: [@BabiesIQ](https://t.me/BabiesIQ)
- GitHub: Open a private security advisory

Do **not** post vulnerability details in public issues. You will receive an acknowledgment within 48 hours and a resolution timeline within 7 days of confirmed impact.

</details>

---

## ✦ Credits & Legal

<div align="center">

| | |
|---|---|
| **Project** | SPOTIFY_MUSIC — Public Music Bot Repository |
| **Author** | BabiesIQ Team |
| **Telegram** | [@BabiesIQ](https://t.me/BabiesIQ) |
| **API** | [www.babyapi.pro](https://www.babyapi.pro) |
| **Repository** | [github.com/BABY-MUSIC/SPOTIFY_MUSIC](https://github.com/BABY-MUSIC/SPOTIFY_MUSIC) |
| **License** | MIT |

</div>

<details>
<summary><b>📄 Legal Notice</b> — click to expand</summary>

<br/>

- Use, upload, and modification are **at your own risk**
- Only editing `.env` / `config` vars is permitted — **do not modify core files**
- Keep the original copyright header intact if you fork this project
- The developers are **not responsible** for bans, API blocks, or any resulting damage
- Internal protection mechanisms may exist — unauthorized modifications may cause the system to stop functioning
- **Use only the official API**: [www.babyapi.pro](https://www.babyapi.pro)

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for full terms.

</details>

---

<div align="center">

```
▓▒░ ʙ ᴀ ʙ ɪ ᴇ sＩＱ ░▒▓  s ᴇ ᴄ ᴜ ʀ ᴇ  ▓▒░ ɴ ᴇ ᴛ ᴡ ᴏ ʀ ᴋ ░▒▓
```

**Made with ♪ by BabiesIQ Team**

[![Telegram](https://img.shields.io/badge/Join-BabiesIQ-1DB954?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/BabiesIQ)
[![API](https://img.shields.io/badge/Powered%20by-BabyAPI-1DB954?style=for-the-badge&logo=fastapi&logoColor=white)](https://www.babyapi.pro)

</div>
