<div align="center">

# SPOTIFY_MUSIC 🎧

**A blazing-fast Telegram Music Bot — powered by Go**

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org/)
[![License](https://img.shields.io/badge/License-GPLv3-success?style=flat-square)](LICENSE)
[![Stars](https://img.shields.io/github/stars/BabiesIQ/SPOTIFY_MUSIC?style=flat-square&color=yellow)](https://github.com/BabiesIQ/SPOTIFY_MUSIC/stargazers)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docs.docker.com/)

Stream music directly in your Telegram voice chats — YouTube, Spotify links, direct URLs, and more.  
Built for speed and reliability using **gotdbot**, **ntgcalls**, and **MongoDB**.

<a href="https://heroku.com/deploy?template=https://github.com/BabiesIQ/SPOTIFY_MUSIC">
  <img src="https://www.herokucdn.com/deploy/button.svg" alt="Deploy to Heroku" height="30">
</a>

</div>

---

## 🎯 What It Does

| Feature | Details |
|---------|---------|
| 🎵 Multi-source playback | YouTube, Spotify, M3U8 streams, direct links |
| 🎛 Full playback control | Play, pause, resume, seek, skip, loop, shuffle |
| 📋 Queue management | Multi-track queue with skip & reorder support |
| 🔊 Volume & mute | Per-chat mute/unmute and volume control |
| 👥 Multi-assistant | Up to 10 concurrent session strings |
| 🌍 Localization ready | Multi-language support baked in |
| 🐳 Docker native | One-command deployment |

---

## ⚡ Quick Start

### 1. Get Your Credentials

- **Telegram API** → [my.telegram.org](https://my.telegram.org) → `API_ID` and `API_HASH`
- **Bot Token** → [@BotFather](https://t.me/BotFather) on Telegram
- **MongoDB URI** → [MongoDB Atlas](https://www.mongodb.com/cloud/atlas) (free tier)
- **Session String** → Generate with any Pyrogram/Telethon session generator

### 2. Clone & Configure

```bash
git clone https://github.com/BabiesIQ/SPOTIFY_MUSIC.git
cd SPOTIFY_MUSIC
cp sample.env .env
nano .env   # fill in your values
```

### 3. Deploy

**Docker (recommended):**
```bash
docker compose up -d --build
```

**Manual (Go):**
```bash
go run setup_ntgcalls.go
go mod tidy
go build -o spotify-music .
./spotify-music
```

> 📖 Full deployment guide → [docs/deployment.md](docs/deployment.md)

---

## 🔑 Environment Variables

<details>
<summary>Click to expand all variables</summary>

| Variable | Description | Required |
|----------|-------------|:--------:|
| `API_ID` | Telegram API ID from my.telegram.org | ✅ |
| `API_HASH` | Telegram API Hash | ✅ |
| `TOKEN` | Bot token from @BotFather | ✅ |
| `MONGO_URI` | MongoDB connection string | ✅ |
| `OWNER_ID` | Your Telegram user ID | ❌ |
| `LOGGER_ID` | Group ID where bot sends logs | ✅ |
| `SESSION_TYPE` | `pyrogram`, `telethon`, or `gogram` | ✅ |
| `STRING1` | First assistant session string | ✅ |
| `STRING2`–`STRING10` | Additional session strings | ❌ |
| `MAX_FILE_SIZE` | Max download size in bytes (default: 500MB) | ❌ |
| `API_URL` | External music API base URL | ❌ |
| `API_KEY` | Key for external music API | ❌ |
| `COOKIES_URL` | Cookie file URL for authenticated streams | ❌ |
| `SUPPORT_GROUP` | Your Telegram support group link | ❌ |
| `SUPPORT_CHANNEL` | Your Telegram channel link | ❌ |
| `SONG_DURATION_LIMIT` | Max song length in seconds | ❌ |
| `DEVS` | Space-separated developer Telegram IDs | ❌ |

</details>

---

## 🤖 Bot Commands

<details>
<summary>Admin & sudo commands</summary>

| Command | Action |
|---------|--------|
| `/play [name or URL]` | Play audio in voice chat |
| `/vplay [name or URL]` | Play video in voice chat |
| `/skip` | Skip to next track |
| `/pause` / `/resume` | Pause or resume playback |
| `/end` | Stop and clear the queue |
| `/mute` / `/unmute` | Mute or unmute the assistant |
| `/volume [0-200]` | Set playback volume |
| `/seek [seconds]` | Jump to position in track |
| `/queue` | Show the current queue |
| `/shuffle` | Shuffle the queue |
| `/loop` | Toggle loop mode |
| `/auth` / `/unauth` | Grant or revoke user permissions |
| `/settings` | Open bot settings panel |

</details>

<details>
<summary>User commands</summary>

| Command | Action |
|---------|--------|
| `/start` | Check if the bot is alive |
| `/ping` | Measure bot latency |
| `/help` | Show the help menu |

</details>

---

## 📚 Documentation

| Guide | Description |
|-------|-------------|
| [Deployment Guide](docs/deployment.md) | Docker, manual, Heroku, systemd — all options |
| [Cookie Authentication](docs/session-cookies.md) | Set up browser cookies for authenticated streams |

---

## 🤝 Contributing

Pull requests are welcome! Please open an issue first to discuss major changes.

1. Fork the repo
2. Create your branch: `git checkout -b feat/your-feature`
3. Commit your changes: `git commit -m "feat: add your feature"`
4. Push: `git push origin feat/your-feature`
5. Open a Pull Request

---

<div align="center">

Built with ❤️ by [@BabiesIQ](https://github.com/BabiesIQ)

</div>
